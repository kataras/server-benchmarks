package main

import (
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/x/sessions"
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
		ID    int `json:"id"`
		Count int `json:"count"`
	}
)

func handler(ctx buffalo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Buffalo does not support parameter type-based routing.
		return err
	}

	var in []testInput
	if err := ctx.Bind(&in); err != nil {
		return err
	}

	return ctx.Render(200, render.JSON(testOutput{
		ID:    id,
		Count: len(in),
	}))
}

func main() {
	app := buffalo.New(buffalo.Options{
		Env:          "production", // default was "development".
		Addr:         ":5000",
		LogLvl:       logger.ErrorLevel, // default was "DebugLevel".
		WorkerOff:    true,
		SessionStore: sessions.Null{}, // disable sessions.
	})

	app.POST("/{id}", handler)

	app.Serve()
}
