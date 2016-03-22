package core

import (
	"fmt"
	"guber"
	"path"
)

type AppResource struct {
	c *Client
}

// implements Model interface (currently no methods)
type App struct {
	r    *AppResource
	Name string `json:"name"`
}

type AppList struct {
	Items []*App `json:"items"`
}

// Model
//==============================================================================

// Relations--------------------------------------------------------------------
func (m *App) Components() *ComponentResource {
	return &ComponentResource{c: m.r.c, App: m}
}

// Operations-------------------------------------------------------------------
func (m *App) createNamespace() error {

	fmt.Println("app", *m)

	namespace := &guber.Namespace{
		Metadata: &guber.Metadata{
			Name: m.Name,
		},
	}

	fmt.Println("namespace", *namespace)

	_, err := m.r.c.K8S.Namespaces().Create(namespace)
	return err
}

func (m *App) deleteNamespace() error {
	_, err := m.r.c.K8S.Namespaces().Delete(m.Name)
	return err
}

func (m *App) ProvisionSecret(repo *ImageRepo) error {
	// TODO not sure i've been consistent with error handling -- this strategy is
	// useful when there could be multiple types of errors, alongside the
	// expectation of an error when something doesn't exist
	secret, err := m.r.c.K8S.Secrets(m.Name).Get(repo.Name)
	if secret != nil {
		return nil
	} else if err != nil {
		return err
	}
	_, err = m.r.c.K8S.Secrets(m.Name).Create(repo.AsKubeSecret())
	return err
}

// Resource (Collection)
//==============================================================================

// EtcdKey is a part of Resource interface
func (r *AppResource) EtcdKey(modelName string) string {
	return path.Join("/apps", modelName)
}

// InitializeModel is a part of Resource interface
func (r *AppResource) InitializeModel(m Model) {
	model := m.(*App)
	model.r = r
}

func (r *AppResource) List() (*AppList, error) {
	list := new(AppList)
	err := r.c.DB.List(r, list)
	return list, err
}

func (r *AppResource) New() *App {
	return &App{r: r}
}

func (r *AppResource) Create(m *App) (*App, error) {
	if err := r.c.DB.Create(r, m.Name, m); err != nil {
		return nil, err
	}

	// TODO for error handling and retries, we may want to do this in a job and
	// utilize a Status field
	if err := m.createNamespace(); err != nil {
		panic(err) // TODO see above
	}

	return m, nil
}

func (r *AppResource) Get(name string) (*App, error) {
	m := r.New()
	err := r.c.DB.Get(r, name, m)
	return m, err
}

// No update for App

func (r *AppResource) Delete(name string) error {
	// NOTE we do a Get here to verify the record exists before trying to delete
	// dependencies, and also to be able to call model-level methods.
	m, err := r.Get(name)
	if err != nil {
		return err
	}

	if err = m.deleteNamespace(); err != nil {
		panic(err) // TODO see Create method about using jobs/status field for this
	}

	return r.c.DB.Delete(r, name)
}
