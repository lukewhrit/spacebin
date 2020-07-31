package document

import "github.com/gofiber/fiber"

// Register contains all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/api/v1/documents")

	api.Get("/*", func(c *fiber.Ctx) {
		c.SendStatus(501)
	})
}
