package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/config"
)

// Start initializes the server
func Start() {
	app := fiber.New(&fiber.Settings{
		Prefork: config.GetConfig().Server.Prefork,
	})

	registerMiddlewares(app)
	registerRoutes(app)

	address := fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)

	log.Fatal(app.Listen(address))
}
