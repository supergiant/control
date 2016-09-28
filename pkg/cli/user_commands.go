package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListUsers(c *cli.Context) (err error) {
	list := new(model.UserList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).Users.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateUser(c *cli.Context) error {
	item := new(model.User)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Users.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetUser(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.User)
	if err := newClient(c).Users.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateUser(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.User)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Users.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteUser(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.User)
	if err := newClient(c).Users.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
