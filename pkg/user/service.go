package user

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

const prefix = "/user/"

// Service contains business logic related to users
type Service struct {
	Repository storage.Interface
}

func (s *Service) GetByToken(ctx context.Context, apiToken string) (*User, error) {
	logrus.Debug("user.Service.GetByToken start")
	result, err := s.Repository.GetAll(ctx, prefix)
	if err != nil {
		return nil, err
	}
	for _, rawUser := range result {
		u := new(User)
		err = json.NewDecoder(bytes.NewReader(rawUser)).Decode(u)
		if err != nil {
			logrus.Warningf("failed to convert stored data to cloud account struct")
			logrus.Debugf("corrupted data: %s", rawUser)
			continue
		}
		if u.APIToken == apiToken {
			return u, nil
		}
	}
	logrus.Debug("user.Service.GetByToken end")
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
