package document

import (
	"fmt"

	"github.com/gofiber/fiber"
)

// Register contains all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/api/v1/document")

	api.Post("/", func(c *fiber.Ctx) {
		id, err := NewDocument("this is a test", "txt")

		if err != nil {
			c.JSON(&fiber.Map{
				"status":  c.Fasthttp.Response.StatusCode,
				"payload": fiber.Map{},
				"error":   err.Error(),
			})

			return
		}

		document, err := GetDocument(id)

		if err != nil {
			fmt.Println(err)

			c.JSON(&fiber.Map{
				"status":  c.Fasthttp.Response.StatusCode,
				"payload": fiber.Map{},
				"error":   err.Error(),
			})

			return
		}

		c.Status(201).JSON(&fiber.Map{
			"status":  c.Fasthttp.Response.StatusCode(),
			"payload": document,
			"error":   fiber.Map{},
		})
	})
}
