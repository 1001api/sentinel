package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/entities"
	"github.com/hubkudev/sentinel/views/pages"
	"github.com/hubkudev/sentinel/views/pages/misc"
)

type WebService interface {
	SendLandingPage(ctx *fiber.Ctx) error
	SendLoginPage(ctx *fiber.Ctx) error
	SendDashboardPage(ctx *fiber.Ctx) error
	SendProjectsPage(ctx *fiber.Ctx) error
	SendAPIKeysPage(ctx *fiber.Ctx) error
	SendTOSPage(ctx *fiber.Ctx) error
	SendAuthRedirectPage(ctx *fiber.Ctx) error
}

type WebServiceImpl struct {
	UserService    UserService
	ProjectService ProjectService
}

func (s *WebServiceImpl) SendLandingPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.IndexPage())
}

func (s *WebServiceImpl) SendLoginPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.LoginPage())
}

func (s *WebServiceImpl) SendDashboardPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)
	return configs.Render(c, pages.DashboardPage(user))
}

func (s *WebServiceImpl) SendProjectsPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)

	projects, err := s.ProjectService.GetAllProjects(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.ProjectsPage(user, projects))
}

func (s *WebServiceImpl) SendAPIKeysPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*entities.User)

	publicKey, err := s.UserService.GetPublicKey(user.ID)
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
