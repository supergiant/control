package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/storage"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

const prefix = "/user/"

// Service contains business logic related to users
type Service struct {
	Repository storage.Interface
}

func (s *Service) RegisterUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	logrus.Debug("user.Service.Authenticate start")
	obj, err := s.Repository.Get(ctx, prefix, username, reflect.TypeOf((*User)(nil)))
	if err != nil {
		return err
	}
	user, ok := obj.(User)
	if !ok {
		return errors.New("corrupted data")
	}
	logrus.Debug("user.Service.Authenticate end")
	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}
