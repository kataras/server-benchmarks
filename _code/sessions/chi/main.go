package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pressly/chi"

	uuid "github.com/iris-contrib/go.uuid"
)

// HTTP session management backed by gorilla/sessions
// as chi does not contain a session manager, most of its users use gorilla.
func main() {
	r := chi.NewRouter()
	storeKey := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(storeKey)

	r.Get("/sessions", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		id, _ := uuid.NewV4()
		session.Values["ID"] = id.String()

		session.Values["name"] = "John Doe"
		if err := session.Save(r, w); err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			return
		}

		v := session.Values["name"]
		name, ok := v.(string)
		if !ok {
			w.WriteHeader(500)
			fmt.Fprint(w, "server error")
			return
		}

		fmt.Fprint(w, name)
	})

	http.ListenAndServe(":5000", r)
}
