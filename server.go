package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/database"
	"github.com/spacebin-org/curiosity/database/models"
	"github.com/spacebin-org/curiosity/middlewares"
)

func initDatabase() {
	var err error

	// Connect to database
	database.DBConn, err = gorm.Open(
		config.GetDatabase().Dialect,
		config.GetDatabase().ConnectionURI,
	)

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}

	// Setup database
	database.DBConn.CreateTable(&models.Document{})
	database.DBConn.AutoMigrate(&models.Document{})

	database.DBConn.Create(&models.Document{
		Key:     "abcdef",
		Content: "this is a test",
	})
}

func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Couldn't load configuration file: %v", err)
	}

	// Initialize application
	app := fiber.New()

	// Register middleware and endpoints
	registerMiddlewares(app)
	registerEndpoints(app)

	// Initialize Database
	initDatabase()

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

func registerMiddlewares(app *fiber.App) {
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
}

func registerEndpoints(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"message": "Hello, World 👋!",
		})
	})
}
