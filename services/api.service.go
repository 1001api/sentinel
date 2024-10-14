package services

import (
	"bytes"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/entities"
	"github.com/hubkudev/sentinel/views/pages"
)

type APIService interface {
	CreateProject(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
	LiveEvents(ctx *fiber.Ctx) error
	LiveEventDetail(ctx *fiber.Ctx) error
	GetEventSummary(c *fiber.Ctx) error
	GetEventSummaryDetail(c *fiber.Ctx) error
}

type APIServiceImpl struct {
	ProjectService ProjectService
	EventService   EventService
}

func (s *APIServiceImpl) CreateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	name := c.FormValue("project_name")
	desc := c.FormValue("project_desc")

	if name == "" {
		return c.SendString("Project name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	if len(desc) > 200 {
		return c.SendString("Maximum description length is 200 characters")
	}

	_, err := s.ProjectService.CreateProject(name, desc, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) UpdateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	name := c.FormValue("project_name")
	desc := c.FormValue("project_desc")
	projectID := c.FormValue("project_id")

	if name == "" {
		return c.SendString("Project name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	if len(desc) > 200 {
		return c.SendString("Maximum description length is 200 characters")
	}

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.UpdateProject(name, desc, projectID, user.ID); err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	projectID := c.FormValue("project_id")

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.DeleteProject(user.ID, projectID); err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) LiveEvents(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)

	events, err := s.EventService.GetLiveEvents(context.Background(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventLiveTableRow(events)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) LiveEventDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	events, err := s.EventService.GetLiveEventDetail(context.Background(), projectID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventDetailTableRow(events)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) GetEventSummary(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventSummary(context.Background(), projectID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.ProjectSummaryText(summary)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) GetEventSummaryDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventDetailSummary(context.Background(), projectID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventDetailSummarySection(summary)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}
