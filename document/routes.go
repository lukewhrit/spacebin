package document

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/structs"
)

// Register contains all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/api/v1/document")

	api.Post("/", func(c *fiber.Ctx) {
		id, err := NewDocument("this is a test", "txt")

		if err != nil {
			c.JSON(&structs.Response{
				Status:  c.Fasthttp.Response.StatusCode(),
				Payload: structs.Payload{},
				Error:   err.Error(),
			})

			return
		}

		document, err := GetDocument(id)

		if err != nil {
			fmt.Println(err)

			c.JSON(&structs.Response{
				Status:  c.Fasthttp.Response.StatusCode(),
				Payload: structs.Payload{},
				Error:   err.Error(),
			})

			return
		}

		c.Status(201).JSON(&structs.Response{
			Status: c.Fasthttp.Response.StatusCode(),
			Payload: structs.Payload{
				ID:          &document.ID,
				ContentHash: b64.StdEncoding.EncodeToString([]byte(document.Content)),
			},
			Error: "",
		})
	})
}
