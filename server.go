package main

import (
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
	"github.com/spacebin-org/curiosity/document"
	"github.com/spacebin-org/curiosity/middlewares"
)

func initDatabase() {
	var err error

	// Connect to database
	database.DBConn, err = gorm.Open(
		config.GetConfig().Database.Dialect,
		config.GetConfig().Database.ConnectionURI,
	)

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}

	// Setup database
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
