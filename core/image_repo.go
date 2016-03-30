package core

import (
	"path"

	"github.com/supergiant/supergiant/types"
)

type ImageRepoCollection struct {
	core *Core
}

type ImageRepoResource struct {
	collection *ImageRepoCollection
	*types.ImageRepo
}

// NOTE this does not inherit from types like model does; all we need is a List
// object, internally, that has a slice of our composed model above.
type ImageRepoList struct {
	Items []*ImageRepoResource `json:"items"`
}

// EtcdKey implements the Collection interface.
func (c *ImageRepoCollection) EtcdKey(name types.ID) string {
	key := "/image_repos/dockerhub"
	if name != nil {
		key = path.Join(key, *name)
	}
	return key
}

// InitializeResource implements the Collection interface.
func (c *ImageRepoCollection) InitializeResource(r Resource) {
	resource := r.(*ImageRepoResource)
	resource.collection = c
}

// List returns an ImageRepoList.
func (c *ImageRepoCollection) List() (*ImageRepoList, error) {
	list := new(ImageRepoList)
	err := c.core.DB.List(c, list)
	return list, err
}

// New initializes an ImageRepo with a pointer to the Collection.
func (c *ImageRepoCollection) New() *ImageRepoResource {
	return new(ImageRepoResource)
}

// Create takes an ImageRepo and creates it in etcd. It also creates a Kubernetes
// Namespace with the name of the ImageRepo.
func (c *ImageRepoCollection) Create(r *ImageRepoResource) (*ImageRepoResource, error) {
	if err := c.core.DB.Create(c, r.Name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Get takes a name and returns an ImageRepoResource if it exists.
func (c *ImageRepoCollection) Get(name types.ID) (*ImageRepoResource, error) {
	r := c.New()
	if err := c.core.DB.Get(c, name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================

// PersistableObject satisfies the Resource interface
func (r *ImageRepoResource) PersistableObject() interface{} {
	return r.ImageRepo
}

// Delete deletes the ImageRepo in etcd.
func (r *ImageRepoResource) Delete() error {
	return r.collection.core.DB.Delete(r.collection, r.Name)
}
