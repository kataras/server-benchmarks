package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
)

func main() {
	martini.Env = martini.Prod

	app := martini.New()
	r := martini.NewRouter()
	r.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		fmt.Fprintf(w, "Hello %s", params["name"])
	})
	app.MapTo(r, (*martini.Routes)(nil))
	app.Action(r.Handle)
	app.RunOnAddr(":5000")
}
