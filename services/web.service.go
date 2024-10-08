package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/views/pages"
)

type WebService interface {
	SendLandingPage(ctx *fiber.Ctx) error
	SendDashboardPage(ctx *fiber.Ctx) error
}

type WebServiceImpl struct{}

func (s *WebServiceImpl) SendLandingPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.IndexPage())
}

func (s *WebServiceImpl) SendDashboardPage(c *fiber.Ctx) error {
	return configs.Render(c, pages.DashboardPage())
}
