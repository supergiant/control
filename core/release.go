package core

import (
	"fmt"
	"path"
	"time"

	"github.com/supergiant/guber"
	"github.com/supergiant/supergiant/types"
)

type ReleaseCollection struct {
	core      *Core
	Component *ComponentResource
}

type ReleaseResource struct {
	collection *ReleaseCollection
	*types.Release
}

type ReleaseList struct {
	Items []*ReleaseResource `json:"items"`
}

// EtcdKey implements the Collection interface.
func (c *ReleaseCollection) EtcdKey(id string) string {
	return path.Join("/releases", c.Component.App().Name, c.Component.Name, id)
}

// InitializeResource implements the Collection interface.
func (c *ReleaseCollection) InitializeResource(r Resource) {
	resource := r.(*ReleaseResource)
	resource.collection = c
}

// List returns an ReleaseList.
func (c *ReleaseCollection) List() (*ReleaseList, error) {
	list := new(ReleaseList)
	err := c.core.DB.List(c, list)
	return list, err
}

// New initializes an Release with a pointer to the Collection.
func (c *ReleaseCollection) New() *ReleaseResource {
	return &ReleaseResource{
		collection: c,
		Release: &types.Release{
			ID: newReleaseID(),
		},
	}
}

// Create takes an Release and creates it in etcd.
func (c *ReleaseCollection) Create(r *ReleaseResource) (*ReleaseResource, error) {
	if err := c.core.DB.Create(c, r.ID, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Get takes an id and returns an ReleaseResource if it exists.
func (c *ReleaseCollection) Get(id string) (*ReleaseResource, error) {
	r := c.New()

	// TODO
	if id == "" {
		panic("id is nil")
	}

	if err := c.core.DB.Get(c, id, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================

// Delete removes all assets (volumes, pods, etc.) and deletes the Release in
// etcd.
func (r *ReleaseResource) Delete() error {
	if err := r.deleteServices(); err != nil {
		return err
	}

	c := make(chan error)
	for _, instance := range r.Instances().List().Items {
		go func(instance *InstanceResource) {
			c <- instance.Delete()
		}(instance)
	}
	for i := 0; i < r.InstanceCount; i++ {
		if err := <-c; err != nil {
			return err
		}
	}

	return r.collection.core.DB.Delete(r.collection, r.ID)
}

func newReleaseID() string {
	return time.Now().Format("20060102150405")
}

func (r *ReleaseResource) App() *AppResource {
	return r.Component().App()
}

func (r *ReleaseResource) Component() *ComponentResource {
	return r.collection.Component
}

func (r *ReleaseResource) Instances() *InstanceCollection {
	return &InstanceCollection{
		core:    r.collection.core,
		Release: r,
	}
}

func (r *ReleaseResource) IsStarted() bool {
	for _, instance := range r.Instances().List().Items {
		if !instance.IsStarted() {
			return false
		}
	}
	return true
}

func (r *ReleaseResource) IsStopped() bool {
	for _, instance := range r.Instances().List().Items {
		if !instance.IsStopped() {
			return false
		}
	}
	return true
}

func (r *ReleaseResource) imageRepoNames() (repoNames []string) { // TODO convert Image into Value object w/ repo, image, version
	uniqRepoNames := make(map[string]bool)
	for _, container := range r.Containers {
		repoName := ImageRepoName(container)
		if _, ok := uniqRepoNames[repoName]; !ok {
			uniqRepoNames[repoName] = true
			repoNames = append(repoNames, repoName)
		}
	}
	return repoNames
}

// TODO make sub-method on container, extract guts of for loop here
func (r *ReleaseResource) containerPorts(public bool) (ports []*types.Port) {

	// TODO these will need to be unique -------------------------------------------------

	for _, container := range r.Containers {
		for _, port := range container.Ports {
			if port.Public == public {
				ports = append(ports, port)
			}
		}
	}
	return ports
}

// Operations-------------------------------------------------------------------
func (r *ReleaseResource) getService(name string) (*guber.Service, error) {
	return r.collection.core.K8S.Services(r.App().Name).Get(name)
}

func (r *ReleaseResource) provisionService(name string, svcType string, svcPorts []*guber.ServicePort) error {
	if service, _ := r.getService(name); service != nil {
		return nil // already created
	}

	service := &guber.Service{
		Metadata: &guber.Metadata{
			Name: name,
		},
		Spec: &guber.ServiceSpec{
			Type: svcType,
			Selector: map[string]string{
				"service": r.Component().Name,
			},
			Ports: svcPorts,
		},
	}
	_, err := r.collection.core.K8S.Services(r.App().Name).Create(service)
	return err
}

func (r *ReleaseResource) externalServiceName() string {
	return fmt.Sprintf("%s-public", r.Component().Name)
}

func (r *ReleaseResource) internalServiceName() string {
	return r.Component().Name
}

// Exposed to fetch IPs
func (r *ReleaseResource) ExternalService() (*guber.Service, error) {
	return r.getService(r.externalServiceName())
}

func (r *ReleaseResource) InternalService() (*guber.Service, error) {
	return r.getService(r.internalServiceName())
}

func (r *ReleaseResource) provisionExternalService() error {
	var ports []*guber.ServicePort
	for _, port := range r.containerPorts(true) {
		ports = append(ports, AsKubeServicePort(port))
	}
	return r.provisionService(r.externalServiceName(), "NodePort", ports)
}

func (r *ReleaseResource) provisionInternalService() error {
	var ports []*guber.ServicePort
	for _, port := range r.containerPorts(false) {
		ports = append(ports, AsKubeServicePort(port))
	}
	return r.provisionService(r.internalServiceName(), "ClusterIP", ports)
}

func (r *ReleaseResource) deleteServices() (err error) {
	if _, err = r.collection.core.K8S.Services(r.App().Name).Delete(r.externalServiceName()); err != nil {
		return err
	}
	if _, err = r.collection.core.K8S.Services(r.App().Name).Delete(r.internalServiceName()); err != nil {
		return err
	}
	return nil
}

// NOTE it seems weird here, but "Provision" == "CreateUnlessExists"
func (r *ReleaseResource) provisionSecrets() error {
	repos, err := r.imageRepos()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		if err := r.App().ProvisionSecret(repo); err != nil {
			return err
		}
	}
	return nil
}

// Provision creates needed assets for all instances. It does not actually
// start instances -- that is handled by deploy.go.
func (r *ReleaseResource) Provision() error {
	if err := r.provisionSecrets(); err != nil {
		return err
	}

	// Create Services
	if err := r.provisionInternalService(); err != nil {
		return err
	}
	if err := r.provisionExternalService(); err != nil {
		return err
	}

	// Concurrently provision volumes
	// which is not actually concurrent... just sends all requests, and then
	// loops waiting, which prevents concurrently polling while waiting.

	for _, vol := range r.volumes() {
		if err := vol.Provision(); err != nil {
			return err
		}
	}
	for _, vol := range r.volumes() {
		if err := vol.WaitForAvailable(); err != nil {
			return err
		}
	}

	// c := make(chan error)
	// for _, instance := range r.Instances().List().Items {
	// 	go func(instance *InstanceResource) { // NOTE we have to pass instance here, else every goroutine hits the same instance
	// 		c <- instance.ProvisionVolumes()
	// 	}(instance)
	// }
	// for i := 0; i < r.InstanceCount; i++ {
	// 	if err := <-c; err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (r *ReleaseResource) volumes() (vols []*AwsVolume) {
	for _, instance := range r.Instances().List().Items {
		vols = append(vols, instance.Volumes()...)
	}
	return vols
}

func (r *ReleaseResource) imageRepos() (repos []*ImageRepoResource, err error) { // Not returning ImageRepoResource, since they are defined before hand
	for _, repoName := range r.imageRepoNames() {
		repo, err := r.collection.core.ImageRepos().Get(repoName)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// TODO naming inconsistencies for kube definitions of resources
// ImagePullSecrets returns repo names defined for Kube pods
func (r *ReleaseResource) ImagePullSecrets() (pullSecrets []*guber.ImagePullSecret, err error) {
	repos, err := r.imageRepos()
	if err != nil {
		return pullSecrets, err
	}
	for _, repo := range repos {
		pullSecrets = append(pullSecrets, AsKubeImagePullSecret(repo))
	}
	return pullSecrets, nil
}
