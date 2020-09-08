/*
 * Copyright 2020 Luke Whrit, Jack Dorland; The Spacebin Authors

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 *  Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	api.Get("/:id/raw", func(c *fiber.Ctx) {
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

			c.Status(201).Send(document.Content)
		}
	})
}
