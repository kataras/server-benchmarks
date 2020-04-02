package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"

	uuid "github.com/iris-contrib/go.uuid"
)

func main() {
	app := iris.New()
	sessionsManager := sessions.New(sessions.Config{
		Cookie: "session",
		SessionIDGenerator: func(ctx iris.Context) string {
			id, _ := uuid.NewV4()
			return id.String()
		}},
	)

	app.Get("/sessions", func(ctx iris.Context) {
		session := sessionsManager.Start(ctx)
		session.Set("name", "John Doe")
		// Notes:
		// Save is done automatically on Iris.
		// Error reporting is done automatically too.

		name := session.GetString("name")
		ctx.WriteString(name)
	})

	app.Listen(":5000")
}
