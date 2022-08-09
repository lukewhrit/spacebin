/*
 * Copyright 2020-2022 Luke Whrit, Jack Dorland

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package document

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/coral-dev/spirit/internal/pkg/config"
	"github.com/coral-dev/spirit/internal/pkg/domain"
	"github.com/gofiber/fiber/v2"
)

// Register loads all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/v1/documents")

	api.Post("/", func(c *fiber.Ctx) error {
		b := new(CreateRequest)

		// Validate and parse body
		if err := c.BodyParser(b); err != nil {
			return fiber.NewError(400, err.Error())
		}

		if err := b.Validate(); err != nil {
			return fiber.NewError(400, err.Error())
		}

		// Create and retrieve document
		id, err := NewDocument(b.Content, b.Extension)

		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		document, err := GetDocument(id)

		if err != nil {
			return fiber.NewError(500, err.Error())
		}

		hash := md5.Sum([]byte(document.Content))

		c.Status(201).JSON(&domain.Response{
			Status: c.Response().StatusCode(),
			Payload: domain.Payload{
				ID:          &document.ID,
				ContentHash: hex.EncodeToString(hash[:]),
			},
			Error: "",
		})

		return nil
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		if c.Params("id") != "" && len(c.Params("id")) == config.Config.Documents.IDLength {
			document, err := GetDocument(c.Params("id"))

			if err != nil {
				return fiber.NewError(404, err.Error())
			}

			c.Status(200).JSON(&domain.Response{
				Status: c.Response().StatusCode(),
				Payload: domain.Payload{
					ID:        &document.ID,
					Content:   &document.Content,
					Extension: &document.Extension,
					CreatedAt: &document.CreatedAt,
					UpdatedAt: &document.UpdatedAt,
				},
				Error: "",
			})
		} else {
			return fiber.NewError(400)
		}

		return nil
	})

	api.Get("/:id/raw", func(c *fiber.Ctx) (err error) {
		if c.Params("id") != "" && len(c.Params("id")) == config.Config.Documents.IDLength {
			document, err := GetDocument(c.Params("id"))

			if err != nil {
				return fiber.NewError(404, err.Error())
			}

			c.Status(200).SendString(document.Content)
		} else {
			return fiber.NewError(400)
		}

		return nil
	})

}
