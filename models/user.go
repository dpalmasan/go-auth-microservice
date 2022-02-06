package models

import "github.com/go-auth-microservice/types"

type UserModel interface {
	GetByEmail(email string) (types.User, error)
	Add(user types.User) (types.User, error)
}
