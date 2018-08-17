package user

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

// User is the representation of supergiant user
type User struct {
	Login             string `json:"login" valid:"required, length(1|32)"`
	EncryptedPassword []byte `json:"encrypted_password" valid:"-"`
	Password          string `json:"password" valid:"required, length(8|24), printableascii"`
}

func (u *User) encryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = ""
	u.EncryptedPassword = hashedPassword
	return nil
}

func (u *User) ToJSON() []byte {
	js, _ := json.Marshal(u)
	return js
}

func FromJSON(raw []byte) (*User, error) {
	usr := new(User)
	err := json.Unmarshal(raw, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
