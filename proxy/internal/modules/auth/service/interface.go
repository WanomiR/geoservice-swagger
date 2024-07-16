package service

import (
	"net/http"
	"proxy/internal/modules/auth/entities"
)

//go:generate mockgen -source=./interface.go -destination=../controller/mock_service/mock_service.go
type Authenticator interface {
	Register(user entities.User) error
	Authenticate(user entities.User) (string, error)
	RequireAuthentication(http.Handler) http.Handler
}
