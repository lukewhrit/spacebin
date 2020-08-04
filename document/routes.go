package document

import "github.com/gofiber/fiber"

// Register loads all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/api/v1/documents")

	registerCreate(api)
}
