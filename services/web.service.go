package services

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/configs"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/views/pages"
	"github.com/hubkudev/sentinel/views/pages/misc"
)

type WebService interface {
	SendLandingPage(ctx *fiber.Ctx) error
	SendLoginPage(ctx *fiber.Ctx) error
	SendDashboardPage(ctx *fiber.Ctx) error
	SendEventsPage(ctx *fiber.Ctx) error
	SendEventDetailPage(ctx *fiber.Ctx) error
	SendProjectsPage(ctx *fiber.Ctx) error
	SendAPIKeysPage(ctx *fiber.Ctx) error
	SendTOSPage(ctx *fiber.Ctx) error
	SendAuthRedirectPage(ctx *fiber.Ctx) error
}

type WebServiceImpl struct {
	UserService    UserService
	ProjectService ProjectService
	EventService   EventService
}

func (s *WebServiceImpl) SendLandingPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.IndexPage())
}

func (s *WebServiceImpl) SendLoginPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.LoginPage())
}

func (s *WebServiceImpl) SendDashboardPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	return configs.Render(c, pages.DashboardPage(user))
}

func (s *WebServiceImpl) SendEventsPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	events, err := s.EventService.GetLiveEvents(context.Background(), user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	projects, err := s.ProjectService.GetAllProjects(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return configs.Render(c, pages.EventsPage(pages.EventsPageProps{
		User:     user,
		Events:   events,
		Projects: projects,
	}))
}

func (s *WebServiceImpl) SendEventDetailPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ProjectID is required"})
	}

	project, err := s.ProjectService.GetProjectByID(projectID, user.ID.String())
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	summary, err := s.EventService.GetEventDetailSummary(context.Background(), projectID, user.ID.String())
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	weeklyEvents, err := s.EventService.GetWeeklyEventsChart(context.Background(), projectID, user.ID.String())
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.EventDetailPage(pages.EventDetailPageProps{
		User:             user,
		Project:          project,
		Summary:          summary,
		WeeklyEventChart: weeklyEvents,
	}))
}

func (s *WebServiceImpl) SendProjectsPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	projects, err := s.ProjectService.GetAllProjects(user.ID.String())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.ProjectsPage(user, projects))
}

func (s *WebServiceImpl) SendAPIKeysPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	publicKey, err := s.UserService.GetPublicKey(user.ID.String())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.APIKeysPage(user, publicKey))
}

func (s *WebServiceImpl) SendTOSPage(c *fiber.Ctx) error {
	return configs.Render(c, misc.TOSPage())
}

func (s *WebServiceImpl) SendAuthRedirectPage(c *fiber.Ctx) error {
	return configs.Render(c, misc.AuthSuccessPage())
}
