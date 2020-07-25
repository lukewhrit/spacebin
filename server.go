package main

import "github.com/gofiber/fiber"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(fiber.Map{
			"message": "Hello, World 👋!",
		})
	})

	app.Listen(3000)
}
