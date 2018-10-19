package jwt

import (
	"bytes"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

func TestNewTokenService(t *testing.T) {
	ttl := int64(600)
	secret := []byte(`secret`)
	ts := NewTokenService(ttl, secret)

	if ts.tokenTTL != ttl {
		t.Errorf("expected ttl %d actual %d", ttl, ts.tokenTTL)
	}

	if !bytes.EqualFold(ts.secretKey, secret) {
		t.Errorf("expected secret %s actual %s",
			string(secret), string(ts.secretKey))
	}
}

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

func TestTokenService_ValidateErrExpiredNotFound(t *testing.T) {
	secret := []byte(`secret`)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"issued_at": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(secret)

	ts := &TokenService{
		60,
		secret,
	}

	claims, err := ts.Validate(tokenString)

	if err != sgerrors.ErrInvalidCredentials {
		t.Errorf("Wrong error ")
	}

	if claims != nil {
		t.Error("Claims must be nil")
	}
}

func TestTokenService_ValidateExpired(t *testing.T) {
	secret := []byte(`secret`)
	userId := "root"

	ts := &TokenService{
		-1,
		secret,
	}

	token, _ := ts.Issue(userId)

	claims, err := ts.Validate(token)

	if err != sgerrors.ErrTokenExpired {
		t.Errorf("Wrong error expected %v actual %v",
			sgerrors.ErrTokenExpired, err)
	}

	if claims != nil {
		t.Error("Claims must be nil")
	}
}
