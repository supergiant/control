package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListEntrypoints(c *cli.Context) (err error) {
	list := new(model.EntrypointList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).Entrypoints.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateEntrypoint(c *cli.Context) error {
	item := new(model.Entrypoint)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Entrypoints.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetEntrypoint(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Entrypoint)
	if err := newClient(c).Entrypoints.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateEntrypoint(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Entrypoint)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Entrypoints.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteEntrypoint(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Entrypoint)
	if err := newClient(c).Entrypoints.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
