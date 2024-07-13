package service

import (
	"net/http"
	"proxy/modules/model"
	"proxy/modules/service/geoservice"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*geoservice.Address, error)
	GeoCode(lat, lng string) ([]*geoservice.Address, error)
}

type UserAuthenticator interface {
	Register(user model.User) error
	Authenticate(user model.User) (string, error)
	RequireAuthentication(http.Handler) http.Handler
}

type ProxyReverser interface {
	ProxyReverse(http.Handler) http.Handler
}
