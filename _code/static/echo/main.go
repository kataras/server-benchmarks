package main

import "github.com/labstack/echo/v4"

func main() {
	app := echo.New()
	app.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "Index")
	})

	app.Start(":5000")
}
