package cloud_account

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const prefix = "cloudaccount/"

// ETCDRepository is implementation of cloud_account.Repository
type ETCDRepository struct {
	keysAPI client.KeysAPI
}

// NewRepository is function constructor
func NewRepository(cl client.Client) Repository {
	return &ETCDRepository{
		keysAPI: client.NewKeysAPI(cl),
	}
}

// GetAll - retrieves all accounts, returns empty slice if none
func (r *ETCDRepository) GetAll(ctx context.Context) ([]CloudAccount, error) {
	accounts := make([]CloudAccount, 0)
	resp, err := r.keysAPI.Get(ctx, prefix, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		if client.IsKeyNotFound(err) {
			return accounts, nil
		}
		return accounts, errors.WithStack(err)
	}
	if resp.Node != nil {
		for _, v := range resp.Node.Nodes {
			ca := CloudAccount{}
			err = json.NewDecoder(strings.NewReader(v.Value)).Decode(&ca)
			if err != nil {
				logrus.Warningf("corrupted data in etcd node %s", v.Key)
				continue
			}
			accounts = append(accounts, ca)
		}
	}
	return accounts, nil
}

// Create stores account in etcd, if account with such name is already present it will return an error
func (r *ETCDRepository) Create(ctx context.Context, acc *CloudAccount) error {
	rawJSON, err := json.Marshal(acc)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = r.keysAPI.Create(ctx, prefix+acc.Name, string(rawJSON))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Get retrieves account from etcd returns nil if account is not found
func (r *ETCDRepository) Get(ctx context.Context, accountName string) (*CloudAccount, error) {
	resp, err := r.keysAPI.Get(ctx, prefix+accountName, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	ac := &CloudAccount{}
	err = json.NewDecoder(strings.NewReader(resp.Node.Value)).Decode(ac)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ac, nil
}

// Update overwrites the cloud account record, account name should not be changed.
func (r *ETCDRepository) Update(ctx context.Context, acc *CloudAccount) error {
	rawJSON, err := json.Marshal(acc)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = r.keysAPI.Set(ctx, prefix+acc.Name, string(rawJSON), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete removes amn account from the etcd storage, this is idempotent operation.
func (r *ETCDRepository) Delete(ctx context.Context, accountName string) error {
	_, err := r.keysAPI.Delete(ctx, prefix+accountName, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			return nil
		}
		return errors.WithStack(err)
	}
	return nil
}
