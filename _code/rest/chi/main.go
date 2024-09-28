package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		ID    int `json:"id"`
		Count int `json:"count"`
	}
)

const contentTypeKey = "Content-Type"
const contentTypeValue = "application/json; charset=utf-8"

func handler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB.

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
		ID:    id,
		Count: len(in),
	})
}

func main() {
	r := chi.NewRouter()
	r.Post("/{id}", handler)
	http.ListenAndServe(":5000", r)
}
