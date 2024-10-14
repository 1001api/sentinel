package services

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	repositories "github.com/hubkudev/sentinel/repos"
)

type EventService interface {
	CreateEvent(c *fiber.Ctx) error
	GetLiveEvents(ctx context.Context, userID string) ([]entities.Event, error)
	GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]entities.Event, error)
	GetEventSummary(ctx context.Context, projectID string, userID string) (*entities.EventSummary, error)
	GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error)
	GetWeeklyEventsChart(ctx context.Context, projectID string, userID string) (*entities.EventSummaryChart, error)
}

type EventServiceImpl struct {
	UtilService UtilService
	ProjectRepo repositories.ProjectRepository
	EventRepo   repositories.EventRepository
}

func (s *EventServiceImpl) CreateEvent(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
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
	exist, _ := s.ProjectRepo.CheckWithinUserID(context.Background(), input.ProjectID, user.ID)
	if !exist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project not found"})
	}

	if err := s.EventRepo.CreateEvent(context.Background(), &input, user.ID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *EventServiceImpl) GetLiveEvents(ctx context.Context, userID string) ([]entities.Event, error) {
	return s.EventRepo.GetLiveEvents(ctx, userID)
}

func (s *EventServiceImpl) GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]entities.Event, error) {
	return s.EventRepo.GetLiveEventDetail(ctx, projectID, userID)
}

func (s *EventServiceImpl) GetEventSummary(ctx context.Context, projectID string, userID string) (*entities.EventSummary, error) {
	return s.EventRepo.GetEventSummary(ctx, projectID, userID)
}

func (s *EventServiceImpl) GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error) {
	return s.EventRepo.GetEventDetailSummary(ctx, projectID, userID)
}

func (s *EventServiceImpl) GetWeeklyEventsChart(ctx context.Context, projectID string, userID string) (*entities.EventSummaryChart, error) {
	events, err := s.EventRepo.GetWeeklyEventsTime(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	totalWeekly, err := s.EventRepo.GetWeeklyEventsTotal(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	summary := entities.EventSummaryChart{
		Total: totalWeekly,
		Time:  events,
	}

	return &summary, nil
}
