package user

import (
	"golang.org/x/crypto/bcrypt"
)

// User is the representation of supergiant user
type User struct {
	APIToken          string
	Login             string
	EncryptedPassword []byte
	Password          string
}

func (m *User) encryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = ""
	m.EncryptedPassword = hashedPassword
	return nil
}
