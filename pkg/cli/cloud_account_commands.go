package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListCloudAccounts(c *cli.Context) (err error) {
	list := new(model.CloudAccountList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).CloudAccounts.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateCloudAccount(c *cli.Context) error {
	item := new(model.CloudAccount)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).CloudAccounts.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetCloudAccount(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.CloudAccount)
	if err := newClient(c).CloudAccounts.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateCloudAccount(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.CloudAccount)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).CloudAccounts.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteCloudAccount(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.CloudAccount)
	if err := newClient(c).CloudAccounts.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
