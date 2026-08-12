package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// "GET /{$}" matches exactly "/", like the routers of the other apps.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Index")
	})

	http.ListenAndServe(":5000", mux)
}
