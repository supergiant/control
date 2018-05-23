package user

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const prefix = "users/"

// ETCDRepository is an implementation of user.Repository
type ETCDRepository struct {
	keysAPI client.KeysAPI
}

// Get retrieves user from etcd store if not found returns nil
func (r *ETCDRepository) Get(ctx context.Context, login string) (*User, error) {
	resp, err := r.keysAPI.Get(ctx, prefix+login, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	usr := new(User)
	err = json.NewDecoder(strings.NewReader(resp.Node.Value)).Decode(usr)
	if err != nil {
		logrus.Warningf("corrupted data in etcd node %s", resp.Node.Key)
	}
	return usr, nil
}

// Create adds user to the etcd
func (r *ETCDRepository) Create(ctx context.Context, u *User) error {
	rawJSON, err := json.Marshal(u)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = r.keysAPI.Create(ctx, prefix+u.Login, string(rawJSON))
	return errors.WithStack(err)
}

// GetAll retrieves all users from the etcd store, returns empty slice if no users are present
func (r *ETCDRepository) GetAll(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	resp, err := r.keysAPI.Get(ctx, prefix, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		if client.IsKeyNotFound(err) {
			logrus.Warningf("No users present in the system!")
			return users, nil
		}
		return nil, errors.WithStack(err)
	}

	for _, v := range resp.Node.Nodes {
		usr := new(User)
		err = json.NewDecoder(strings.NewReader(v.Value)).Decode(&usr)
		if err != nil {
			logrus.Warningf("corrupted data in etcd node %s", v.Key)
			continue
		}
		users = append(users, *usr)
	}

	return users, nil
}
