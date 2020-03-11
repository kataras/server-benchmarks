package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode("release")

	app := gin.New()
	app.GET("/hello/:name", func(ctx *gin.Context) {
		ctx.String(200, "Hello %s", ctx.Param("name"))
	})

	app.Run(":5000")
}
