package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
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

func handler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
	r := chi.NewRouter()
	r.Post("/{id}", handler)
	http.ListenAndServe(":5000", r)
}
