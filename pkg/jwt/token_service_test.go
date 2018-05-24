package jwt

import (
	"testing"
)

func TestTokenService(t *testing.T) {
	ts := &TokenService{
		60,
		[]byte("secret key"),
	}

	userId := "user_id"

	tokenString, err := ts.Issue(userId)

	if err != nil {
		t.Error(err)
		return
	}

	claims, err := ts.Validate(tokenString)

	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := claims["user_id"]; !ok {
		t.Errorf("user_id not found in token claims")
		return
	}
}
