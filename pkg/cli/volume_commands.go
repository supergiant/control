package cli

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func ListVolumes(c *cli.Context) (err error) {
	list := new(model.VolumeList)
	list.Filters, err = listFilters(c)
	if err != nil {
		return err
	}
	if err = newClient(c).Volumes.List(list); err != nil {
		return err
	}
	return printList(c, list)
}

func CreateVolume(c *cli.Context) error {
	item := new(model.Volume)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Volumes.Create(item); err != nil {
		return err
	}
	return printObj(item)
}

func GetVolume(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Volume)
	if err := newClient(c).Volumes.Get(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func UpdateVolume(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Volume)
	if err := decodeInputFileInto(c, item); err != nil {
		return err
	}
	if err := newClient(c).Volumes.Update(&id, item); err != nil {
		return err
	}
	return printObj(item)
}

func DeleteVolume(c *cli.Context) error {
	id := c.Int64("id")
	item := new(model.Volume)
	if err := newClient(c).Volumes.Delete(&id, item); err != nil {
		return err
	}
	return printObj(item)
}
