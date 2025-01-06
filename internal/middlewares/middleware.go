package middlewares

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/dto"
	"github.com/hubkudev/sentinel/internal/entities"
	"github.com/hubkudev/sentinel/internal/services"
	"github.com/hubkudev/sentinel/views/pages"
)

type Middleware interface {
	ProtectedRoute(c *fiber.Ctx) error
	InternalRoute(c *fiber.Ctx) error
	APIPublicRoute(c *fiber.Ctx) error
	APIPrivateRoute(c *fiber.Ctx) error
	UnProtectedRoute(c *fiber.Ctx) error
	LiveEventsCache(c *fiber.Ctx) error
	LiveEventCache(c *fiber.Ctx) error
	LiveEventSummaryCache(c *fiber.Ctx) error
	LiveEventDetailSummaryCache(c *fiber.Ctx) error
	JSONWeeklyEventCache(c *fiber.Ctx) error
	JSONEventTypeCache(c *fiber.Ctx) error
	JSONEventLabelCache(c *fiber.Ctx) error
}

type MiddlewareImpl struct {
	UserService    services.UserService
	CacheService   services.CacheService
	SessionStorage *session.Store
}

func InitMiddleware(
	userService services.UserService,
	sessionStore *session.Store,
	cacheService services.CacheService,
) MiddlewareImpl {
	return MiddlewareImpl{
		UserService:    userService,
		SessionStorage: sessionStore,
		CacheService:   cacheService,
	}
}

func (m *MiddlewareImpl) ProtectedRoute(c *fiber.Ctx) error {
	sess, err := m.SessionStorage.Get(c)
	if err != nil {
		log.Println(err)
		return redirectTo(c, "/login")
	}

	userID, ok := sess.Get("ID").(string)
	if !ok || userID == "" {
		return redirectTo(c, "/login")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return redirectTo(c, "/login")
	}

	// check if the user is exist in the database
	exist, err := m.UserService.FindByID(userUUID)
	if exist == nil {
		return redirectTo(c, "/login")
	}

	c.Locals("user", exist)

	return c.Next()
}

func (m *MiddlewareImpl) UnProtectedRoute(c *fiber.Ctx) error {
	sess, err := m.SessionStorage.Get(c)
	if err != nil {
		c.Locals("user", nil)
		return c.Next()
	}

	userID, ok := sess.Get("ID").(string)
	if !ok || userID == "" {
		c.Locals("user", nil)
		return c.Next()
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return redirectTo(c, "/login")
	}

	// check if the user is exist in the database
	exist, err := m.UserService.FindByID(userUUID)
	if exist == nil {
		c.Locals("user", nil)
		return c.Next()
	}

	c.Locals("user", exist)
	return c.Next()
}

func (m *MiddlewareImpl) APIPublicRoute(c *fiber.Ctx) error {
	var key dto.KeyPayload

	if err := c.BodyParser(&key); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "PublicKey field is required",
		})
	}

	if key.PublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "PublicKey field is required",
		})
	}

	// check if the key is exist in the database
	exist, err := m.UserService.FindByPublicKey(key.PublicKey)
	if exist == nil || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "invalid public key",
		})
	}

	c.Locals("user", exist)

	return c.Next()
}

func (m *MiddlewareImpl) APIPrivateRoute(c *fiber.Ctx) error {
	var key dto.PrivateKeyPayload

	if err := c.BodyParser(&key); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "PrivateKey field is required",
		})
	}

	if key.PrivateKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "PrivateKey field is required",
		})
	}

	// check if the key is valid in the database
	exist, err := m.UserService.FindByPrivateKey(key.PrivateKey)
	if exist == nil || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "invalid private key",
		})
	}

	c.Locals("user", exist)

	return c.Next()
}

func (m *MiddlewareImpl) InternalRoute(c *fiber.Ctx) error {
	// get passphrase header
	token := string(c.Request().Header.Peek("passphrase"))

	// validate token/passphrase is the same as SENTINEL_INTERNAL_PASS=
	// please set this in .env to allow internal route to be secured internally.
	internalPass := os.Getenv("SENTINEL_INTERNAL_PASS")
	if token != internalPass {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var payload dto.UserIDPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "UserID field is required",
		})
	}

	if payload.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "UserID field is required",
		})
	}

	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid UserID",
		})
	}

	// check if the key is valid in the database
	exist, err := m.UserService.FindByID(userUUID)
	if exist == nil || err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "invalid UserID",
		})
	}

	c.Locals("user", exist)

	return c.Next()
}

func (m *MiddlewareImpl) LiveEventsCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_LIVE_EVENTS(user.ID))
	if err != nil {
		log.Println("error getting cached live events:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var events []gen.GetLiveEventsRow
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&events); err != nil {
			log.Println("error decoding cached live events:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		buf := bytes.Buffer{}
		eventRows := pages.EventLiveTableRow(events)
		eventRows.Render(context.Background(), &buf)

		return c.SendString(buf.String())
	}

	return c.Next()
}

func (m *MiddlewareImpl) LiveEventCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	strategy := c.Query("strategy", "last_hour")

	cached, err := m.CacheService.GetCache(configs.CACHE_LIVE_EVENT(user.ID, projectUUID, strategy))
	if err != nil {
		log.Println("error getting cached live event:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var events []gen.GetLiveEventsDetailRow
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&events); err != nil {
			log.Println("error decoding cached live events:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		buf := bytes.Buffer{}
		eventRows := pages.EventDetailTableRow(events)
		eventRows.Render(context.Background(), &buf)

		return c.SendString(buf.String())
	}

	return c.Next()
}

func (m *MiddlewareImpl) LiveEventSummaryCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_LIVE_EVENT_SUMMARY(user.ID, projectUUID))
	if err != nil {
		log.Println("error getting cached live event summary:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var summary gen.GetEventSummaryRow
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&summary); err != nil {
			log.Println("error decoding cached live events:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		buf := bytes.Buffer{}
		eventRows := pages.ProjectSummaryText(&summary)
		eventRows.Render(context.Background(), &buf)

		return c.SendString(buf.String())
	}

	return c.Next()
}

func (m *MiddlewareImpl) LiveEventDetailSummaryCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_LIVE_EVENT_DETAIL_SUMMARY(user.ID, projectUUID))
	if err != nil {
		log.Println("error getting cached live event detail summary:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var summary entities.EventDetail
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&summary); err != nil {
			log.Println("error decoding cached live event detail summary:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		buf := bytes.Buffer{}
		eventRows := pages.EventDetailSummarySection(&summary)
		eventRows.Render(context.Background(), &buf)

		return c.SendString(buf.String())
	}

	return c.Next()
}

func (m *MiddlewareImpl) JSONWeeklyEventCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_JSON_WEEKLY_EVENT_CHART(user.ID, projectUUID))
	if err != nil {
		log.Println("error getting json weekly event cache:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var weeklyEvent entities.EventSummaryChart
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&weeklyEvent); err != nil {
			log.Println("error decoding cached json weekly event cache:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(weeklyEvent)
	}

	return c.Next()
}

func (m *MiddlewareImpl) JSONEventTypeCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_JSON_EVENT_TYPE_CHART(user.ID, projectUUID))
	if err != nil {
		log.Println("error getting json event type cache:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var types []gen.GetPercentageEventsTypeRow
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&types); err != nil {
			log.Println("error decoding cached json event type:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(types)
	}

	return c.Next()
}

func (m *MiddlewareImpl) JSONEventLabelCache(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*gen.FindUserByIDRow)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString("ProjectID is required")
	}

	cached, err := m.CacheService.GetCache(configs.CACHE_JSON_EVENT_LABEL_CHART(user.ID, projectUUID))
	if err != nil {
		log.Println("error getting json event label cache:", err)
	}

	if len(cached) > 0 {
		// deserialize gob object from redis
		var labels []gen.GetPercentageEventsLabelRow
		if err := gob.NewDecoder(bytes.NewReader(cached)).Decode(&labels); err != nil {
			log.Println("error decoding cached json event label:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(labels)
	}

	return c.Next()
}

func redirectTo(c *fiber.Ctx, path string) error {
	htmx := c.Get("HX-Request")

	// if request is not coming from HTMX, redirect normally
	if htmx == "" && htmx != "true" {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect(path)
	}

	// if request is coming from htmx, redirect using htmx way
	c.Set("HX-Redirect", path)
	return c.SendStatus(fiber.StatusOK)
}
