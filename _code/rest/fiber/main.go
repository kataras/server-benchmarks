package main

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type (
	testInput struct {
		Name     string  `json:"name"`
		Language string  `json:"language"`
		ID       string  `json:"id"`
		Bio      string  `json:"bio"`
		Version  float64 `json:"version"`
	}

	testOutput struct {
		ID      int    `json:"id"`
		Count   int    `json:"count"`
		FirstID string `json:"first_id"`
	}
)

func handler(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// * Fiber does not support parameter type-based routing.
		return c.SendStatus(fiber.StatusNotFound)
	}

	var in []testInput
	if err := c.Bind().JSON(&in); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
	})
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // 2MB.
	})

	app.Post("/:id", handler)

	app.Listen(":5000", fiber.ListenConfig{DisableStartupMessage: true})
}
