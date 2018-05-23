package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type TokenService struct {
	tokenTTL  int64
	secretKey []byte
}

func NewTokenService(tokenTTL int64, secret []byte) TokenService {
	return TokenService{
		tokenTTL:  tokenTTL,
		secretKey: secret,
	}
}

func (ts TokenService) Issue(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
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

func (ts TokenService) Validate(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return ts.secretKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		e, ok := claims["expires_at"].(float64)
		expiresAt := int64(e)

		if !ok {
			return errors.New("Token malformed")
		}

		if int64(expiresAt) < time.Now().Unix() {
			return errors.New("Token has been expired")
		}
	} else {
		return errors.New("Error while converting to jwt claims map")
	}

	return nil
}
