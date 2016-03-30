package client

import (
	"path"

	"github.com/supergiant/supergiant/types"
)

type Component types.Component

type ComponentCollection struct {
	client *Client

	App *AppResource
}

type ComponentResource struct {
	collection *ComponentCollection
	*Component
}

type ComponentList struct {
	Items []*ComponentResource
}

func (c *ComponentCollection) path() string {
	return path.Join("apps", *c.App.Name, "components")
}

func (r *ComponentResource) path() string {
	return path.Join(r.collection.path(), *r.Name)
}

// Collection-level
//==============================================================================
func (c *ComponentCollection) New(m *Component) *ComponentResource {
	return &ComponentResource{c, m}
}

func (c *ComponentCollection) List() (*ComponentList, error) {
	list := new(ComponentList)
	if _, err := c.client.Get(c.path(), list); err != nil {
		return nil, err
	}
	// see TODO in instance.go
	for _, component := range list.Items {
		component.collection = c
	}
	return list, nil
}

func (c *ComponentCollection) Create(m *Component) (*ComponentResource, error) {
	r := c.New(m)
	if err := c.client.Post(c.path(), m, r.Component); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *ComponentCollection) Get(name types.ID) (*ComponentResource, error) {
	m := &Component{
		PersistableComponent: &types.PersistableComponent{ // TODO any way to not make this so weird?
			Name: name,
		},
	}
	r := c.New(m)
	if found, err := c.client.Get(r.path(), r.Component); err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return r, nil
}

// Resource-level
//==============================================================================
// func (r *ComponentResource) Update(m *Component) (*ComponentResource, error) {
//   r.collection.client.
// }

func (r *ComponentResource) Delete() (bool, error) {
	return r.collection.client.Delete(r.path())
}

// Relations
func (r *ComponentResource) Releases() *ReleaseCollection {
	return &ReleaseCollection{
		client:    r.collection.client,
		App:       r.collection.App,
		Component: r,
	}
}

func (r *ComponentResource) CurrentRelease() (*ReleaseResource, error) {
	if r.CurrentReleaseTimestamp == nil { // will be empty on first release
		return nil, nil
	}
	return r.Releases().Get(r.CurrentReleaseTimestamp)
}

func (r *ComponentResource) TargetRelease() (*ReleaseResource, error) {
	if r.TargetReleaseTimestamp == nil {
		return nil, nil
	}
	return r.Releases().Get(r.TargetReleaseTimestamp)
}
