package controller

import "net/http"

type Authenticator interface {
	Register(w http.ResponseWriter, r *http.Request)
	Authenticate(w http.ResponseWriter, r *http.Request)
}
