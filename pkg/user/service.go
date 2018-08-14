package user

import (
	"context"

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

// NewService is a constructor function for user.Service
func NewService(storagePrefix string, repository storage.Interface) *Service {
	return &Service{
		storagePrefix: storagePrefix,
		repository:    repository,
	}
}

// Create is used to register new user
func (s *Service) Create(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user create: user can't be nil")
	}
	err := user.encryptPassword()
	if err != nil {
		return err
	}

	if _, err := s.repository.Get(ctx, s.storagePrefix, user.Login); err != nil {
		if !sgerrors.IsNotFound(err) {
			return errors.Wrap(err, "user get")
		}
	}
	err = s.repository.Put(ctx, s.storagePrefix, user.Login, user.ToJSON())
	return err
}

// Authenticate checks if password stored in db is the same as in request
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
	user, err := FromJSON(rawJSON)
	if err != nil {
		return errors.Wrap(err, "user authenticate: unmarshall user")
	}

	if err := bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password)); err != nil {
		return sgerrors.ErrInvalidCredentials
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]*User, error) {
	res, err := s.repository.GetAll(ctx, s.storagePrefix)
	if err != nil {
		return nil, err
	}

	usrs := make([]*User, 0)
	for _, v := range res {
		u, err := FromJSON(v)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, u)
	}
	return usrs, nil
}
