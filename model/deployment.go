package model

import "path"

type DeploymentResource struct {
	c         *Client
	Component *Component
}

type Deployment struct {
	r  *DeploymentResource
	ID string `json:"type"`
}

type DeploymentList struct {
	Items []*Deployment
}

// Model (Entity)
//==============================================================================

// Relations--------------------------------------------------------------------
func (m *Deployment) Component() *Component {
	return m.r.Component
}

func (m *Deployment) Instances() (instances []*Instance, err error) {
	for n := 1; n < m.Component().Instances+1; n++ {
		instances = append(instances, &Instance{c: m.r.c, Deployment: m})
	}
	return instances, nil
}

// Resource (Collection)
//==============================================================================
func (r *DeploymentResource) EtcdKey(id string) string {
	return path.Join("/deployments", id)
}

func (r *DeploymentResource) List() (list *DeploymentList, err error) {
	err = r.c.DB.List(r, list)
	return list, err
}

func (r *DeploymentResource) Create(m *Deployment) (*Deployment, error) {
	err := r.c.DB.Create(r, m.ID, m)
	return m, err
}

func (r *DeploymentResource) Get(id string) (m *Deployment, err error) {
	err = r.c.DB.Get(r, id, m)
	return m, err
}

// No update for Deployment

func (r *DeploymentResource) Delete(id string) error {
	return r.c.DB.Delete(r, id)
}
