package main

import "github.com/gofiber/fiber/v3"

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Index")
	})

	app.Listen(":5000", fiber.ListenConfig{DisableStartupMessage: true})
}
