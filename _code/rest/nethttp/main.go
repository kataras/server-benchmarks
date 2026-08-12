package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func handler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB.

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// * net/http does not support parameter type-based routing.
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var in []testInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(testOutput{
		ID:      id,
		Count:   len(in),
		FirstID: in[0].ID,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{id}", handler)

	http.ListenAndServe(":5000", mux)
}
