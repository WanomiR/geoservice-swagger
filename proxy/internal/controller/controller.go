package controller

import "net/http"

type UserController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Authenticate(w http.ResponseWriter, r *http.Request)
	RequireAuthentication(next http.Handler) http.Handler
}

type AddressController interface{}
