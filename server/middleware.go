package server

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"
	"github.com/spacebin-org/curiosity/config"
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
