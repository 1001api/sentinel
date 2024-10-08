package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/services"
)

func InitAuthRoute(app *fiber.App, authService services.AuthService) {
	auth := app.Group("auth")

	auth.Get("/google", authService.GoogleLogin)
	auth.Get("/google/callback", authService.GoogleCallback)
	auth.Get("/logout", authService.Logout)
}
