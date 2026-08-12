package main

import "github.com/gofiber/fiber/v3"

func main() {
	app := fiber.New()

	app.Get("/hello/:name", func(c fiber.Ctx) error {
		return c.SendString("Hello " + c.Params("name"))
	})

	app.Listen(":5000", fiber.ListenConfig{DisableStartupMessage: true})
}
