package server

import (
	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/document"
)

func registerRoutes(app *fiber.App) {
	document.Register(app)
}
