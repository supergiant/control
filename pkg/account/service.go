package account

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/supergiant/supergiant/pkg/clouds"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
)

// Service holds all business logic related to cloud accounts
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

const DefaultStoragePrefix = "/supergiant/account/"

// GetAll retrieves cloud accounts from underlying storage, returns empty slice if none found
func (s *Service) GetAll(ctx context.Context) ([]model.CloudAccount, error) {

	accounts := make([]model.CloudAccount, 0)
	res, err := s.repository.GetAll(ctx, s.storagePrefix)
	if err != nil {
		return accounts, err
	}
	for _, v := range res {
		ca := new(model.CloudAccount)
		err = json.NewDecoder(bytes.NewReader(v)).Decode(ca)
		if err != nil {
			logrus.Warningf("failed to convert stored data to cloud account struct")
			logrus.Debugf("corrupted data: %s", string(v))
			continue
		}
		accounts = append(accounts, *ca)
	}

	return accounts, nil
}

// Get retrieves a user by it's accountName, returns nil if not found
func (s *Service) Get(ctx context.Context, accountName string) (*model.CloudAccount, error) {

	res, err := s.repository.Get(ctx, s.storagePrefix, accountName)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, sgerrors.ErrNotFound
	}

	ca := &model.CloudAccount{}
	err = json.NewDecoder(bytes.NewReader(res)).Decode(ca)
	if err != nil {
		logrus.Warning("failed to convert stored data to cloud acccount struct")
		return nil, errors.WithStack(err)
	}

	if ca.Credentials == nil {
		ca.Credentials = make(map[string]string, 0)
	}

	return ca, nil
}

// Create stores user in the underlying storage
func (s *Service) Create(ctx context.Context, account *model.CloudAccount) error {
	switch account.Provider {
	case clouds.DigitalOcean:
		if account.Credentials[clouds.DigitalOceanAccessToken] == "" {
			return errors.Wrap(sgerrors.ErrInvalidCredentials, "no digital ocean's access token provided")
		}
	case clouds.AWS:
		if account.Credentials[clouds.AWSAccessKeyID] == "" ||
			account.Credentials[clouds.AWSSecretKey] == "" {
			return errors.Wrap(sgerrors.ErrInvalidCredentials, "both aws access key and key id should be provided")
		}
	default:
		return sgerrors.ErrUnsupportedProvider
	}

	rawJSON, err := json.Marshal(account)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.repository.Put(ctx, s.storagePrefix, account.Name, rawJSON)

	return err
}

// Update cloud account
func (s *Service) Update(ctx context.Context, account *model.CloudAccount) error {
	rawJSON, err := json.Marshal(account)
	if err != nil {
		return errors.WithStack(err)
	}

	oldAcc, err := s.Get(ctx, account.Name)
	if err != nil {
		return err
	}
	if oldAcc.Name != account.Name || oldAcc.Provider != account.Provider {
		return errors.New("account name or provider can't be changed")
	}

	err = s.repository.Put(ctx, s.storagePrefix, account.Name, rawJSON)

	return err
}

// Delete cloud account by name
func (s *Service) Delete(ctx context.Context, accountName string) error {
	return s.repository.Delete(ctx, s.storagePrefix, accountName)
}
