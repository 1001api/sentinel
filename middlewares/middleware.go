package middlewares

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hubkudev/sentinel/dto"
	repositories "github.com/hubkudev/sentinel/repos"
)

type Middleware interface {
	ProtectedRoute(c *fiber.Ctx) error
	APIProtectedRoute(c *fiber.Ctx) error
}

type MiddlewareImpl struct {
	UserRepo       repositories.UserRepository
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
	exist, err := m.UserRepo.FindByID(context.Background(), userID.(string))
	if exist == nil {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/login")
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
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// check if the key is exist in the database
	exist, err := m.UserRepo.FindByPublicKey(context.Background(), key.PublicKey)
	if exist == nil || err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user", exist)

	return c.Next()
}
