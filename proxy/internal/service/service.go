package service

import (
	"net/http"
	"proxy/internal/entities"
	"proxy/internal/service/geoservice"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*geoservice.Address, error)
	GeoCode(lat, lng string) ([]*geoservice.Address, error)
}

type Authenticator interface {
	Register(user entities.User) error
	Authenticate(user entities.User) (string, error)
	RequireAuthentication(http.Handler) http.Handler
}

type ProxyReverser interface {
	ProxyReverse(http.Handler) http.Handler
}
