package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type TokenService struct {
	tokenTTL  int64
	secretKey []byte
}

func NewTokenService(tokenTTL int64, secret []byte) *TokenService {
	return &TokenService{
		tokenTTL:  tokenTTL,
		secretKey: secret,
	}
}

func (ts TokenService) Issue(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		// TODO(stgleb): Pass list of access here
		"accesses":   []string{"edit", "view"},
		"user_id":    userId,
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().Unix() + ts.tokenTTL,
	})

	tokenString, err := token.SignedString(ts.secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ts TokenService) Validate(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		e, ok := claims["expires_at"].(float64)

		if !ok {
			return nil, sgerrors.ErrInvalidCredentials
		}

		expiresAt := int64(e)

		if int64(expiresAt) < time.Now().Unix() {
			return nil, sgerrors.ErrTokenExpired
		}

		return claims, nil
	} else {
		return nil, errors.New("Error while converting to jwt claims map")
	}

	return nil, nil
}
