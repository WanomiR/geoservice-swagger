package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"proxy/internal/controller"
	"proxy/internal/entities"
	"proxy/internal/repository/dbrepo"
	"proxy/internal/service"
	"proxy/internal/service/auth"
	"proxy/internal/service/geoservice"
	"proxy/internal/service/reverse"
	"proxy/utils/readresponder"
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
	config       Config
	proxyService service.ProxyReverser
	authService  service.Authenticator
	controller   controller.Controller
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("listening on port " + a.config.port)
	return http.ListenAndServe(":"+a.config.port, a.routes())
}

func (a *App) readConfig() error {
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

func (a *App) init() error {

	if err := a.readConfig(); err != nil {
		return err
	}

	a.proxyService = reverse.NewProxyReverse(a.config.proxyHost, a.config.proxyPort)

	db := dbrepo.NewMapDBRepo(entities.User{Email: "admin@example.com", Password: "password"})
	a.authService = auth.NewUserAuth("HS256", a.config.jwtSecret, auth.WithDatabase(db))

	rr := readresponder.NewReadRespond(readresponder.WithMaxBytes(1 << 20))
	geo := geoservice.NewGeoService(a.config.apiKey, a.config.secretKey)

	a.controller = controller.NewAppController(
		controller.WithResponder(rr),
		controller.WithAuthenticator(a.authService),
		controller.WithGeoService(geo),
	)

	return nil
}

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(a.proxyService.ProxyReverse)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", a.controller.Register)
		r.Post("/login", a.controller.Authenticate)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello from API"))
		})

		r.Route("/address", func(r chi.Router) {
			r.Use(a.authService.RequireAuthentication)
			r.Post("/search", a.controller.AddressSearch)
			r.Post("/geocode", a.controller.AddressGeocode)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	return r
}
