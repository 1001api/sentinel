package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/entities"
)

type APIService interface {
	CreateProject(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
}

type APIServiceImpl struct {
	ProjectService ProjectService
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
