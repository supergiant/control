package client

import (
	"fmt"
	"path"
	"time"

	"github.com/supergiant/supergiant/types"
)

type Instance types.Instance

type InstanceCollection struct {
	client *Client

	App       *AppResource
	Component *ComponentResource
	Release   *ReleaseResource
}

type InstanceResource struct {
	collection *InstanceCollection
	*Instance
}

type InstanceList struct {
	Items []*InstanceResource
}

func (c *InstanceCollection) path() string {
	return path.Join("apps", *c.App.Name, "components", *c.Component.Name, "releases", *c.Release.Timestamp, "instances")
}

func (r *InstanceResource) path() string {
	// TODO instance ID should probably just be a string
	return path.Join(r.collection.path(), *r.ID)
}

// Collection-level
//==============================================================================
func (c *InstanceCollection) New(m *Instance) *InstanceResource {
	return &InstanceResource{c, m}
}

func (c *InstanceCollection) List() (*InstanceList, error) {
	list := new(InstanceList)
	if _, err := c.client.Get(c.path(), list); err != nil {
		return nil, err
	}

	// TODO
	// We need some way, like we do in core/, of initializing the collection
	// object on each deserialized resource in a list. With Get & Create, we call
	// New() which handles that. We don't call New() when Listing, though,
	// because the items are deserialized directly onto a containing List object.
	//
	// This is different than core/, because you have to first unmarshal before
	// you can iterate through the underlying resources.
	//
	// But this may be fine for now.
	for _, instance := range list.Items {
		instance.collection = c
	}
	return list, nil
}

// func (c *InstanceCollection) Create(m *Instance) (*InstanceResource, error) {
// 	r := c.New(m)
// 	if err := c.client.Post(c.path(), m, r.Instance); err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }

func (c *InstanceCollection) Get(id types.ID) (*InstanceResource, error) {
	m := &Instance{
		ID: id,
	}
	r := c.New(m)
	if err := r.Reload(); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================
// func (r *InstanceResource) Delete() (bool, error) {
// 	return r.collection.client.Delete(r.path())
// }

func (r *InstanceResource) Reload() error {
	_, err := r.collection.client.Get(r.path(), r.Instance)
	return err
}

func (r *InstanceResource) Start() error {
	return r.collection.client.Post(r.path()+"/start", nil, nil)
}

func (r *InstanceResource) Stop() error {
	return r.collection.client.Post(r.path()+"/stop", nil, nil)
}

func (r *InstanceResource) WaitForStarted() error {
	desc := fmt.Sprintf("Instance start: %s", r.Name)
	return waitFor(desc, 120*time.Second, 3*time.Second, func() (bool, error) {
		if err := r.Reload(); err != nil {
			return false, err
		}
		return r.Status == "STARTED", nil
	})
}

func (r *InstanceResource) WaitForStopped() error {
	desc := fmt.Sprintf("Instance stop: %s", r.Name)
	return waitFor(desc, 120*time.Second, 3*time.Second, func() (bool, error) {
		if err := r.Reload(); err != nil {
			return false, err
		}
		return r.Status == "STOPPED", nil
	})
}
