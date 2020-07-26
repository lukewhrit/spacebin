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
		log.Fatalf("Couldn't load configuration file: %v", err)
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
	registerEndpoints(app)

	listenString := fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort())

	// Only listen with TLS if cert & key are provided
	if config.GetTLS().Cert != "" && config.GetTLS().Key != "" {
		cert, err := tls.LoadX509KeyPair(config.GetTLS().Cert, config.GetTLS().Key)

		if err != nil {
			log.Fatal(err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		log.Fatal(app.Listen(listenString, tlsConfig))
	} else {
		log.Fatal(app.Listen(listenString))
	}
}

func registerEndpoints(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"message": "Hello, World 👋!",
		})
	})
}
