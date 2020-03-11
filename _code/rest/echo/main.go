package main

import (
	"strconv"

	"github.com/labstack/echo/v4"
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

func handler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Echo does not support parameter type-based routing.
		return ctx.String(404, err.Error())
	}

	var in testInput
	if err := ctx.Bind(&in); err != nil {
		return ctx.String(400, err.Error())
	}

	return ctx.JSON(200, testOutput{
		ID:   id,
		Name: in.Email,
	})
}

func main() {
	app := echo.New()
	app.POST("/:id", handler)
	app.Start(":5000")
}
