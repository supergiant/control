package fake_client

import "github.com/supergiant/supergiant/pkg/model"

type Users struct {
	Collection
	RegenerateAPITokenFn func(interface{}, *model.User) error
}

func (c *Users) RegenerateAPIToken(id interface{}, m *model.User) error {
	if c.RegenerateAPITokenFn == nil {
		return nil
	}
	return c.RegenerateAPITokenFn(id, m)
}
