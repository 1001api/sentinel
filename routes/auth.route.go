package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/services"
)

func InitAuthRoute(app *fiber.App, authService services.AuthService) {
	auth := app.Group("auth")

	auth.Post("/register-first-time", authService.RegisterFirstUser)
	auth.Post("/login", authService.Login)
	auth.Get("/logout", authService.Logout)
}
