package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListEntrypointListeners(c *cli.Context) (err error) {
	list := new(model.EntrypointListenerList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).EntrypointListeners.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateEntrypointListener(c *cli.Context) error {
	item := new(model.EntrypointListener)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).EntrypointListeners.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetEntrypointListener(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.EntrypointListener)
	if err := newClient(c).EntrypointListeners.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateEntrypointListener(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.EntrypointListener)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).EntrypointListeners.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteEntrypointListener(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.EntrypointListener)
	if err := newClient(c).EntrypointListeners.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
