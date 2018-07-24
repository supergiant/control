package user

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
)

const DefaultStoragePrefix = "/supergiant/user/"

// Service contains business logic related to users
type Service struct {
	storagePrefix string
	repository    storage.Interface
}

func NewService(storagePrefix string, repository storage.Interface) *Service {
	return &Service{
		storagePrefix: storagePrefix,
		repository:    repository,
	}
}

func (s *Service) Create(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user create: user can't be nil")
	}
	err := user.encryptPassword()
	if err != nil {
		return err
	}
	data, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "user create")
	}
	if _, err := s.repository.Get(ctx, s.storagePrefix, user.Login); err != nil {
		if !sgerrors.IsNotFound(err) {
			return errors.Wrap(err, "user get")
		}
	}
	err = s.repository.Put(ctx, s.storagePrefix, user.Login, data)
	return err
}

func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return sgerrors.ErrInvalidCredentials
	}

	rawJSON, err := s.repository.Get(ctx, s.storagePrefix, username)
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
