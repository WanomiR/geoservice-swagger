package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

type Config struct {
	jwtSecret string
	apiKey    string
	secretKey string
	host      string
	port      string
	proxyHost string
	proxyPort string
}

type App struct {
	serviceProvider *serviceProvider
	config          Config
	router          *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("listening on port " + a.config.port)
	return http.ListenAndServe(":"+a.config.port, a.router)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initRoutes,
	}

	for _, init := range inits {
		err := init(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	a.config.jwtSecret = os.Getenv("JWT_SECRET")
	a.config.apiKey = os.Getenv("DADATA_API_KEY")
	a.config.secretKey = os.Getenv("DADATA_SECRET_KEY")
	a.config.host = os.Getenv("HOST")
	a.config.port = os.Getenv("PORT")
	a.config.proxyHost = os.Getenv("PROXY_HOST")
	a.config.proxyPort = os.Getenv("PROXY_PORT")

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = &serviceProvider{}
	return nil
}

func (a *App) initRoutes(_ context.Context) error {
	r := chi.NewRouter()
	proxy := a.serviceProvider.ProxyService(a.config.proxyHost, a.config.proxyPort)

	r.Use(proxy.ProxyReverse)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", a.serviceProvider.UserController(a.config.jwtSecret).Register)
		r.Post("/login", a.serviceProvider.UserController(a.config.jwtSecret).Authenticate)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello from API"))
		})

		//r.Route("/address", func(r chi.Router) {
		//	r.Use(a.RequireAuthentication)
		//	r.Post("/search", g.HandleAddressSearch)
		//	r.Post("/geocode", g.HandleAddressGeocode)
		//})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	a.router = r

	return nil
}
