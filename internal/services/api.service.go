package services

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/configs"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/constants"
	"github.com/hubkudev/sentinel/internal/entities"
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
	CreateAPIKey(ctx *fiber.Ctx) error
	DeleteAPIKey(ctx *fiber.Ctx) error
}

type APIServiceImpl struct {
	ProjectService  ProjectService
	EventService    EventService
	DownloadService DownloadService
	CacheService    CacheService
	KeyService      KeyService
}

func InitAPIService(
	projectService ProjectService,
	eventService EventService,
	downloadService DownloadService,
	cacheService CacheService,
	keyService KeyService,
) APIServiceImpl {
	return APIServiceImpl{
		ProjectService:  projectService,
		EventService:    eventService,
		DownloadService: downloadService,
		CacheService:    cacheService,
		KeyService:      keyService,
	}
}

func (s *APIServiceImpl) CreateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	name := c.FormValue("project_name")
	desc := c.FormValue("project_desc")
	url := c.FormValue("project_url")

	if name == "" {
		return c.SendString("Project name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	if len(desc) > 200 {
		return c.SendString("Maximum description length is 200 characters")
	}

	if len(url) > 255 {
		return c.SendString("Maximum length of url is 255 characters")
	}

	_, err := s.ProjectService.CreateProject(context.Background(), name, desc, url, user.ID)
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
	url := c.FormValue("project_url")
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

	if len(url) > 255 {
		return c.SendString("Maximum length of url is 255 characters")
	}

	if projectID == "" {
		return c.SendString("Project ID required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.UpdateProject(context.Background(), name, desc, url, projectUUID, user.ID); err != nil {
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	if err := s.ProjectService.DeleteProject(context.Background(), user.ID, projectUUID); err != nil {
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	sizeInKB, err := s.ProjectService.CountProjectSize(context.Background(), projectUUID, user.ID)
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	lastTime, err := s.ProjectService.LastProjectDataReceived(context.Background(), projectUUID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("N/A")
	}

	buf := bytes.Buffer{}

	text := pages.ProjectLastDataText(lastTime)
	text.Render(context.Background(), &buf)

	return c.SendString(buf.String())
}

func (s *APIServiceImpl) StartDownloadEvent(c *fiber.Ctx) error {
	downloadFormat := c.FormValue("format", "csv")
	downloadInterval := c.FormValue("interval", "csv")
	projectID := c.FormValue("id")
	if projectID == "" {
		return c.SendString("Project ID required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	c.Set(
		"HX-Redirect",
		fmt.Sprintf("/api/event/download/finish/%s?format=%s&interval=%s", projectUUID.String(), downloadFormat, downloadInterval),
	)
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) FinishDownloadEvent(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.SendString("Project ID required")
	}

	downloadFormat := c.Query("format", "csv")
	downloadInterval := c.Query("interval", "last_7_days")

	enabled, ok := constants.DownloadFormats[downloadFormat]
	if !enabled || !ok {
		return c.SendString("Download format is not available")
	}

	intervalInt, ok := constants.Intervals[downloadInterval]
	if !ok {
		return c.SendString("Download interval is not available")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.SendString("Project ID required")
	}

	// check if project id is available within user context
	exist, err := s.ProjectService.GetProjectByID(context.Background(), projectUUID, user.ID)
	if exist == nil || err != nil {
		return c.SendString("Project not found")
	}

	// ensure directory exist, if exist, it will be no op
	if err := os.MkdirAll("temp", os.ModePerm); err != nil {
		log.Println("Error creating temp folder:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}

	filename := fmt.Sprintf("%s_%s.%s", time.Now().Format("02/01/2006_15:04:05"), downloadInterval, downloadFormat)
	tempFile, err := os.CreateTemp("temp", fmt.Sprintf("temp-*.%s", downloadFormat))
	if err != nil {
		log.Println("Error creating temp file:", err)
		return c.SendString("Unable to download your file, please try again later.")
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	header, body, err := s.DownloadService.DownloadEventData(context.Background(), projectUUID, user.ID, int32(intervalInt))
	if err != nil {
		return err
	}

	var writeErr error

	switch downloadFormat {
	case "csv":
		// write CSV
		writeErr = s.DownloadService.WriteCSV(tempFile, header, body)
	case "xlsx":
		// write XLSX
		writeErr = s.DownloadService.WriteXLSX(tempFile, header, body)
	case "json":
		// write JSON
		writeErr = s.DownloadService.WriteJSON(tempFile, header, body)
	case "html":
		// write HTML
		writeErr = s.DownloadService.WriteHTML(tempFile, header, body)
	case "pdf":
		// write PDF
		writeErr = s.DownloadService.WritePDF(tempFile, header, body)
	}

	if writeErr != nil {
		log.Println(err)
		return c.SendString("Unable to download your file, please try again later.")
	}

	return c.Download(tempFile.Name(), filename)
}

func (s *APIServiceImpl) LiveEvents(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	events, err := s.EventService.GetLiveEvents(context.Background(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	if len(events) > 0 {
		log.Printf("Caching Live Events [%s]", user.ID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(events); err != nil {
			log.Println("Error encoding events cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(configs.CACHE_LIVE_EVENTS(user.ID), cacheBuf.Bytes(), 30*time.Minute); err != nil {
			log.Println("Error saving live events as a cache:", err)
		}
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	var events []gen.GetLiveEventsDetailRow
	strategy := c.Query("strategy", entities.EventsByLastHour)

	switch strategy {
	case entities.EventsByLastHour:
		events, err = s.EventService.GetLiveEventDetail(context.Background(), projectUUID, user.ID, entities.EventsByLastHour, 100)
	case entities.EventsByLastN:
		events, err = s.EventService.GetLiveEventDetail(context.Background(), projectUUID, user.ID, entities.EventsByLastN, 100)
	}

	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	if len(events) > 0 {
		log.Printf("Caching Live Event Detail [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(events); err != nil {
			log.Println("Error encoding live event detail cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(configs.CACHE_LIVE_EVENT(user.ID, projectUUID, strategy), cacheBuf.Bytes(), 30*time.Minute); err != nil {
			log.Println("Error saving live event detail as a cache:", err)
		}
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventSummary(context.Background(), projectUUID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	if summary != nil {
		log.Printf("Caching Event Summary [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(summary); err != nil {
			log.Println("Error encoding event summary cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(configs.CACHE_LIVE_EVENT_SUMMARY(user.ID, projectUUID), cacheBuf.Bytes(), 30*time.Minute); err != nil {
			log.Println("Error saving live event summary as a cache:", err)
		}
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	summary, err := s.EventService.GetEventDetailSummary(context.Background(), projectUUID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	if summary != nil {
		log.Printf("Caching Live Event Detail Summary [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(summary); err != nil {
			log.Println("Error encoding live event detail summary cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(
			configs.CACHE_LIVE_EVENT_DETAIL_SUMMARY(user.ID, projectUUID),
			cacheBuf.Bytes(),
			30*time.Minute,
		); err != nil {
			log.Println("Error saving live event detail as a cache:", err)
		}
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

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	weeklyEvents, err := s.EventService.GetWeeklyEventsChart(context.Background(), projectUUID, user.ID)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if weeklyEvents != nil {
		log.Printf("Caching JSON Weekly Event Chart [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(weeklyEvents); err != nil {
			log.Println("Error encoding JSON weekly event chart cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(
			configs.CACHE_JSON_WEEKLY_EVENT_CHART(user.ID, projectUUID),
			cacheBuf.Bytes(),
			30*time.Minute,
		); err != nil {
			log.Println("Error saving JSON weekly event chart as a cache:", err)
		}
	}

	return c.JSON(weeklyEvents)
}

func (s *APIServiceImpl) JSONEventTypeChart(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	percentage, err := s.EventService.GetEventTypeChart(context.Background(), projectUUID, user.ID)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(percentage) > 0 {
		log.Printf("Caching JSON Event Type Chart [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(percentage); err != nil {
			log.Println("Error encoding JSON event type chart cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(
			configs.CACHE_JSON_EVENT_TYPE_CHART(user.ID, projectUUID),
			cacheBuf.Bytes(),
			30*time.Minute,
		); err != nil {
			log.Println("Error saving JSON event type chart as a cache:", err)
		}
	}

	return c.JSON(percentage)
}

func (s *APIServiceImpl) JSONEventLabelChart(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	labels, err := s.EventService.GetEventLabelChart(context.Background(), projectUUID, user.ID)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(labels) > 0 {
		log.Printf("Caching JSON Event Label Chart [%s/%s]", user.ID, projectUUID)

		// save events as a gob bytes
		var cacheBuf bytes.Buffer
		if err := gob.NewEncoder(&cacheBuf).Encode(labels); err != nil {
			log.Println("Error encoding JSON event label chart cache:", err)
		}

		// Cache events into redis for 30 minutes exp time
		// What if the new data comes?
		// The invalidation process happens at the CreateEvent,
		// after successful creation of the event.
		if err := s.CacheService.SetCache(
			configs.CACHE_JSON_EVENT_LABEL_CHART(user.ID, projectUUID),
			cacheBuf.Bytes(),
			30*time.Minute,
		); err != nil {
			log.Println("Error saving JSON event label chart as a cache:", err)
		}
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

func (s *APIServiceImpl) CreateAPIKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	name := c.FormValue("key_name")

	if name == "" {
		return c.SendString("Key name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	_, err := s.KeyService.CreateAPIKey(context.Background(), name, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) UpdateAPIKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	name := c.FormValue("key_name")

	if name == "" {
		return c.SendString("Key name is required")
	}

	if len(name) > 64 {
		return c.SendString("Maximum name length is 64 characters")
	}

	_, err := s.KeyService.CreateAPIKey(context.Background(), name, user.ID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}

func (s *APIServiceImpl) DeleteAPIKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	key := c.FormValue("key_id")

	if key == "" {
		return c.SendString("Key id is required")
	}

	keyID, err := strconv.Atoi(key)
	if err != nil {
		return c.SendString("Key id is required")
	}

	if err := s.KeyService.DeleteKey(context.Background(), user.ID, keyID); err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(fiber.StatusOK)
}
