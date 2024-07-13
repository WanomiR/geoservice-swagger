package repository

import (
	"proxy/modules/model"
)

type DatabaseRepo interface {
	GetUserByEmail(string) (model.User, error)
	InsertUser(model.User) error
}
