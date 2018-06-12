package user

import (
	"golang.org/x/crypto/bcrypt"
)

// User is the representation of supergiant user
type User struct {
	Login             string `json:"login" valid:"required, length(1|32)"`
	EncryptedPassword []byte `json:"encrypted_password" valid:"-"`
	Password          string `json:"password" valid:"required, length(8|24), printableascii"`
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
