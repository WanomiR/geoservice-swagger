package repository

import (
	"proxy/internal/modules/auth/entities"
)

//go:generate mockgen -source=./repository.go -destination=./mock/mock_repo.go
type DatabaseRepo interface {
	GetUserByEmail(string) (entities.User, error)
	InsertUser(entities.User) error
}
