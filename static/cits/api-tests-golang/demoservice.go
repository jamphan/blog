package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Request-Id", "Some-Request-ID")
		w.Header().Set("Service-Name", "Demo")

		accept := r.Header.Get("Accept")
		if accept != "text/plain" {
			w.WriteHeader(400)
			w.Write([]byte("Must set accept to text/plain"))
			return
		}
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":3000", r)
}