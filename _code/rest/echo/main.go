package main

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func handler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Echo does not support parameter type-based routing.
		return ctx.String(404, err.Error())
	}

	var in []testInput
	if err := ctx.Bind(&in); err != nil {
		return ctx.String(400, err.Error())
	}

	return ctx.JSON(200, testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
	})
}

func main() {
	app := echo.New()
	app.Server.ReadHeaderTimeout = 20 * time.Second
	app.Server.WriteTimeout = 20 * time.Second

	app.Use(middleware.BodyLimit("2M"))
	app.POST("/:id", handler)
	app.Start(":5000")
}
