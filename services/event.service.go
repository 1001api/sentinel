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
