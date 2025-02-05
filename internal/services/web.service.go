package services

import (
	"context"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/configs"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/entities"
	"github.com/hubkudev/sentinel/views/pages"
	"github.com/hubkudev/sentinel/views/pages/misc"
)

type WebService interface {
	SendLandingPage(ctx *fiber.Ctx) error
	SendLoginPage(ctx *fiber.Ctx) error
	SendCreateFirstTimeUserPage(ctx *fiber.Ctx) error
	SendDashboardPage(ctx *fiber.Ctx) error
	SendPricingPage(ctx *fiber.Ctx) error
	SendEventsPage(ctx *fiber.Ctx) error
	SendEventDetailPage(ctx *fiber.Ctx) error
	SendProjectsPage(ctx *fiber.Ctx) error
	SendAPIKeysPage(ctx *fiber.Ctx) error
	SendTOSPage(ctx *fiber.Ctx) error
	SendAuthRedirectPage(ctx *fiber.Ctx) error
}

type WebServiceImpl struct {
	UserService    UserService
	ProjectService ProjectService
	EventService   EventService
	KeyService     KeyService
	AggrService    AggrService
}

func InitWebService(
	userService UserService,
	projectService ProjectService,
	eventService EventService,
	keyService KeyService,
	aggrService AggrService,
) WebServiceImpl {
	return WebServiceImpl{
		UserService:    userService,
		ProjectService: projectService,
		EventService:   eventService,
		KeyService:     keyService,
		AggrService:    aggrService,
	}
}

func (s *WebServiceImpl) SendLandingPage(c *fiber.Ctx) error {
	user := c.Locals("user")
	var payload *gen.FindUserByIDRow
	if user != nil {
		payload = user.(*gen.FindUserByIDRow)
	}
	return configs.Render(c, pages.IndexPage(payload))
}

func (s *WebServiceImpl) SendLoginPage(c *fiber.Ctx) error {
	adminExist, err := s.UserService.CheckAdminExist()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// if admin does not exist, redirect to first time page;
	if !adminExist {
		return c.Redirect("/first-time")
	}

	return configs.Render(c, pages.LoginPage())
}

func (s *WebServiceImpl) SendCreateFirstTimeUserPage(c *fiber.Ctx) error {
	adminExist, err := s.UserService.CheckAdminExist()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// if admin exist, redirect to login page;
	// prevent from creating a new user from this page.
	if adminExist {
		return c.Redirect("/login")
	}

	return configs.Render(c, pages.CreateUserFirstTimePage())
}

func (s *WebServiceImpl) SendDashboardPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	return configs.Render(c, pages.DashboardPage(user))
}

func (s *WebServiceImpl) SendPricingPage(c *fiber.Ctx) error {
	user := c.Locals("user")
	var payload *gen.FindUserByIDRow
	if user != nil {
		payload = user.(*gen.FindUserByIDRow)
	}
	return configs.Render(c, pages.PricingPage(payload))
}

func (s *WebServiceImpl) SendEventsPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	events, err := s.EventService.GetLiveEvents(context.Background(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	projects, err := s.ProjectService.GetAllProjects(context.Background(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return configs.Render(c, pages.EventsPage(pages.EventsPageProps{
		User:     user,
		Events:   events,
		Projects: projects,
	}))
}

func (s *WebServiceImpl) SendEventDetailPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ProjectID is required"})
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	type result struct {
		data interface{}
		err  error
	}

	var wg sync.WaitGroup
	multiChan := make([]chan result, 5)
	for i := range multiChan {
		multiChan[i] = make(chan result, 1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		project, err := s.ProjectService.GetProjectByID(context.Background(), projectUUID, user.ID)
		multiChan[0] <- result{data: project, err: err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		summary, err := s.AggrService.GetDetailSummary(context.Background(), projectUUID, user.ID)
		multiChan[1] <- result{data: summary, err: err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		weeklyEvents, err := s.EventService.GetWeeklyEventsChart(context.Background(), projectUUID, user.ID)
		multiChan[2] <- result{data: weeklyEvents, err: err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		eventTypeChart, err := s.EventService.GetEventTypeChart(context.Background(), projectUUID, user.ID)
		multiChan[3] <- result{data: eventTypeChart, err: err}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		eventLabelChart, err := s.EventService.GetEventLabelChart(context.Background(), projectUUID, user.ID)
		multiChan[4] <- result{data: eventLabelChart, err: err}
	}()

	// wait to goroutines to finish
	wg.Wait()

	// close all channels
	for _, ch := range multiChan {
		close(ch)
	}

	var response pages.EventDetailPageProps

	// retrieve all data from channels
	for i, ch := range multiChan {
		res := <-ch

		// if contains err, return early
		if res.err != nil {
			log.Println(res.err)
			return c.SendStatus(fiber.StatusNotFound)
		}

		// map data based on the channel index
		switch i {
		case 0: // Project
			response.Project = res.data.(*gen.FindProjectByIDRow)
		case 1: // Summary
			response.Summary = res.data.(*entities.EventDetail)
		case 2: // WeeklyEventChart
			response.WeeklyEventChart = res.data.(*entities.EventSummaryChart)
		case 3: // EventTypeChart
			response.EventTypeChart = res.data.([]gen.GetPercentageEventsTypeRow)
		case 4: // EventLabelChart
			response.EventLabelChart = res.data.([]gen.GetPercentageEventsLabelRow)
		}
	}

	response.User = user
	return configs.Render(c, pages.EventDetailPage(response))
}

func (s *WebServiceImpl) SendProjectsPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	projects, err := s.ProjectService.GetAllProjects(context.Background(), user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.ProjectsPage(user, projects))
}

func (s *WebServiceImpl) SendAPIKeysPage(c *fiber.Ctx) error {
	user := c.Locals("user").(*gen.FindUserByIDRow)

	publicKey, err := s.UserService.GetPublicKey(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	privateKeys, err := s.KeyService.GetAllKeys(context.Background(), user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return configs.Render(c, pages.APIKeysPage(user, publicKey, privateKeys))
}

func (s *WebServiceImpl) SendTOSPage(c *fiber.Ctx) error {
	return configs.Render(c, misc.TOSPage())
}

func (s *WebServiceImpl) SendAuthRedirectPage(c *fiber.Ctx) error {
	return configs.Render(c, misc.AuthSuccessPage())
}
