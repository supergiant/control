package model

import (
	"fmt"
	"guber"
	"path"
)

type ComponentResource struct {
	c   *Client
	App *App // Component belongs_to App
}

type Component struct {
	r *ComponentResource

	Name      string `json:"name"`
	Instances int    `json:"instances"`
	// TODO kinda weird,
	// you choose a container that has the deploy file, and then reference it as a command
	CustomDeployScript *CustomDeployScript `json:"custom_deploy_script"`

	Blueprint *Blueprint `json:"blueprint"`

	CurrentReleaseID    int    `json:"current_release_id"`
	ActiveDeploymentID  string `json:"active_deployment_id"`
	StandbyDeploymentID string `json:"standby_deployment_id"`
}

type ComponentList struct {
	Items []*Component
}

// TODO implement...
type CustomDeployScript struct {
	Image   string `json:"image"`
	Command string `json:"command"`
}

// Model (Entity)
//==============================================================================

// Util-------------------------------------------------------------------------
func (m *Component) imageRepoNames() (repoNames []string) { // TODO convert Image into Value object w/ repo, image, version
	var uniqRepoNames map[string]bool
	for _, container := range m.Blueprint.Containers {
		repoName := container.ImageRepoName()
		if _, ok := uniqRepoNames[repoName]; ok {
			uniqRepoNames[repoName] = true
			repoNames = append(repoNames, repoName)
		}
	}
	return repoNames
}

// TODO make sub-method on container, extract guts of for loop here
func (m *Component) containerPorts(public bool) (ports []*Port) {
	for _, container := range m.Blueprint.Containers {
		for _, port := range container.Ports {
			if port.Public == public {
				ports = append(ports, port)
			}
		}
	}
	return ports
}

func (m *Component) externalPorts() []*Port {
	return m.containerPorts(true)
}

func (m *Component) internalPorts() []*Port {
	return m.containerPorts(false)
}

func (m *Component) externalServicePorts() (ports []*guber.ServicePort) {
	for _, port := range m.containerPorts(true) {
		ports = append(ports, port.AsKubeServicePort())
	}
	return ports
}

func (m *Component) internalServicePorts() (ports []*guber.ServicePort) {
	for _, port := range m.containerPorts(false) {
		ports = append(ports, port.AsKubeServicePort())
	}
	return ports
}

// HasExternalPorts checks for any public ports in any of the defined Containers
func (m *Component) HasExternalPorts() bool {
	return len(m.externalPorts()) > 0
}

// HasInternalPorts checks for any non-public ports in any of the defined Containers
func (m *Component) HasInternalPorts() bool {
	return len(m.internalPorts()) > 0
}

// Operations-------------------------------------------------------------------
func (m *Component) provisionService(name string, svcType string, svcPorts []*guber.ServicePort) error {
	if service, _ := m.r.c.K8S.Services(m.App().Name).Get(name); service != nil {
		return nil // already created
	}

	service := &guber.Service{
		Metadata: &guber.Metadata{
			Name: name,
		},
		Spec: &guber.ServiceSpec{
			Type: svcType,
			Selector: map[string]string{
				"deployment": m.ActiveDeploymentID,
			},
			Ports: m.externalServicePorts(),
		},
	}
	_, err := m.r.c.K8S.Services(m.App().Name).Create(service)
	return err
}

func (m *Component) ProvisionExternalService() error {
	serviceName := fmt.Sprintf("%s-public", m.Name)
	return m.provisionService(serviceName, "NodePort", m.externalServicePorts())
}

func (m *Component) ProvisionInternalService() error {
	serviceName := m.Name
	return m.provisionService(serviceName, "NodePort", m.internalServicePorts())
}

// NOTE it seems weird here, but "Provision" == "CreateUnlessExists"
func (m *Component) ProvisionSecrets() error {
	repos, err := m.ImageRepos()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		if err := m.App().ProvisionSecret(repo); err != nil {
			return err
		}
	}
	return nil
}

// Relations--------------------------------------------------------------------
func (m *Component) App() *App {
	return m.r.App
}

func (m *Component) ImageRepos() (repos []*ImageRepo, err error) { // Not returning ImageRepoResource, since they are defined before hand
	for _, repoName := range m.imageRepoNames() {
		repo, err := m.r.c.ImageRepos().Get(repoName)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// TODO naming inconsistencies for kube definitions of resources
// ImagePullSecrets returns repo names defined for Kube pods
func (m *Component) ImagePullSecrets() (pullSecrets []*guber.ImagePullSecret, err error) {
	repos, err := m.ImageRepos()
	if err != nil {
		return pullSecrets, err
	}

	for _, repo := range repos {
		pullSecrets = append(pullSecrets, repo.AsKubeImagePullSecret())
	}
	return pullSecrets, nil
}

func (m *Component) deployments() *DeploymentResource {
	return &DeploymentResource{c: m.r.c, Component: m}
}

func (m *Component) ActiveDeployment() (*Deployment, error) {
	return m.deployments().Get(m.ActiveDeploymentID)
}

// Resource (Collection)
//==============================================================================
func (r *ComponentResource) EtcdKey(modelName string) string {
	return path.Join("/components", r.App.Name, modelName)
}

func (r *ComponentResource) List() (list *ComponentList, err error) {
	err = r.c.DB.List(r, list)
	return list, err
}

func (r *ComponentResource) Create(m *Component) (*Component, error) {
	err := r.c.DB.Create(r, m.Name, m)
	return m, err
}

func (r *ComponentResource) Get(name string) (m *Component, err error) {
	err = r.c.DB.Get(r, name, m)
	return m, err
}

// No update for Component (yet)

func (r *ComponentResource) Delete(name string) error {
	return r.c.DB.Delete(r, name)
}
