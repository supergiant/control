package user

import "testing"

func TestEncryptPassword(t *testing.T) {
	u := &User{
		Login:    "root",
		Password: "1234",
	}

	err := u.encryptPassword()

	if err != nil {
		t.Errorf("Error must be nil actual %v", err)
	}
}
