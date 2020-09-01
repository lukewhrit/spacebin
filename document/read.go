package document

import (
	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/structs"
)

func registerRead(api fiber.Router) {
	api.Get("/:id", func(c *fiber.Ctx) {
		if c.Params("id") != "" && len(c.Params("id")) == config.Config.Documents.IDLength {
			document, err := GetDocument(c.Params("id"))

			if err != nil {
				c.Status(500).JSON(&structs.Response{
					Status:  c.Fasthttp.Response.StatusCode(),
					Payload: structs.Payload{},
					Error:   err.Error(),
				})

				return
			}

			c.Status(201).JSON(&structs.Response{
				Status: c.Fasthttp.Response.StatusCode(),
				Payload: structs.Payload{
					ID:        &document.ID,
					Content:   &document.Content,
					Extension: &document.Extension,
					CreatedAt: &document.CreatedAt,
					UpdatedAt: &document.UpdatedAt,
				},
				Error: "",
			})
		}
	})
}
