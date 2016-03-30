package core

import (
	"fmt"
	"path"

	"github.com/supergiant/supergiant/types"
)

type ComponentCollection struct {
	core *Core
	App  *AppResource
}

type ComponentResource struct {
	collection *ComponentCollection
	*types.Component
}

type ComponentList struct {
	Items []*ComponentResource `json:"items"`
}

// EtcdKey implements the Collection interface.
func (c *ComponentCollection) EtcdKey(name types.ID) string {
	key := path.Join("/components", *c.App.Name)
	if name != nil {
		key = path.Join(key, *name)
	}
	return key
}

// InitializeResource implements the Collection interface.
func (c *ComponentCollection) InitializeResource(r Resource) {
	resource := r.(*ComponentResource)
	resource.collection = c

	// TODO it seems wrong this is called here -- execessive to have to load the
	// current Release, Entrypoints, and Kube Services just to render a
	// Component.
	// However, it's rare a Component is loaded out of the context of its
	// Release. We will change this when we see issues.
	if err := resource.decorate(); err != nil {
		panic(err)
	}
}

// List returns an ComponentList.
func (c *ComponentCollection) List() (*ComponentList, error) {
	list := new(ComponentList)
	err := c.core.DB.List(c, list)
	return list, err
}

// New initializes an Component with a pointer to the Collection.
func (c *ComponentCollection) New() *ComponentResource {
	return new(ComponentResource)
}

// Create takes an Component and creates it in etcd.
func (c *ComponentCollection) Create(r *ComponentResource) (*ComponentResource, error) {
	if err := c.core.DB.Create(c, r.Name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Get takes a name and returns an ComponentResource if it exists.
func (c *ComponentCollection) Get(name types.ID) (*ComponentResource, error) {
	r := c.New()
	if err := c.core.DB.Get(c, name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================

func (r *ComponentResource) decorate() error {
	release, err := r.CurrentRelease()
	if err != nil {
		return err
	}
	if release == nil {
		return nil
	}
	r.Addresses = &types.ComponentAddresses{
		External: release.ExternalAddresses(),
		Internal: release.InternalAddresses(),
	}
	return nil
}

// PersistableObject satisfies the Resource interface
func (r *ComponentResource) PersistableObject() interface{} {
	return r.PersistableComponent
}

// Save saves the Component in etcd through an update.
func (r *ComponentResource) Save() error {
	return r.collection.core.DB.Update(r.collection, r.Name, r)
}

// Delete cascades delete calls to current and target releases, and deletes the
// Component in etcd.
//
// TODO this should somehow stop any ongoing tasks related to the Component.
func (r *ComponentResource) Delete() error {

	// TODO we should really be going through and deleting all instances... it would just be a lot of requests

	current, err := r.CurrentRelease()
	if current != nil {
		// TODO should do something more formal here...
		fmt.Println(err)
		current.Delete()
	}
	target, err := r.TargetRelease()
	if target != nil {
		// TODO
		fmt.Println(err)
		target.Delete()
	}
	return r.collection.core.DB.Delete(r.collection, r.Name)
}

func (r *ComponentResource) App() *AppResource {
	return r.collection.App
}

func (r *ComponentResource) Releases() *ReleaseCollection {
	return &ReleaseCollection{
		core:      r.collection.core,
		Component: r,
	}
}

// TODO starting to think this should return err if it doesn't exist. We should
// expect the user to have checked if the ID is present.
func (r *ComponentResource) CurrentRelease() (*ReleaseResource, error) {
	if r.CurrentReleaseTimestamp == nil { // not yet released
		return nil, nil
	}
	return r.Releases().Get(r.CurrentReleaseTimestamp)
}

func (r *ComponentResource) TargetRelease() (*ReleaseResource, error) {
	if r.TargetReleaseTimestamp == nil { // something probably went wrong...
		return nil, nil
	}
	return r.Releases().Get(r.TargetReleaseTimestamp)
}
