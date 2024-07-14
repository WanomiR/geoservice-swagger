package controller

import (
	"net/http"
	"proxy/internal/service"
	"proxy/utils/readresponder"
)

type Controller interface {
	Register(w http.ResponseWriter, r *http.Request)
	Authenticate(w http.ResponseWriter, r *http.Request)
	AddressSearch(w http.ResponseWriter, r *http.Request)
	AddressGeocode(w http.ResponseWriter, r *http.Request)
}

type AppController struct {
	readResponder readresponder.ReadResponder
	authService   service.Authenticator
	geoService    service.GeoProvider
}

type AppControllerOption func(*AppController)

func WithResponder(readResponder readresponder.ReadResponder) AppControllerOption {
	return func(c *AppController) {
		c.readResponder = readResponder
	}
}

func WithAuthenticator(authService service.Authenticator) AppControllerOption {
	return func(c *AppController) {
		c.authService = authService
	}
}

func WithGeoService(geoService service.GeoProvider) AppControllerOption {
	return func(c *AppController) {
		c.geoService = geoService
	}
}

func NewAppController(options ...AppControllerOption) *AppController {
	controller := &AppController{}

	for _, option := range options {
		option(controller)
	}

	return controller
}
