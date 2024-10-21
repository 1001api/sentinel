package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/services"
)

type Middleware interface {
	ProtectedRoute(c *fiber.Ctx) error
	APIProtectedRoute(c *fiber.Ctx) error
	UnProtectedRoute(c *fiber.Ctx) error
}

type MiddlewareImpl struct {
	UserService    services.UserService
	SessionStorage *session.Store
}

func (m *MiddlewareImpl) ProtectedRoute(c *fiber.Ctx) error {
	sess, err := m.SessionStorage.Get(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/login")
	}

	userID := sess.Get("ID")
	if userID == nil {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/login")
	}

	// check if the user is exist in the database
	exist, err := m.UserService.FindByID(userID.(string))
	if exist == nil {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/login")
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

	userID := sess.Get("ID")
	if userID == nil {
		c.Locals("user", nil)
		return c.Next()
	}

	// check if the user is exist in the database
	exist, err := m.UserService.FindByID(userID.(string))
	if exist == nil {
		c.Locals("user", nil)
		return c.Next()
	}

	c.Locals("user", exist)
	return c.Next()
}

func (m *MiddlewareImpl) APIProtectedRoute(c *fiber.Ctx) error {
	var key dto.KeyPayload

	if err := c.BodyParser(&key); err != nil {
		// send raw error (unprocessable entity)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if key.PublicKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "valid PublicKey is required",
		})
	}

	// check if the key is exist in the database
	exist, err := m.UserService.FindByPublicKey(key.PublicKey)
	if exist == nil || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "valid PublicKey is required",
		})
	}

	c.Locals("user", exist)

	return c.Next()
}
