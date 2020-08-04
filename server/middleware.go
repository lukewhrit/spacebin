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
		Level: config.GetConfig().Server.CompresssionLevel,
	}))

	app.Use(limiter.New(limiter.Config{
		Timeout: config.GetConfig().Server.Ratelimits.Duration,
		Max:     config.GetConfig().Server.Ratelimits.Requests,
	}))

	app.Use(cors.New())
	app.Use(SecurityHeaders())
	app.Use(middleware.Logger())
}
