package main

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/logger"

	uuid "github.com/iris-contrib/go.uuid"
)

func main() {
	app := buffalo.New(buffalo.Options{
		Env:       "production", // default was "development".
		Addr:      ":5000",
		LogLvl:    logger.ErrorLevel, // default was "DebugLevel".
		WorkerOff: true,
	})

	app.GET("/sessions", func(ctx buffalo.Context) error {
		session := ctx.Session()
		id, _ := uuid.NewV4()
		session.Set("ID", id.String())

		session.Set("name", "John Doe")
		if err := session.Save(); err != nil {
			ctx.Render(500, render.String(err.Error()))
			return nil
		}

		v := session.Get("name")
		name, ok := v.(string)
		if !ok {
			ctx.Render(500, render.String("server error"))
			return nil
		}

		return ctx.Render(200, render.String(name))
	})

	app.Serve()
}
