package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
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

func handler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		// * Chi does not support parameter type-based routing.
		w.WriteHeader(404)
		return
	}

	var in testInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(testOutput{
		ID:   id,
		Name: in.Email,
	})
}

func main() {
	martini.Env = martini.Prod

	app := martini.New()
	r := martini.NewRouter()
	r.Post("/:id", handler)
	app.MapTo(r, (*martini.Routes)(nil))
	app.Action(r.Handle)
	app.RunOnAddr(":5000")
}
