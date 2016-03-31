package core

import (
	"path"

	"github.com/supergiant/guber"
	"github.com/supergiant/supergiant/types"
)

type AppCollection struct {
	core *Core
}

type AppResource struct {
	collection *AppCollection
	*types.App
}

// NOTE this does not inherit from types like model does; all we need is a List
// object, internally, that has a slice of our composed model above.
type AppList struct {
	Items []*AppResource `json:"items"`
}

// EtcdKey implements the Collection interface.
func (c *AppCollection) EtcdKey(name types.ID) string {
	if name == nil {
		return "/apps"
	}
	return path.Join("/apps", *name)
}

// InitializeResource implements the Collection interface.
func (c *AppCollection) InitializeResource(r Resource) {
	resource := r.(*AppResource)
	resource.collection = c
}

// List returns an AppList.
func (c *AppCollection) List() (*AppList, error) {
	list := new(AppList)
	err := c.core.DB.List(c, list)
	return list, err
}

// New initializes an App with a pointer to the Collection.
func (c *AppCollection) New() *AppResource {
	return &AppResource{
		App: &types.App{
			Meta: types.NewMeta(),
		},
	}
}

// Create takes an App and creates it in etcd. It also creates a Kubernetes
// Namespace with the name of the App.
func (c *AppCollection) Create(r *AppResource) (*AppResource, error) {
	if err := c.core.DB.Create(c, r.Name, r); err != nil {
		return nil, err
	}

	// TODO for error handling and retries, we may want to do this in a task and
	// utilize a Status field
	if err := r.createNamespace(); err != nil {
		panic(err)
	}
	return r, nil
}

// Get takes a name and returns an AppResource if it exists.
func (c *AppCollection) Get(name types.ID) (*AppResource, error) {
	r := c.New()
	if err := c.core.DB.Get(c, name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================

// PersistableObject satisfies the Resource interface
func (r *AppResource) PersistableObject() interface{} {
	return r.App
}

// Delete cascades deletes to all Components, deletes the Kube Namespace, and
// deletes the App in etcd.
func (r *AppResource) Delete() error {
	components, err := r.Components().List()
	if err != nil {
		return err
	}
	for _, component := range components.Items {
		if err := component.Delete(); err != nil {
			return err
		}
	}
	if err := r.deleteNamespace(); err != nil {
		return err
	}
	return r.collection.core.DB.Delete(r.collection, r.Name)
}

// Components returns a ComponentCollection with a pointer to the AppResource.
func (r *AppResource) Components() *ComponentCollection {
	return &ComponentCollection{
		core: r.collection.core,
		App:  r,
	}
}

func (r *AppResource) createNamespace() error {
	namespace := &guber.Namespace{
		Metadata: &guber.Metadata{
			Name: *r.Name,
		},
	}
	_, err := r.collection.core.K8S.Namespaces().Create(namespace)
	return err
}

func (r *AppResource) deleteNamespace() error {
	_, err := r.collection.core.K8S.Namespaces().Delete(*r.Name)
	return err
}

func (r *AppResource) ProvisionSecret(repo *ImageRepoResource) error {
	// TODO not sure i've been consistent with error handling -- this strategy is
	// useful when there could be multiple types of errors, alongside the
	// expectation of an error when something doesn't exist
	secret, err := r.collection.core.K8S.Secrets(*r.Name).Get(*repo.Name)

	if err != nil {
		return err
	} else if secret != nil {
		return nil
	}
	_, err = r.collection.core.K8S.Secrets(*r.Name).Create(asKubeSecret(repo))
	return err
}
