package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	proxy := NewReverseProxy("hugo", "1313")
	r.Use(proxy.ReverseProxy)

	r.Route("/api/address", func(r chi.Router) {
		r.Post("/search", handleAddressSearch)
		r.Post("/geocode", handleAddressGeocode)
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from API"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
