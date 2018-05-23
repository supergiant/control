package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	APIToken          string
	Login             string
	EncryptedPassword []byte
	Password          string
}

type Repository interface {
	GetAll(context.Context) ([]User, error)
	Get(ctx context.Context, login string) (*User, error)
	Create(context.Context, *User) (error)
}

func (m *User) encryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.EncryptedPassword = hashedPassword
	return nil
}
