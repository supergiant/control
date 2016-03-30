package client

import (
	"path"

	"github.com/supergiant/supergiant/types"
)

type Release types.Release

type ReleaseCollection struct {
	client *Client

	App       *AppResource
	Component *ComponentResource
}

type ReleaseResource struct {
	collection *ReleaseCollection
	*Release
}

type ReleaseList struct {
	Items []*ReleaseResource
}

func (c *ReleaseCollection) path() string {
	return path.Join("apps", *c.App.Name, "components", *c.Component.Name, "releases")
}

func (r *ReleaseResource) path() string {
	return path.Join(r.collection.path(), *r.Timestamp)
}

// Collection-level
//==============================================================================
func (c *ReleaseCollection) New(m *Release) *ReleaseResource {
	return &ReleaseResource{c, m}
}

func (c *ReleaseCollection) List() (*ReleaseList, error) {
	list := new(ReleaseList)
	if _, err := c.client.Get(c.path(), list); err != nil {
		return nil, err
	}
	// see TODO in instance.go
	for _, release := range list.Items {
		release.collection = c
	}
	return list, nil
}

// func (c *ReleaseCollection) Create(m *Release) (*ReleaseResource, error) {
// 	r := c.New(m)
// 	if err := c.client.Post(c.path(), m, r.Release); err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }

func (c *ReleaseCollection) Get(timestamp types.ID) (*ReleaseResource, error) {
	m := &Release{
		Timestamp: timestamp,
	}
	r := c.New(m)
	if found, err := c.client.Get(r.path(), r.Release); err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return r, nil
}

// Resource-level
//==============================================================================
func (r *ReleaseResource) Delete() (bool, error) {
	return r.collection.client.Delete(r.path())
}

// Relations
func (r *ReleaseResource) Instances() *InstanceCollection {
	return &InstanceCollection{
		client:    r.collection.client,
		App:       r.collection.App,
		Component: r.collection.Component,
		Release:   r,
	}
}
