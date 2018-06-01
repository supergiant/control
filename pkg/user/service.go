package user

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/supergiant/supergiant/pkg/storage"
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
	logrus.Debug("user.Service.Authenticate start")
	rawJSON, err := s.Repository.Get(ctx, prefix, username)
	if err != nil {
		return err
	}
	user := new(User)
	err = json.NewDecoder(bytes.NewReader(rawJSON)).Decode(user)
	if err != nil {
		return errors.WithStack(err)
	}
	logrus.Debug("user.Service.Authenticate end")
	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}
