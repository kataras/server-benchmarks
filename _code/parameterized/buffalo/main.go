package main

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/x/sessions"
)

func main() {
	app := buffalo.New(buffalo.Options{
		Env:          "production", // default was "development".
		Addr:         ":5000",
		LogLvl:       logger.ErrorLevel, // default was "DebugLevel".
		WorkerOff:    true,
		SessionStore: sessions.Null{}, // disable sessions.
	})

	app.GET("/hello/{name}", func(ctx buffalo.Context) error {
		return ctx.Render(200, render.String("Hello %s", ctx.Param("name")))
	})

	app.Serve()
}
