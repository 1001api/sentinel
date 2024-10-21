package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/middlewares"
	"github.com/hubkudev/sentinel/services"
)

func InitAPIRoute(app *fiber.App, m middlewares.Middleware, apiService services.APIService, eventService services.EventService) {
	api := app.Group("api")

	project := api.Group("project")
	project.Post("/create", m.ProtectedRoute, apiService.CreateProject)
	project.Put("/update", m.ProtectedRoute, apiService.UpdateProject)
	project.Delete("/delete", m.ProtectedRoute, apiService.DeleteProject)
	project.Get("/size/:id", m.ProtectedRoute, apiService.CountProjectSize)
	project.Get("/last-data-retrieved/:id", m.ProtectedRoute, apiService.LastDataRetrieved)

	event := api.Group("event")
	event.Get("/live", m.ProtectedRoute, apiService.LiveEvents)
	event.Get("/live/:id", m.ProtectedRoute, apiService.LiveEventDetail)
	event.Get("/summary/detail/:id", m.ProtectedRoute, apiService.GetEventSummaryDetail)
	event.Get("/summary/:id", m.ProtectedRoute, apiService.GetEventSummary)
	event.Get("/monthly/count", m.ProtectedRoute, apiService.CountMonthlyEvents)

	json := api.Group("json")
	json.Get("/event/chart/:id", m.ProtectedRoute, apiService.JSONWeeklyEventChart)
	json.Get("/event-type/chart/:id", m.ProtectedRoute, apiService.JSONEventTypeChart)
	json.Get("/event-label/chart/:id", m.ProtectedRoute, apiService.JSONEventLabelChart)

	v1 := api.Group("v1")
	v1.Post("/event", m.APIProtectedRoute, eventService.CreateEvent)
}
