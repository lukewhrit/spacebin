package main

import (
	"log"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"

	"github.com/spacebin-org/curiosity/config"
)

func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Was Unable to load configuration: %v", err)
	}

	// Initialize application
	app := fiber.New()

	// Setup middlewares
	app.Use(middleware.Compress(middleware.CompressConfig{
		Level: config.CompressLevel(),
	}))

	app.Use(limiter.New(limiter.Config{
		Timeout: config.Ratelimits().Duration,
		Max:     config.Ratelimits().Requests,
	}))

	app.Use(cors.New())

	app.Use(middleware.Logger())

	// Register endpoints
	endpoints(app)

	app.Listen(config.Port())
}

func endpoints(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"message": "Hello, World 👋!",
		})
	})
}
