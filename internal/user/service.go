package user

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

const prefix = "/supergiant/user/"

// Service contains business logic related to users
type Service struct {
	Repository storage.Interface
}

func (s *Service) Create(ctx context.Context, user *User) error {
	logrus.Debug("user.Service.Create start")
	logrus.Debug("user.Service.Create end")
	return nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	logrus.Debug("user.Service.Authenticate start")
	rawJSON, err := s.Repository.Get(ctx, prefix, username)
	if err != nil {
		return err
	}
	user := new(User)
	err = json.Unmarshal(rawJSON, user)
	if err != nil {
		return errors.WithStack(err)
	}
	logrus.Debug("user.Service.Authenticate end")
	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}
