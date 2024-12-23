package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubkudev/sentinel/internal/middlewares"
	"github.com/hubkudev/sentinel/internal/services"
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
	event.Get("/live", m.ProtectedRoute, m.LiveEventsCache, apiService.LiveEvents)
	event.Get("/live/:id", m.ProtectedRoute, m.LiveEventCache, apiService.LiveEventDetail)
	event.Get("/summary/detail/:id", m.ProtectedRoute, m.LiveEventDetailSummaryCache, apiService.GetEventSummaryDetail)
	event.Get("/summary/:id", m.ProtectedRoute, m.LiveEventSummaryCache, apiService.GetEventSummary)
	event.Get("/monthly/count", m.ProtectedRoute, apiService.CountMonthlyEvents)
	event.Post("/download/start", m.ProtectedRoute, apiService.StartDownloadEvent)
	event.Get("/download/finish/:id", m.ProtectedRoute, apiService.FinishDownloadEvent)

	json := api.Group("json")
	json.Get("/event/chart/:id", m.ProtectedRoute, m.JSONWeeklyEventCache, apiService.JSONWeeklyEventChart)
	json.Get("/event-type/chart/:id", m.ProtectedRoute, m.JSONEventTypeCache, apiService.JSONEventTypeChart)
	json.Get("/event-label/chart/:id", m.ProtectedRoute, m.JSONEventLabelCache, apiService.JSONEventLabelChart)

	key := api.Group("key")
	key.Post("/create", m.ProtectedRoute, apiService.CreateAPIKey)
	key.Delete("/delete", m.ProtectedRoute, apiService.DeleteAPIKey)

	v1 := api.Group("v1")
	v1.Post("/event", m.APIProtectedRoute, eventService.CreateEvent)
}
