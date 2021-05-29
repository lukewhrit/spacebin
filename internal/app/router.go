/*
 * Copyright 2020-2021 Luke Whrit, Jack Dorland

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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spacebin-org/spirit/internal/pkg/config"
	"github.com/spacebin-org/spirit/internal/pkg/document"
)

func registerRouter(app *fiber.App) {
	// Setup middlewares
	app.Use(compress.New(compress.Config{
		Level: config.Config.Server.CompresssionLevel,
	}))

	app.Use(limiter.New(limiter.Config{
		Duration: config.Config.Server.Ratelimits.Duration,
		Max:      config.Config.Server.Ratelimits.Requests,
	}))

	app.Use(cors.New())
	app.Use(logger.New())

	// Custom middleware to set security-related headers
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers
		c.Set("X-Download-Options", "noopen")
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer-when-downgrade")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';")

		// Go to next middleware
		return c.Next()
	})

	document.Register(app)
}
