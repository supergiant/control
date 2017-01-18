package core

import (
	"errors"

	"github.com/supergiant/supergiant/pkg/model"
)

type CloudAccounts struct {
	Collection
}

func (c *CloudAccounts) Create(m *model.CloudAccount) error {
	// NOTE we have to do pre-validation here in order to make sure provider is correct
	if err := validateFields(m); err != nil {
		return err
	}

	if err := c.provider(m).ValidateAccount(m); err != nil {
		return &ErrorValidationFailed{err}
	}
	return c.Collection.Create(m)
}

func (c *CloudAccounts) Delete(id *int64, m *model.CloudAccount) error {
	if err := c.Core.DB.Where("cloud_account_name = ?", m.Name).Find(&m.Kubes); err != nil {
		return err
	}
	if len(m.Kubes) > 0 {
		return &ErrorValidationFailed{errors.New("Cannot delete CloudAccount that has active Kubes")}
	}
	return c.Collection.Delete(id, m)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (c *CloudAccounts) provider(m *model.CloudAccount) Provider {
	switch m.Provider {
	case "aws":
		return c.Core.AWSProvider(m.Credentials)
	case "digitalocean":
		return c.Core.DOProvider(m.Credentials)
	case "openstack":
		return c.Core.OSProvider(m.Credentials)
	case "gce":
		return c.Core.GCEProvider(m.Credentials)
	default:
		panic("Could not load provider interface for " + m.Provider)
	}
}
