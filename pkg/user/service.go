package user

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Service contains business logic related to users
type Service struct {
	Repository Repository
}

func (s *Service) GetByToken(ctx context.Context, apiToken string) (*User, error) {
	users, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.APIToken == apiToken {
			return &u, nil
		}
	}
	return nil, errors.Errorf("user with api token %s not found", apiToken)
}

func (s *Service) RegisterUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	user, err := s.Repository.Get(ctx, username)

	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}
