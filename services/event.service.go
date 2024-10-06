package services

import "github.com/gofiber/fiber/v2"

type EventService interface {
	HelloWorld(c *fiber.Ctx) error
}

type EventServiceImpl struct{}

func (s *EventServiceImpl) HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
