package repository

import (
	"proxy/internal/entities"
)

type DatabaseRepo interface {
	GetUserByEmail(string) (entities.User, error)
	InsertUser(entities.User) error
}
