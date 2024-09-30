package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func handler(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2<<20) // 2MB.

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// * Gin does not support parameter type-based routing.
		ctx.Status(404)
		return
	}

	var in []testInput
	if err := ctx.BindJSON(&in); err != nil {
		ctx.Status(400)
		return
	}

	ctx.JSON(200, testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
	})
}

func main() {
	gin.SetMode("release")

	app := gin.New()
	app.POST("/:id", handler)
	app.Run(":5000")
}
