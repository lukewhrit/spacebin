package routes

import "github.com/gofiber/fiber"

// Document related routes
func Document(app *fiber.App) {
	api := app.Group("/api/v1/documents")

	api.Get("/*", func(c *fiber.Ctx) {
		c.SendStatus(501)
	})
}
