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
		Prefork: config.Config.Server.Prefork,
	})

	registerMiddlewares(app)
	registerRoutes(app)

	address := fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port)

	log.Fatal(app.Listen(address))
}
