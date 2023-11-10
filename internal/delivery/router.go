package delivery

import (
	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(handler *fiber.App, services *service.Services) {
	handler.Use(setRequestID)
	handler.Use(logRequest)

	handler.Get("/health", func(c *fiber.Ctx) error { return c.SendString("200") })
	v1 := handler.Group("/api/v1")
	{
		newLinkRoutes(v1.Group("/shortener"), services.Link)
	}
}
