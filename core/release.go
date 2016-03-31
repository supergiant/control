package core

import (
	"fmt"
	"path"
	"strconv"
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

	// TODO these are shared between releases, it's kinda funky right now
	externalService *guber.Service
	internalService *guber.Service
	imageRepos      []*ImageRepoResource
	entrypoints     map[string]*EntrypointResource // a map for quick lookup
}

type ReleaseList struct {
	Items []*ReleaseResource `json:"items"`
}

// EtcdKey implements the Collection interface.
func (c *ReleaseCollection) EtcdKey(timestamp types.ID) string {
	key := path.Join("/releases", *c.Component.App().Name, *c.Component.Name)
	if timestamp != nil {
		key = path.Join(key, *timestamp)
	}
	return key
}

// InitializeResource implements the Collection interface.
func (c *ReleaseCollection) InitializeResource(r Resource) {
	resource := r.(*ReleaseResource)
	resource.collection = c

	// TODO
	// We do this here because this is called when pulling from the DB. If it's
	// being pulled from the DB, it can be assumed to have services.
	// Still really sloppy, since there could be an error.
	svc, err := resource.getService(resource.externalServiceName())
	if err != nil {
		panic(err)
	}
	resource.externalService = svc

	svc, err = resource.getService(resource.internalServiceName())
	if err != nil {
		panic(err)
	}
	resource.internalService = svc

	repos, err := resource.getImageRepos()
	if err != nil {
		panic(err)
	}
	resource.imageRepos = repos

	resource.entrypoints = resource.getEntrypoints()
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
		Release: &types.Release{
			Meta: types.NewMeta(),
		},
	}
}

// Create takes an Release and creates it in etcd.
func (c *ReleaseCollection) Create(r *ReleaseResource) (*ReleaseResource, error) {
	r.Timestamp = newReleaseTimestamp()
	if err := c.core.DB.Create(c, r.Timestamp, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Get takes an id and returns an ReleaseResource if it exists.
func (c *ReleaseCollection) Get(id types.ID) (*ReleaseResource, error) {
	r := c.New()
	if err := c.core.DB.Get(c, id, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================

// PersistableObject satisfies the Resource interface
func (r *ReleaseResource) PersistableObject() interface{} {
	return r.Release
}

// Delete removes all assets (volumes, pods, etc.) and deletes the Release in
// etcd.
func (r *ReleaseResource) Delete() error {
	if err := r.deleteServices(); err != nil {
		return err
	}
	if err := r.removeExternalPortsFromEntrypoint(); err != nil {
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

	return r.collection.core.DB.Delete(r.collection, r.Timestamp)
}

func newReleaseTimestamp() types.ID {
	stamp := time.Now().Format("20060102150405")
	return &stamp
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
	for _, container := range r.Containers {
		repoNames = append(repoNames, ImageRepoName(container))
	}
	return uniqStrs(repoNames)
}

func (r *ReleaseResource) getEntrypoints() map[string]*EntrypointResource { // TODO convert Image into Value object w/ repo, image, version
	entrypoints := make(map[string]*EntrypointResource)
	for _, port := range r.containerPorts(true) {
		if port.EntrypointDomain == nil {
			continue
		}
		entrypoint, err := r.collection.core.Entrypoints().Get(port.EntrypointDomain)
		if err != nil {
			panic(err) // TODO
		}
		entrypoints[*port.EntrypointDomain] = entrypoint
	}
	return entrypoints
}

func (r *ReleaseResource) containerPorts(public bool) (ports []*types.Port) {
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
	return r.collection.core.K8S.Services(*r.App().Name).Get(name)
}

func (r *ReleaseResource) provisionService(name string, svcType string, svcPorts []*guber.ServicePort) (*guber.Service, error) {
	if service, _ := r.getService(name); service != nil {
		return service, nil // already created
	}

	service := &guber.Service{
		Metadata: &guber.Metadata{
			Name: name,
		},
		Spec: &guber.ServiceSpec{
			Type: svcType,
			Selector: map[string]string{
				"service": *r.Component().Name,
			},
			Ports: svcPorts,
		},
	}
	return r.collection.core.K8S.Services(*r.App().Name).Create(service)
}

func (r *ReleaseResource) externalServiceName() string {
	return *r.Component().Name + "-public"
}

func (r *ReleaseResource) internalServiceName() string {
	return *r.Component().Name
}

func (r *ReleaseResource) provisionExternalService() error {
	var ports []*guber.ServicePort
	for _, port := range r.containerPorts(true) {
		ports = append(ports, asKubeServicePort(port))
	}
	svc, err := r.provisionService(r.externalServiceName(), "NodePort", ports)
	if err != nil {
		return err
	}
	// TODO repeated in initialization
	r.externalService = svc
	return nil
}

func (r *ReleaseResource) provisionInternalService() error {
	var ports []*guber.ServicePort
	for _, port := range r.containerPorts(false) {
		ports = append(ports, asKubeServicePort(port))
	}
	svc, err := r.provisionService(r.internalServiceName(), "ClusterIP", ports)
	if err != nil {
		return err
	}
	// TODO repeated in initialization
	r.internalService = svc
	return nil
}

func (r *ReleaseResource) deleteServices() (err error) {
	if _, err = r.collection.core.K8S.Services(*r.App().Name).Delete(r.externalServiceName()); err != nil {
		return err
	}
	if _, err = r.collection.core.K8S.Services(*r.App().Name).Delete(r.internalServiceName()); err != nil {
		return err
	}
	return nil
}

// NOTE it seems weird here, but "Provision" == "CreateUnlessExists"
func (r *ReleaseResource) provisionSecrets() error {
	for _, repo := range r.imageRepos {
		if err := r.App().ProvisionSecret(repo); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReleaseResource) nodePortFor(number int) (p int) {
	for _, port := range r.externalService.Spec.Ports {
		if port.Port == number {
			p = port.NodePort
		}
	}
	return p
}

func (r *ReleaseResource) elbPortFor(port *types.Port) (elbPort int) {
	if port.PreserveNumber {
		elbPort = port.Number
	} else {
		elbPort = r.nodePortFor(port.Number)
	}
	return elbPort
}

func (r *ReleaseResource) addExternalPortsToEntrypoint() error {
	// TODO repeated, move to Port class
	for _, port := range r.containerPorts(true) {
		if port.EntrypointDomain != nil {
			nodePort := r.nodePortFor(port.Number)
			elbPort := r.elbPortFor(port)
			entrypoint := r.entrypoints[*port.EntrypointDomain]
			if err := entrypoint.AddPort(elbPort, nodePort); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReleaseResource) removeExternalPortsFromEntrypoint() error {
	// TODO repeated, move to Port class
	for _, port := range r.containerPorts(true) {
		if port.EntrypointDomain != nil {
			// nodePort := r.nodePortFor(port.Number)
			elbPort := r.elbPortFor(port)
			entrypoint := r.entrypoints[*port.EntrypointDomain]
			if err := entrypoint.RemovePort(elbPort); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReleaseResource) ExternalAddresses() (pAddrs []*types.PortAddress) {
	for _, port := range r.containerPorts(true) {
		entrypoint := r.entrypoints[*port.EntrypointDomain]
		pAddr := &types.PortAddress{
			Port:    strconv.Itoa(port.Number), // TODO repeated all over the place
			Address: fmt.Sprintf("%s:%d", entrypoint.Address, r.elbPortFor(port)),
		}
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (r *ReleaseResource) InternalAddresses() (pAddrs []*types.PortAddress) {
	for _, port := range r.containerPorts(false) {

		// TODO this should be a method in Guber
		svcDNS := fmt.Sprintf("%s.%s.svc.cluster.local", r.internalServiceName(), *r.App().Name)

		pAddr := &types.PortAddress{
			Port:    strconv.Itoa(port.Number), // TODO repeated all over the place
			Address: fmt.Sprintf("%s:%d", svcDNS, port.Number),
		}
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
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

	if err := r.addExternalPortsToEntrypoint(); err != nil {
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

func (r *ReleaseResource) getImageRepos() (repos []*ImageRepoResource, err error) { // Not returning ImageRepoResource, since they are defined before hand
	for _, repoName := range r.imageRepoNames() {
		repo, err := r.collection.core.ImageRepos().Get(&repoName)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// TODO naming inconsistencies for kube definitions of resources
// ImagePullSecrets returns repo names defined for Kube pods
func (r *ReleaseResource) ImagePullSecrets() (pullSecrets []*guber.ImagePullSecret, err error) { // TODO don't need to return error here it seems
	for _, repo := range r.imageRepos {
		pullSecrets = append(pullSecrets, asKubeImagePullSecret(repo))
	}
	return pullSecrets, nil
}
