package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	_ "proxy/docs"
	"proxy/geoservice"
	"proxy/reverse"
)

// @title Geoservice API
// @version 1.0
// @description Find matching addresses by street name or coordinates

// @host localhost:8080
// @BasePath /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")
	port := os.Getenv("PORT")

	g := geoservice.NewGeoService(apiKey, secretKey)

	r := chi.NewRouter()

	proxy := reverse.NewReverseProxy("hugo", "1313")
	r.Use(proxy.ReverseProxy)

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from API"))
	})

	r.Route("/api/address", func(r chi.Router) {
		r.Post("/search", g.HandleAddressSearch)
		r.Post("/geocode", g.HandleAddressGeocode)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost"+port+"/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(port, r))
}
