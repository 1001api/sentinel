package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/views/pages"
)

type APIService interface {
	CreateProject(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
	CountProjectSize(ctx *fiber.Ctx) error
	CountMonthlyEvents(ctx *fiber.Ctx) error
	LastDataRetrieved(c *fiber.Ctx) error
	LiveEvents(ctx *fiber.Ctx) error
	LiveEventDetail(ctx *fiber.Ctx) error
	GetEventSummary(c *fiber.Ctx) error
	GetEventSummaryDetail(c *fiber.Ctx) error
	StartDownloadEvent(c *fiber.Ctx) error
	FinishDownloadEvent(c *fiber.Ctx) error
	JSONWeeklyEventChart(c *fiber.Ctx) error
	JSONEventTypeChart(c *fiber.Ctx) error
	JSONEventLabelChart(c *fiber.Ctx) error
}

type APIServiceImpl struct {
	ProjectService  ProjectService
	EventService    EventService
	DownloadService DownloadService
}

func (s *APIServiceImpl) CreateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	name := c.FormValue("project_name")
	desc := c.FormValue("project_desc")

	if name == "" {
		return c.SendString("Project name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	if len(desc) > 200 {
		return c.SendString("Maximum description length is 200 characters")
	}

	_, err := s.ProjectService.CreateProject(context.Background(), name, desc, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) UpdateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	name := c.FormValue("project_name")
	desc := c.FormValue("project_desc")
	projectID := c.FormValue("project_id")

	if name == "" {
		return c.SendString("Project name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	if len(desc) > 200 {
		return c.SendString("Maximum description length is 200 characters")
	}

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.UpdateProject(context.Background(), name, desc, projectID, user.ID.String()); err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.FormValue("project_id")

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.DeleteProject(context.Background(), user.ID.String(), projectID); err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) CountProjectSize(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	sizeInKB, err := s.ProjectService.CountProjectSize(context.Background(), projectID, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	text := pages.ProjectSizeText(sizeInKB)
	text.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) LastDataRetrieved(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	lastTime, err := s.ProjectService.LastProjectDataReceived(context.Background(), projectID, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("N/A")
	}

	buf := bytes.Buffer{}

	text := pages.ProjectLastDataText(lastTime)
	text.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) StartDownloadEvent(c *fiber.Ctx) error {
	projectID := c.FormValue("id")
	if projectID == "" {
		return c.SendString("Project ID required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	c.Set("HX-Redirect", fmt.Sprintf("/api/event/download/finish/%s", projectUUID.String()))
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) FinishDownloadEvent(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.SendString("Project ID required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	// check if project id is available within user context
	exist, err := s.ProjectService.GetProjectByID(context.Background(), projectUUID.String(), user.ID.String())
	if exist == nil || err != nil {
		return c.SendString("Project not found")
	}

	filename := fmt.Sprintf("%s.csv", time.Now().Format("02/01/2006_15:04:05"))
	tempFile, err := os.CreateTemp("temp", "temp-*.csv")
	if err != nil {
		log.Println("Error creating temp file:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	header, body, err := s.DownloadService.DownloadEventData(context.Background(), projectUUID, user.ID)
	if err != nil {
		log.Println("Error retrieving Event data:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}

	if err := writer.Write(header); err != nil {
		log.Println("Error writing CSV header:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}

	for _, row := range body {
		if err := writer.Write(row); err != nil {
			log.Println("Error writing CSV row:", err)
			return c.SendString("Unable to download your file, please try again later.")
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Println("Error flushing CSV:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}

	return c.Download(tempFile.Name(), filename)
}

func (s *APIServiceImpl) LiveEvents(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	events, err := s.EventService.GetLiveEvents(context.Background(), user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventLiveTableRow(events)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) LiveEventDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	events, err := s.EventService.GetLiveEventDetail(context.Background(), projectID, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventDetailTableRow(events)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) GetEventSummary(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventSummary(context.Background(), projectID, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.ProjectSummaryText(summary)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) GetEventSummaryDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventDetailSummary(context.Background(), projectID, user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	buf := bytes.Buffer{}

	eventRows := pages.EventDetailSummarySection(summary)
	eventRows.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) JSONWeeklyEventChart(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	weeklyEvents, err := s.EventService.GetWeeklyEventsChart(context.Background(), projectID, user.ID.String())
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(weeklyEvents)
}

func (s *APIServiceImpl) JSONEventTypeChart(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	percentage, err := s.EventService.GetEventTypeChart(context.Background(), projectID, user.ID.String())
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(percentage)
}

func (s *APIServiceImpl) JSONEventLabelChart(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	labels, err := s.EventService.GetEventLabelChart(context.Background(), projectID, user.ID.String())
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(labels)
}

func (s *APIServiceImpl) CountMonthlyEvents(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	monthlyEvents, err := s.EventService.CountUserMonthlyEvents(context.Background(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	return c.SendString(fmt.Sprintf("%d", monthlyEvents))
}
