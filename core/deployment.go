package core

import (
	"math/rand"
	"path"
	"time"
)

const (
	// Deployment IDs are short, since they are used in concatenated asset names
	idChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	idLen   = 4
)

type DeploymentResource struct {
	c         *Client
	Component *Component
}

type Deployment struct {
	r  *DeploymentResource
	ID string `json:"type"`

	// These are needed, since Components and Deployments have a belongs_to style relation
	appName       string `json:"app"`
	componentName string `json:"component"`
}

type DeploymentList struct {
	Items []*Deployment `json:"items"`
}

func newDeploymentID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, idLen)
	for i := 0; i < idLen; i++ {
		result[i] = idChars[rand.Intn(len(idChars))]
	}
	return string(result)
}

// Model (Entity)
//==============================================================================

// Relations--------------------------------------------------------------------
func (m *Deployment) Component() *Component {
	return m.r.Component
}

func (m *Deployment) Instances() (instances []*Instance, err error) {
	for n := 1; n < m.Component().Instances+1; n++ {
		instance := &Instance{c: m.r.c, Deployment: m, ID: n}
		instances = append(instances, instance)
	}
	return instances, nil
}

// Resource (Collection)
//==============================================================================
func (r *DeploymentResource) EtcdKey(id string) string {
	return path.Join("/deployments", id)
}

// InitializeModel is a part of Resource interface
func (r *DeploymentResource) InitializeModel(m Model) {
	model := m.(*Deployment)
	model.r = r
}

func (r *DeploymentResource) New() *Deployment {
	return &Deployment{
		r:             r,
		ID:            newDeploymentID(),
		appName:       r.Component.App().Name,
		componentName: r.Component.Name,
	}
}

func (r *DeploymentResource) List() (*DeploymentList, error) {
	list := new(DeploymentList)
	err := r.c.DB.List(r, list)
	return list, err
}

func (r *DeploymentResource) Create(m *Deployment) (*Deployment, error) {
	err := r.c.DB.Create(r, m.ID, m)
	return m, err
}

func (r *DeploymentResource) Get(id string) (*Deployment, error) {
	m := r.New()
	err := r.c.DB.Get(r, id, m)
	return m, err
}

// No update for Deployment

func (r *DeploymentResource) Delete(id string) error {
	return r.c.DB.Delete(r, id)
}
