package core

import (
	"fmt"
	"path"

	"github.com/supergiant/guber"
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
	Items []*Component `json:"items"`
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
	uniqRepoNames := make(map[string]bool)
	for _, container := range m.Blueprint.Containers {
		repoName := container.ImageRepoName()
		if _, ok := uniqRepoNames[repoName]; !ok {
			uniqRepoNames[repoName] = true
			repoNames = append(repoNames, repoName)
		}
	}
	return repoNames
}

// TODO make sub-method on container, extract guts of for loop here
func (m *Component) containerPorts(public bool) (ports []*Port) {

	// TODO these will need to be unique -------------------------------------------------

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

func (m *Component) hasExternalPorts() bool {
	return len(m.externalPorts()) > 0
}

func (m *Component) hasInternalPorts() bool {
	return len(m.internalPorts()) > 0
}

// Operations-------------------------------------------------------------------
func (m *Component) getService(name string) (*guber.Service, error) {
	return m.r.c.K8S.Services(m.App().Name).Get(name)
}

func (m *Component) provisionService(name string, svcType string, svcPorts []*guber.ServicePort) error {
	if service, _ := m.getService(name); service != nil {
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
			Ports: svcPorts,
		},
	}
	_, err := m.r.c.K8S.Services(m.App().Name).Create(service)
	return err
}

func (m *Component) externalServiceName() string {
	return fmt.Sprintf("%s-public", m.Name)
}

func (m *Component) internalServiceName() string {
	return m.Name
}

// Exposed to fetch IPs
func (m *Component) ExternalService() (*guber.Service, error) {
	return m.getService(m.externalServiceName())
}

func (m *Component) InternalService() (*guber.Service, error) {
	return m.getService(m.internalServiceName())
}

func (m *Component) provisionExternalService() error {
	return m.provisionService(m.externalServiceName(), "NodePort", m.externalServicePorts())
}

func (m *Component) provisionInternalService() error {
	return m.provisionService(m.internalServiceName(), "ClusterIP", m.internalServicePorts())
}

func (m *Component) destroyServices() (err error) {
	if _, err = m.r.c.K8S.Services(m.App().Name).Delete(m.externalServiceName()); err != nil {
		return err
	}
	if _, err = m.r.c.K8S.Services(m.App().Name).Delete(m.internalServiceName()); err != nil {
		return err
	}
	return nil
}

// NOTE it seems weird here, but "Provision" == "CreateUnlessExists"
func (m *Component) provisionSecrets() error {
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

func (m *Component) Provision() error {
	if err := m.provisionSecrets(); err != nil {
		return err
	}

	// Create Services
	if m.hasInternalPorts() {
		if err := m.provisionInternalService(); err != nil {
			return err
		}
	}
	if m.hasExternalPorts() {
		if err := m.provisionExternalService(); err != nil {
			return err
		}
	}

	// NOTE the code below is not tucked inside `deployment.ProvisionInstances()`
	// because I predict eventually wanting to record % progress for each instance
	// without having to expect Deployment to handle progress logic as well.

	// Concurrently provision instances
	deployment, err := m.ActiveDeployment()
	if err != nil {
		return err
	}
	instances, err := deployment.Instances()
	if err != nil {
		return err
	}

	c := make(chan error)
	for _, instance := range instances {
		go func(instance *Instance) { // NOTE we have to pass instance here, else every goroutine hits the same instance
			c <- instance.Provision()
		}(instance)
	}
	for i := 0; i < m.Instances; i++ {
		if err := <-c; err != nil {
			return err
		}
	}
	return nil
}

func (m *Component) Teardown() error { // this is named Teardown() and not Destroy() because we need a method to actually kick off the DestroyComponent job
	if err := m.destroyServices(); err != nil {
		return err
	}

	// TODO this needs to include destrution logic for ALL deployments, if there happens to be more than 1 at the time
	deployment, err := m.ActiveDeployment()
	if err != nil {
		return err
	}
	instances, err := deployment.Instances()
	if err != nil {
		return err
	}

	c := make(chan error)
	for _, instance := range instances {
		go func(instance *Instance) {
			c <- instance.Destroy()
		}(instance)
	}
	for i := 0; i < m.Instances; i++ {
		if err := <-c; err != nil {
			return err
		}
	}
	return nil
}

// NOTE the delete portion of this is handled by the job
func (m *Component) TeardownAndDelete() error {
	msg := &DestroyComponentMessage{
		AppName:       m.App().Name,
		ComponentName: m.Name,
	}
	_, err := m.r.c.Jobs().Start(JobTypeDestroyComponent, msg)
	return err
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

// InitializeModel is a part of Resource interface
func (r *ComponentResource) InitializeModel(m Model) {
	model := m.(*Component)
	model.r = r
}

func (r *ComponentResource) List() (*ComponentList, error) {
	list := new(ComponentList)
	err := r.c.DB.List(r, list)
	return list, err
}

func (r *ComponentResource) New() *Component {
	return &Component{r: r}
}

func (r *ComponentResource) Create(m *Component) (*Component, error) {
	deployment, err := m.deployments().Create(m.deployments().New())
	if err != nil {
		return nil, err
	}

	m.ActiveDeploymentID = deployment.ID

	if err := r.c.DB.Create(r, m.Name, m); err != nil {
		return nil, err
	}

	msg := &CreateComponentMessage{
		AppName:       m.App().Name,
		ComponentName: m.Name,
	}
	if _, err := r.c.Jobs().Start(JobTypeCreateComponent, msg); err != nil {

		// TODO error handling is not great here, considering the model has already been persisted
		return nil, err
	}

	return m, nil
}

func (r *ComponentResource) Get(name string) (*Component, error) {
	m := r.New()
	err := r.c.DB.Get(r, name, m)
	return m, err
}

// No update for Component (yet)

func (r *ComponentResource) Delete(name string) error {
	return r.c.DB.Delete(r, name)
}
