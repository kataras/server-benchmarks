package main

import (
	"fmt"
	"net/http"

	"github.com/pressly/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index")
	})

	http.ListenAndServe(":5000", r)
}
