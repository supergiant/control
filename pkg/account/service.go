package account

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/model"
	"encoding/pem"
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
	logrus.Debug("cloud_account.Service.GetAll start")

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

	logrus.Debug("cloud_account.Service.GetAll end")
	return accounts, nil
}

// Get retrieves a user by it's accountName, returns nil if not found
func (s *Service) Get(ctx context.Context, accountName string) (*model.CloudAccount, error) {
	logrus.Debug("cloud_account.Service.Get start")

	res, err := s.repository.Get(ctx, s.storagePrefix, accountName)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, sgerrors.ErrNotFound
	}

	ca := new(model.CloudAccount)
	err = json.NewDecoder(bytes.NewReader(res)).Decode(ca)
	if err != nil {
		logrus.Warning("failed to convert stored data to cloud acccount struct")
		return nil, errors.WithStack(err)
	}

	logrus.Debug("cloud_account.Service.Get end")
	return ca, nil
}

// Create stores user in the underlying storage
func (s *Service) Create(ctx context.Context, account *model.CloudAccount) error {
	logrus.Debug("cloud_account.Service.Create start")

	rawJSON, err := json.Marshal(account)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.repository.Put(ctx, s.storagePrefix, account.Name, rawJSON)

	logrus.Debug("cloud_account.Service.Create end")
	return err
}

// Update cloud account
func (s *Service) Update(ctx context.Context, account *model.CloudAccount) error {
	logrus.Debug("cloud_account.Service.Update start")

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

	logrus.Debug("cloud_account.Service.Update end")
	return err
}

// Delete cloud account by name
func (s *Service) Delete(ctx context.Context, accountName string) error {
	logrus.Debug("cloud_account.Service.Delete start")
	err := s.repository.Delete(ctx, s.storagePrefix, accountName)
	logrus.Debug("cloud_account.Service.Delete end")
	return err
}
