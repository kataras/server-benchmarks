package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello "+r.PathValue("name"))
	})

	http.ListenAndServe(":5000", mux)
}
