package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/gorilla/securecookie"
	"github.com/martini-contrib/sessions"

	uuid "github.com/iris-contrib/go.uuid"
)

func main() {
	martini.Env = martini.Prod

	app := martini.New()
	storeKey := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(storeKey)
	app.Use(sessions.Sessions("session", store))

	r := martini.NewRouter()
	r.Get("/sessions", func(w http.ResponseWriter, session sessions.Session) {
		id, _ := uuid.NewV4()
		session.Set("ID", id.String())

		session.Set("name", "John Doe")
		v := session.Get("name")
		name, ok := v.(string)
		if !ok {
			w.WriteHeader(500)
			fmt.Fprint(w, "server error")
			return
		}

		fmt.Fprint(w, name)
	})

	app.MapTo(r, (*martini.Routes)(nil))
	app.Action(r.Handle)
	app.RunOnAddr(":5000")
}
