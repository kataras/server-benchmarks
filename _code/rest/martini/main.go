package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
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

const contentTypeKey = "Content-Type"
const contentTypeValue = "application/json; charset=utf-8"

func handler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB.

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		// * Chi does not support parameter type-based routing.
		w.WriteHeader(404)
		return
	}

	var in []testInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(400)
		return
	}

	w.Header().Add(contentTypeKey, contentTypeValue)

	json.NewEncoder(w).Encode(testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
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
