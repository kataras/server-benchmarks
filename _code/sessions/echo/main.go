package main

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

	uuid "github.com/iris-contrib/go.uuid"
)

// Session middleware facilitates HTTP session management backed by gorilla/sessions.
// https://echo.labstack.com/middleware/session
func main() {
	app := echo.New()
	storeKey := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(storeKey)

	app.GET("/sessions", func(ctx echo.Context) error {
		session, _ := store.Get(ctx.Request(), "session")
		id, _ := uuid.NewV4()
		session.Values["ID"] = id.String()

		session.Values["name"] = "John Doe"
		if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
			return ctx.String(500, err.Error())
		}

		v := session.Values["name"]
		name, ok := v.(string)
		if !ok {
			return ctx.String(500, "server error")
		}

		return ctx.String(200, name)
	})

	app.Start(":5000")
}
