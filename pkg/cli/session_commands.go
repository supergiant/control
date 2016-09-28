package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListSessions(c *cli.Context) error {
	list := new(model.SessionList)
	if err := newClient(c).Sessions.List(list); err != nil {
		return err
	}
	return printObj(list)
}

func CreateSession(c *cli.Context) error {
	item := new(model.Session)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Sessions.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetSession(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Session)
	if err := newClient(c).Sessions.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateSession(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Session)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Sessions.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteSession(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Session)
	if err := newClient(c).Sessions.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
