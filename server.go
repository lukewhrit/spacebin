package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"

	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/middlewares"
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
		Level: config.GetCompressLevel(),
	}))

	app.Use(limiter.New(limiter.Config{
		Timeout: config.GetRatelimits().Duration,
		Max:     config.GetRatelimits().Requests,
	}))

	app.Use(cors.New())
	app.Use(middlewares.SecurityHeaders())
	app.Use(middleware.Logger())

	// Register endpoints
	endpoints(app)

	cert, err := tls.LoadX509KeyPair(config.GetTLS().Cert, config.GetTLS().Key)

	if err != nil {
		log.Fatal(err)
	}

	app.Listen(fmt.Sprintf("%s%d", config.GetHost(), config.GetPort()), &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
}

func endpoints(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"message": "Hello, World 👋!",
		})
	})
}
