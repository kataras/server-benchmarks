package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/hello/:name", func(ctx echo.Context) error {
		return ctx.String(200, fmt.Sprintf("Hello %s", ctx.Param("name")))
	})

	app.Start(":5000")
}
