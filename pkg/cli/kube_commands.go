package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListKubes(c *cli.Context) (err error) {
	list := new(model.KubeList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).Kubes.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateKube(c *cli.Context) error {
	item := new(model.Kube)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Kubes.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetKube(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Kube)
	if err := newClient(c).Kubes.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateKube(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Kube)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Kubes.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteKube(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Kube)
	if err := newClient(c).Kubes.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
