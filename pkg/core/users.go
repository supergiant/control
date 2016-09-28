package core

import "github.com/supergiant/supergiant/pkg/model"

type Users struct {
	Collection
}

func (c *Users) RegenerateAPIToken(id *int64, m *model.User) error {
	m.ID = id
	m.GenerateAPIToken()
	return c.Core.DB.Model(m).Update("api_token", m.APIToken)
}
