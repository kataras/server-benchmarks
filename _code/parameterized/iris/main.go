package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()
	app.Get("/hello/{name}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("name"))
	})

	app.Listen(":5000")
}
