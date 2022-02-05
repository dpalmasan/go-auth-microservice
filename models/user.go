package models

import "github.com/go-auth-microservice/types"

type UserModel interface {
	GetUserByEmail(email string) types.User
	AddUser(user types.User)
}
