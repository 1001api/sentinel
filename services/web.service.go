package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/views/pages"
	"github.com/hubkudev/sentinel/views/pages/misc"
)

type WebService interface {
	SendLandingPage(ctx *fiber.Ctx) error
	SendLoginPage(ctx *fiber.Ctx) error
	SendDashboardPage(ctx *fiber.Ctx) error
	SendTOSPage(ctx *fiber.Ctx) error
}

type WebServiceImpl struct{}

func (s *WebServiceImpl) SendLandingPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.IndexPage())
}

func (s *WebServiceImpl) SendLoginPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.LoginPage())
}

func (s *WebServiceImpl) SendDashboardPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.DashboardPage())
}

func (s *WebServiceImpl) SendTOSPage(c *fiber.Ctx) error {
	return configs.Render(c, misc.TOSPage())
}
