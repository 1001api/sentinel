package middlewares

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	repositories "github.com/hubkudev/sentinel/repos"
)

type Middleware interface {
	ProtectedRoute(c *fiber.Ctx) error
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
	exist, err := m.UserRepo.CheckIDExist(context.Background(), userID.(string))
	if !exist {
		return c.Status(fiber.StatusTemporaryRedirect).Redirect("/login")
	}

	c.Locals("userID", userID)

	return c.Next()
}
