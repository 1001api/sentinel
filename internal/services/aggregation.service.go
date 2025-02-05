package services

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/entities"
	"github.com/hubkudev/sentinel/internal/repositories"
)

type AggrService interface {
	GetBriefSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.GetBriefAggrRow, error)
	GetDetailSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*entities.EventDetail, error)
	SaveProjectAggr(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	FindProjectAggr(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, limit int32) ([]entities.ProjectAggr, error)
}

type AggrServiceImpl struct {
	UtilService UtilService
	AggrRepo    repositories.AggrRepo
	ProjectRepo repositories.ProjectRepo
}

func InitAggrService(
	utilService UtilService,
	aggrRepo repositories.AggrRepo,
	projectRepo repositories.ProjectRepo,
) AggrServiceImpl {
	return AggrServiceImpl{
		UtilService: utilService,
		AggrRepo:    aggrRepo,
		ProjectRepo: projectRepo,
	}
}

func (s *AggrServiceImpl) GetBriefSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.GetBriefAggrRow, error) {
	row, err := s.AggrRepo.GetBriefAggr(ctx, &gen.GetBriefAggrParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (s *AggrServiceImpl) GetDetailSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*entities.EventDetail, error) {
	var summary entities.EventDetail

	// event summary total numbering
	tldr, err := s.AggrRepo.GetTotalAggr(ctx, &gen.GetTotalAggrParams{
		ProjectID: projectID,
		UserID:    userID,
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
	sum, err := s.AggrRepo.GetDetailAggr(ctx, &gen.GetDetailAggrParams{
		ProjectID: projectID,
		UserID:    userID,
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

func (s *AggrServiceImpl) SaveProjectAggr(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	summary, err := s.GetDetailSummary(ctx, projectID, userID)
	if err != nil {
		return err
	}

	urlJson, err := json.Marshal(summary.MostVisitedURLs)
	if err != nil {
		return err
	}

	citiesJson, err := json.Marshal(summary.MostCitiesVisited)
	if err != nil {
		return err
	}

	countryJson, err := json.Marshal(summary.MostCountryVisited)
	if err != nil {
		return err
	}

	elementsJson, err := json.Marshal(summary.MostElementsFired)
	if err != nil {
		return err
	}

	lastUsersJson, err := json.Marshal(summary.LastVisitedUsers)
	if err != nil {
		return err
	}

	browsersJson, err := json.Marshal(summary.MostUsedBrowsers)
	if err != nil {
		return err
	}

	eventTypesJson, err := json.Marshal(summary.MostFiredEventType)
	if err != nil {
		return err
	}

	eventLabelsJson, err := json.Marshal(summary.MostFiredEventLabel)
	if err != nil {
		return err
	}

	now := time.Now()

	input := gen.CreateProjectAggrParams{
		ProjectID:            projectID,
		UserID:               userID,
		TotalEvents:          int32(summary.TotalEvents),
		TotalEventTypes:      int32(summary.TotalEventType),
		TotalUniqueUsers:     int32(summary.TotalUniqueUsers),
		TotalLocations:       int32(summary.TotalCountryVisited),
		TotalUniquePageUrls:  int32(summary.TotalPageURL),
		MostVisitedUrls:      urlJson,
		MostVisitedCountries: citiesJson,
		MostVisitedCities:    citiesJson,
		MostVisitedRegions:   countryJson,
		MostFiringElements:   elementsJson,
		LastVisitedUsers:     lastUsersJson,
		MostUsedBrowsers:     browsersJson,
		MostFiredEventTypes:  eventTypesJson,
		MostFiredEventLabels: eventLabelsJson,
		AggregatedAt:         now,
		AggregatedAtStr:      now.Format("02/01/2006"),
	}

	if err := s.ProjectRepo.CreateProjectAggr(ctx, &input); err != nil {
		return err
	}

	return nil
}

func (s *AggrServiceImpl) FindProjectAggr(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, limit int32) ([]entities.ProjectAggr, error) {
	summary, err := s.ProjectRepo.FindProjectAggr(ctx, &gen.FindProjectAggrParams{
		ProjectID: projectID,
		UserID:    userID,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	var results []entities.ProjectAggr

	for _, v := range summary {
		var result entities.ProjectAggr
		result.TotalEvents = v.TotalEvents
		result.TotalEventTypes = v.TotalEventTypes
		result.TotalUniqueUsers = v.TotalUniqueUsers
		result.TotalLocations = v.TotalLocations
		result.TotalUniquePageUrls = v.TotalUniquePageUrls
		result.MostVisitedUrls = s.UtilService.ByteToJSON(v.MostVisitedUrls)
		result.MostVisitedCountries = s.UtilService.ByteToJSON(v.MostVisitedCountries)
		result.MostVisitedCities = s.UtilService.ByteToJSON(v.MostVisitedCities)
		result.MostVisitedRegions = s.UtilService.ByteToJSON(v.MostVisitedRegions)
		result.MostFiringElements = s.UtilService.ByteToJSON(v.MostFiringElements)
		result.LastVisitedUsers = s.UtilService.ByteToJSON(v.LastVisitedUsers)
		result.MostUsedBrowsers = s.UtilService.ByteToJSON(v.MostUsedBrowsers)
		result.MostFiredEventTypes = s.UtilService.ByteToJSON(v.MostFiredEventTypes)
		result.MostFiredEventLabels = s.UtilService.ByteToJSON(v.MostFiredEventLabels)
		result.AggregatedAtStr = v.AggregatedAtStr
		results = append(results, result)
	}

	return results, err
}
