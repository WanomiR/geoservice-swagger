package controller

import (
	"net/http"
)

//go:generate mockgen
type GeoServicer interface {
	AddressSearch(w http.ResponseWriter, r *http.Request)
	AddressGeocode(w http.ResponseWriter, r *http.Request)
}
