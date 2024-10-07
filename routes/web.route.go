package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/views/pages"
)

func InitWebRoute(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return configs.Render(c, pages.IndexPage())
	})
}
