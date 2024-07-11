package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"proxy/auth"
	_ "proxy/docs"
	"proxy/geoservice"
	"proxy/reverse"
)

type Claims map[string]any

var tokenAuth *jwtauth.JWTAuth

// @title Geoservice API
// @version 2.0.0
// @description Geoservice with swagger docs and authentication

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)

	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")
	port := os.Getenv("PORT")

	g := geoservice.NewGeoService(apiKey, secretKey)
	a := auth.NewUserAuth("HS256", jwtSecret)

	r := chi.NewRouter()

	proxy := reverse.NewReverseProxy("hugo", "1313")
	r.Use(proxy.ReverseProxy)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", a.RegisterUser)
		r.Post("/login", a.Authenticate)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello from API"))
		})

		r.Route("/address", func(r chi.Router) {
			r.Use(a.RequireAuthentication)
			r.Post("/search", g.HandleAddressSearch)
			r.Post("/geocode", g.HandleAddressGeocode)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost"+port+"/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(port, r))
}
