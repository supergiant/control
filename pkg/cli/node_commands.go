package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListNodes(c *cli.Context) (err error) {
	list := new(model.NodeList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).Nodes.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateNode(c *cli.Context) error {
	item := new(model.Node)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Nodes.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetNode(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Node)
	if err := newClient(c).Nodes.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateNode(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Node)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Nodes.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteNode(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Node)
	if err := newClient(c).Nodes.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
