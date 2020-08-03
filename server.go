package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"

	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/database"
	"github.com/spacebin-org/curiosity/document"
	"github.com/spacebin-org/curiosity/middlewares"
)

func initDatabase() {
	var err error

	database.DBConn, err = pop.Connect("main")

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}
}

func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Couldn't load configuration file: %v", err)
	}

	// Initialize application
	app := fiber.New(&fiber.Settings{
		Prefork: true,
	})

	// Register middleware and endpoints
	registerMiddlewares(app)
	document.Register(app)

	// Initialize Database
	initDatabase()

	address := fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)

	log.Fatal(app.Listen(address))
}

func registerMiddlewares(app *fiber.App) {
	// Setup middlewares
	app.Use(middleware.Compress(middleware.CompressConfig{
		Level: config.GetConfig().Server.CompresssLevel,
	}))

	app.Use(limiter.New(limiter.Config{
		Timeout: config.GetConfig().Server.Ratelimits.Duration,
		Max:     config.GetConfig().Server.Ratelimits.Requests,
	}))

	app.Use(cors.New())
	app.Use(middlewares.SecurityHeaders())
	app.Use(middleware.Logger())
}
