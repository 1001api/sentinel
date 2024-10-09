package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/middlewares"
	"github.com/hubkudev/sentinel/services"
)

func InitWebRoute(app *fiber.App, m middlewares.Middleware, webService services.WebService) {
	app.Get("/", webService.SendLandingPage)
	app.Get("/login", webService.SendLoginPage)
	app.Get("/dashboard", m.ProtectedRoute, webService.SendDashboardPage)
	app.Get("/projects", m.ProtectedRoute, webService.SendProjectsPage)
	app.Get("/api-keys", m.ProtectedRoute, webService.SendAPIKeysPage)

	misc := app.Group("/misc")
	misc.Get("/tos", webService.SendTOSPage)
	misc.Get("/auth-redirect", webService.SendAuthRedirectPage)
}
