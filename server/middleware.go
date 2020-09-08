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

package server

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"
	"github.com/spacebin-org/spirit/config"
)

func registerMiddlewares(app *fiber.App) {
	// Setup middlewares
	app.Use(middleware.Compress(middleware.CompressConfig{
		Level: config.Config.Server.CompresssionLevel,
	}))

	app.Use(limiter.New(limiter.Config{
		Timeout: config.Config.Server.Ratelimits.Duration,
		Max:     config.Config.Server.Ratelimits.Requests,
	}))

	app.Use(cors.New())
	app.Use(middleware.Logger())

	// Custom middleware to set security-related headers
	app.Use(func(c *fiber.Ctx) {
		// Set some security headers:
		c.Set("X-Download-Options", "noopen")
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer-when-downgrade")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("Cache-Control", "max-age=31536000")

		if config.Config.Server.UseCSP == true {
			c.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';")
		}

		// Go to next middleware:
		c.Next()
	})
}
