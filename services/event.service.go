package services

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	"github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventService interface {
	CreateEvent(c *fiber.Ctx) error
	GetLiveEvents(ctx context.Context, userID string) ([]gen.GetLiveEventsRow, error)
	GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]gen.GetLiveEventsDetailRow, error)
	GetEventSummary(ctx context.Context, projectID string, userID string) (*gen.GetEventSummaryRow, error)
	GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error)
	GetWeeklyEventsChart(ctx context.Context, projectID string, userID string) (*entities.EventSummaryChart, error)
	GetEventTypeChart(ctx context.Context, projectID string, userID string) ([]gen.GetPercentageEventsTypeRow, error)
	GetEventLabelChart(ctx context.Context, projectID string, userID string) ([]gen.GetPercentageEventsLabelRow, error)
	CountUserMonthlyEvents(ctx context.Context, userID uuid.UUID) (int64, error)
}

type EventServiceImpl struct {
	UtilService UtilService
	Repo        *gen.Queries
}

func InitEventService(utilService UtilService, repo *gen.Queries) EventServiceImpl {
	return EventServiceImpl{
		UtilService: utilService,
		Repo:        repo,
	}
}

func (s *EventServiceImpl) CreateEvent(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByPublicKeyRow)
	var input dto.CreateEventInput

	if err := c.BodyParser(&input); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := s.UtilService.ValidateInput(input); err != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	// check if the project id exist within user.
	// if not dont proceed further.
	projectUUID := uuid.MustParse(input.ProjectID)
	exist, _ := s.Repo.CheckProjectWithinUserID(context.Background(), gen.CheckProjectWithinUserIDParams{
		ID:     projectUUID,
		UserID: user.ID,
	})
	if !exist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project not found"})
	}

	// Get user IP
	userIP := c.IP()
	userLoc := s.UtilService.LookupIP(userIP)

	// bind user ip addr and loc
	input.IPAddr = userIP
	if userLoc != nil {
		input.Country = userLoc.Country.Names["en"]
		input.Region = userLoc.Continent.Names["en"]
		input.City = userLoc.City.Names["en"]
	}

	payload := gen.CreateEventParams{
		// i need to insert the dto payload here, but its tedious to do it manually, F
		EventType:        input.EventType,
		EventLabel:       pgtype.Text{String: input.EventLabel, Valid: input.EventLabel != ""},
		PageUrl:          pgtype.Text{String: input.PageURL, Valid: input.PageURL != ""},
		ElementPath:      pgtype.Text{String: input.ElementPath, Valid: input.ElementPath != ""},
		ElementType:      pgtype.Text{String: input.ElementType, Valid: input.ElementType != ""},
		IpAddr:           s.UtilService.ParseIP(input.IPAddr),
		UserAgent:        pgtype.Text{String: input.UserAgent, Valid: input.UserAgent != ""},
		BrowserName:      pgtype.Text{String: input.BrowserName, Valid: input.BrowserName != ""},
		Country:          pgtype.Text{String: input.Country, Valid: input.Country != ""},
		Region:           pgtype.Text{String: input.Region, Valid: input.Region != ""},
		City:             pgtype.Text{String: input.City, Valid: input.City != ""},
		SessionID:        pgtype.Text{String: input.SessionID, Valid: input.SessionID != ""},
		DeviceType:       pgtype.Text{String: input.DeviceType, Valid: input.DeviceType != ""},
		TimeOnPage:       pgtype.Int4{Int32: int32(input.TimeOnPage), Valid: true},
		ScreenResolution: pgtype.Text{String: input.ScreenResolution, Valid: input.ScreenResolution != ""},
		FiredAt:          s.UtilService.ParseTimestamp(input.FiredAt),
		ReceivedAt:       time.Now(),
		UserID:           user.ID,
		ProjectID:        uuid.MustParse(input.ProjectID),
	}

	if err := s.Repo.CreateEvent(context.Background(), payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *EventServiceImpl) GetLiveEvents(ctx context.Context, userID string) ([]gen.GetLiveEventsRow, error) {
	userUUID := uuid.MustParse(userID)
	return s.Repo.GetLiveEvents(ctx, userUUID)
}

func (s *EventServiceImpl) GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]gen.GetLiveEventsDetailRow, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)

	return s.Repo.GetLiveEventsDetail(ctx, gen.GetLiveEventsDetailParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
}

func (s *EventServiceImpl) GetEventSummary(ctx context.Context, projectID string, userID string) (*gen.GetEventSummaryRow, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)

	row, err := s.Repo.GetEventSummary(ctx, gen.GetEventSummaryParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (s *EventServiceImpl) GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)

	var summary entities.EventDetail

	// event summary total numbering
	tldr, err := s.Repo.GetTotalEventSummary(ctx, gen.GetTotalEventSummaryParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	summary.TotalEvents = int(tldr.TotalEvents)
	summary.TotalEventType = int(tldr.TotalEventType)
	summary.TotalUniqueUsers = int(tldr.TotalUniqueUsers)
	summary.TotalCountryVisited = int(tldr.TotalCountryVisited)
	summary.TotalPageURL = int(tldr.TotalPageUrl)

	// event summary detail
	sum, err := s.Repo.GetEventDetailSummary(ctx, gen.GetEventDetailSummaryParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	var castQueryType = map[string]*[]entities.EventTextTotal{
		"most_visited_url":     &summary.MostVisitedURLs,
		"most_visited_country": &summary.MostCountryVisited,
		"most_visited_city":    &summary.MostCitiesVisited,
		"most_used_browser":    &summary.MostUsedBrowsers,
		"most_hit_element":     &summary.MostElementsFired,
		"most_event_type":      &summary.MostFiredEventType,
		"most_event_label":     &summary.MostFiredEventLabel,
	}

	for _, v := range sum {
		switch v.QueryType {
		// case for last visited user since the "total" field need to
		// be converted into time.Time
		case "last_visited_user":
			timestamp, _ := time.Parse("2006-01-02 15:04:05.999999-07", v.Total)
			ip, _, _ := net.ParseCIDR(v.Name.String)
			summary.LastVisitedUsers = append(summary.LastVisitedUsers, entities.EventLastUser{
				IP:        ip,
				Timestamp: timestamp,
			})
		// otherwise, just cast them directly into a map which contains
		// the pointer into respective fields.
		default:
			total, err := strconv.Atoi(v.Total)
			if err != nil {
				total = 0
			}

			// cast query type into corresponding type in the map
			if slice, ok := castQueryType[v.QueryType]; ok {
				*slice = append(*slice, entities.EventTextTotal{
					Name:  v.Name.String,
					Total: total,
				})
			}
		}
	}

	return &summary, nil
}

func (s *EventServiceImpl) GetWeeklyEventsChart(ctx context.Context, projectID string, userID string) (*entities.EventSummaryChart, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)

	events, err := s.Repo.GetWeeklyEvents(ctx, gen.GetWeeklyEventsParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	totalWeekly, err := s.Repo.GetWeeklyEventsTotal(ctx, gen.GetWeeklyEventsTotalParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	summary := entities.EventSummaryChart{
		Total: int(totalWeekly),
		Time:  events,
	}

	return &summary, nil
}

func (s *EventServiceImpl) GetEventTypeChart(ctx context.Context, projectID string, userID string) ([]gen.GetPercentageEventsTypeRow, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)
	return s.Repo.GetPercentageEventsType(ctx, gen.GetPercentageEventsTypeParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
}

func (s *EventServiceImpl) GetEventLabelChart(ctx context.Context, projectID string, userID string) ([]gen.GetPercentageEventsLabelRow, error) {
	projectUUID, userUUID := uuid.MustParse(projectID), uuid.MustParse(userID)
	return s.Repo.GetPercentageEventsLabel(ctx, gen.GetPercentageEventsLabelParams{
		ProjectID: projectUUID,
		UserID:    userUUID,
	})
}

func (s *EventServiceImpl) CountUserMonthlyEvents(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.Repo.CountUserMonthlyEvents(ctx, userID)
}
