package service

import (
	"net/http"
	"proxy/internal/entities"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*entities.Address, error)
	GeoCode(lat, lng string) ([]*entities.Address, error)
}

type Authenticator interface {
	Register(user entities.User) error
	Authenticate(user entities.User) (string, error)
	RequireAuthentication(http.Handler) http.Handler
}

type ProxyReverser interface {
	ProxyReverse(next http.Handler) http.Handler
}
