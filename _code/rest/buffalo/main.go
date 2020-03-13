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
		Email string `json:"email"`
	}

	testOutput struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func handler(ctx buffalo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Buffalo does not support parameter type-based routing.
		return err
	}

	var in testInput
	if err := ctx.Bind(&in); err != nil {
		return err
	}

	return ctx.Render(200, render.JSON(testOutput{
		ID:   id,
		Name: in.Email,
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
