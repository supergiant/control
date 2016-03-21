package model

import "path"

type AppResource struct {
	c *Client
}

// implements Model interface (currently no methods)
type App struct {
	r    *AppResource
	Name string `json:"name"`
}

type AppList struct {
	Items []*App
}

// Model
//==============================================================================

// Relations--------------------------------------------------------------------
func (m *App) Components() *ComponentResource {
	return &ComponentResource{c: m.r.c, App: m}
}

// Operations-------------------------------------------------------------------
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
func (r *AppResource) EtcdKey(modelName string) string {
	return path.Join("/apps", modelName)
}

func (r *AppResource) List() (*AppList, error) {
	list := new(AppList)
	err := r.c.DB.List(r, list)
	return list, err
}

func (r *AppResource) Create(m *App) (*App, error) {
	err := r.c.DB.Create(r, m.Name, m)
	return m, err
}

func (r *AppResource) Get(name string) (m *App, err error) {
	err = r.c.DB.Get(r, name, m)
	return m, err
}

// No update for App

func (r *AppResource) Delete(name string) error {
	return r.c.DB.Delete(r, name)
}
