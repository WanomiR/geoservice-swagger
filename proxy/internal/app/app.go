package app

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy/internal/modules"
	"proxy/internal/utils/readresponder"
	"syscall"
	"time"
)

type App struct {
	server      *http.Server
	signalChan  chan os.Signal
	config      modules.ServicesConfig
	services    *modules.Services
	controllers *modules.Controllers
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Serve() {
	fmt.Println("Started server on port", a.config.Port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Signal() <-chan os.Signal {
	return a.signalChan
}

func (a *App) readConfig(configPath string) error {
	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}

	a.config.Port = os.Getenv("PORT")
	a.config.JwtSecret = os.Getenv("JWT_SECRET")
	a.config.JwtAlg = os.Getenv("JWT_ALG")
	a.config.ApiKey = os.Getenv("DADATA_API_KEY")
	a.config.SecretKey = os.Getenv("DADATA_SECRET_KEY")
	a.config.ProxyHost = os.Getenv("PROXY_HOST")
	a.config.ProxyPort = os.Getenv("PROXY_PORT")

	return nil
}

func (a *App) init() error {

	if err := a.readConfig(".env"); err != nil {
		return err
	}

	a.services = modules.NewServices(a.config)
	rr := readresponder.NewReadRespond(readresponder.WithMaxBytes(1 << 20))
	a.controllers = modules.NewControllers(a.services, rr)

	a.server = &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	return nil
}

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(a.services.Proxy.ProxyReverse)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", a.controllers.Auth.Register)
		r.Post("/login", a.controllers.Auth.Authenticate)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello from API"))
		})

		r.Route("/address", func(r chi.Router) {
			r.Use(a.services.Auth.RequireAuthentication)
			r.Post("/search", a.controllers.Geo.AddressSearch)
			r.Post("/geocode", a.controllers.Geo.AddressGeocode)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", a.config.Port)),
	))

	return r
}
