package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode("release")

	app := gin.New()
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Index")
	})

	app.Run(":5000")
}
