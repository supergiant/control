package user

import "context"

type User struct {
	APIToken string
}

type Repository interface {
	GetAll(context context.Context) ([]User, error)
}
