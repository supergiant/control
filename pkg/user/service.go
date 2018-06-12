package user

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

const prefix = "/supergiant/user/"

// Service contains business logic related to users
type Service struct {
	repository storage.Interface
}

func (s *Service) Create(ctx context.Context, user *User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "user create")
	}
	if _, err := s.repository.Get(ctx, prefix, user.Login); err != nil {
		if !sgerrors.IsNotFound(err) {
			return errors.Wrap(err, "user get")
		}
	}
	err = s.repository.Put(ctx, prefix, user.Login, data)
	return err
}

func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	rawJSON, err := s.repository.Get(ctx, prefix, username)
	if err != nil {
		//If user doesn't exists we still want Forbidden instead of Not Found
		if sgerrors.IsNotFound(err) {
			return sgerrors.ErrInvalidCredentials
		}
		return err
	}
	user := new(User)
	if err = json.Unmarshal(rawJSON, user); err != nil {
		return errors.Wrap(err, "user authenticate unmarshall user")
	}

	if err := bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password)); err != nil {
		return sgerrors.ErrInvalidCredentials
	}
	return nil
}
