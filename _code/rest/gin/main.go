package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
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

func handler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Gin does not support parameter type-based routing.
		ctx.Status(404)
		return
	}

	var in testInput
	if err := ctx.BindJSON(&in); err != nil {
		ctx.Status(400)
		return
	}

	ctx.JSON(200, testOutput{
		ID:   id,
		Name: in.Email,
	})
}

func main() {
	gin.SetMode("release")

	app := gin.New()
	app.POST("/:id", handler)
	app.Run(":5000")
}
