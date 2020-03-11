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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index")
	})

	app.MapTo(r, (*martini.Routes)(nil))
	app.Action(r.Handle)
	app.RunOnAddr(":5000")
}
