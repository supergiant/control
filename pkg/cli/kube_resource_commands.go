package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListKubeResources(c *cli.Context) error {
	list := new(model.KubeResourceList)
	if err := newClient(c).KubeResources.List(list); err != nil {
		return err
	}
	return printObj(list)
}

func CreateKubeResource(c *cli.Context) error {
	item := new(model.KubeResource)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).KubeResources.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetKubeResource(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.KubeResource)
	if err := newClient(c).KubeResources.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateKubeResource(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.KubeResource)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).KubeResources.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteKubeResource(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.KubeResource)
	if err := newClient(c).KubeResources.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func StartKubeResource(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.KubeResource)
	if err := newClient(c).KubeResources.Start(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func StopKubeResource(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.KubeResource)
	if err := newClient(c).KubeResources.Stop(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
