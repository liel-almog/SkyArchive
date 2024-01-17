package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/be-file-streaming/routes"
)

func setupRouter(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello World!",
		})
	})

	api := app.Group("/api")
	routes.NewUploadRouter(api)
}
