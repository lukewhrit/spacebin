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

package app

import (
	"github.com/coral-dev/spirit/internal/pkg/config"
	"github.com/coral-dev/spirit/internal/pkg/domain"
	"github.com/gofiber/fiber/v2"
)

// Start initializes the server
func Start() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: config.Config.Server.Prefork,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 status code
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			// Return statuscode with error message
			return c.Status(code).JSON(&domain.Response{
				Error:   err.Error(),
				Payload: domain.Payload{},
				Status:  code,
			})
		},
	})

	registerRouter(app)

	return app
}
