package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/supergiant/guber"
	"github.com/supergiant/supergiant/types"
)

type InstanceCollection struct {
	core    *Core
	Release *ReleaseResource
}

type InstanceResource struct {
	collection *InstanceCollection
	*types.Instance
}

type InstanceList struct {
	Items []*InstanceResource `json:"items"`
}

func (c *InstanceCollection) app() *AppResource {
	return c.Release.Component().App()
}

func (c *InstanceCollection) component() *ComponentResource {
	return c.Release.Component()
}

// List returns an InstanceList.
func (c *InstanceCollection) List() *InstanceList {
	list := new(InstanceList)
	list.Items = make([]*InstanceResource, 0)
	for i := 0; i < c.Release.InstanceCount; i++ {
		id := strconv.Itoa(i)
		list.Items = append(list.Items, c.New(&id))
	}
	return list
}

// New initializes an Instance with a pointer to the Collection.
func (c *InstanceCollection) New(id types.ID) *InstanceResource {
	r := &InstanceResource{
		collection: c,
		Instance: &types.Instance{
			ID: id,
		},
	}
	// TODO not consistent with the setter approach
	r.BaseName = *r.Component().Name + "-" + *r.ID
	r.Name = r.BaseName + *r.Release().Timestamp
	r.setStatus()
	return r
}

// Get takes an id and returns an InstanceResource if it exists.
func (c *InstanceCollection) Get(id types.ID) (*InstanceResource, error) {
	index, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}
	maxIndex := c.Release.InstanceCount - 1
	if index < 0 || index > maxIndex {
		return nil, fmt.Errorf("%d for Instance ID is out of range; Highest ID is %d", index, maxIndex)
	}
	return c.New(id), nil
}

// Resource-level
//==============================================================================

func (r *InstanceResource) App() *AppResource {
	return r.collection.app()
}

func (r *InstanceResource) Component() *ComponentResource {
	return r.collection.component()
}

func (r *InstanceResource) setStatus() {
	pod, err := r.pod()
	if err != nil {
		panic(err) // TODO
	}
	if pod != nil {
		if pod.Status.Phase == "Running" {
			r.Status = types.InstanceStatusStarted
			return
		}
	}
	r.Status = types.InstanceStatusStopped
}

func (r *InstanceResource) IsStarted() bool {
	return r.Status == types.InstanceStatusStarted
}

func (r *InstanceResource) IsStopped() bool {
	return r.Status == types.InstanceStatusStopped
}

// Delete tears down the instance
func (r *InstanceResource) Delete() (err error) {
	if err = r.deleteReplicationControllerAndPod(); err != nil {
		return err
	}
	if err = r.DeleteVolumes(); err != nil {
		return err
	}
	return nil
}

// The following 2 are only diff from Provision() and Delete() in that they do
// not delete the create or delete the volumes.
func (r *InstanceResource) Start() error {
	for _, vol := range r.Volumes() {
		if err := vol.WaitForAvailable(); err != nil {
			return err
		}
	}
	if err := r.provisionReplicationController(); err != nil {
		return err
	}
	return nil
}

func (r *InstanceResource) Stop() error {
	if err := r.deleteReplicationControllerAndPod(); err != nil {
		return err
	}
	for _, vol := range r.Volumes() {
		if err := vol.WaitForAvailable(); err != nil {
			return err
		}
	}

	// TODO need a wait in here (optional) for the pod to be deleted

	return nil
}

func (r *InstanceResource) Release() *ReleaseResource {
	return r.collection.Release
}

func (r *InstanceResource) Volumes() (vols []*AwsVolume) {
	for _, blueprint := range r.Release().Volumes {
		vol := &AwsVolume{
			core:      r.collection.core,
			Blueprint: blueprint,
			Instance:  r,
		}
		vols = append(vols, vol)
	}
	return vols
}

func (r *InstanceResource) kubeVolumes() (vols []*guber.Volume) {
	for _, vol := range r.Volumes() {
		vols = append(vols, asKubeVolume(vol))
	}
	return vols
}

func (r *InstanceResource) kubeContainers() (containers []*guber.Container) {
	for _, blueprint := range r.Release().Containers {
		containers = append(containers, asKubeContainer(blueprint, r))
	}
	return containers
}

func (r *InstanceResource) replicationController() (*guber.ReplicationController, error) {
	return r.collection.core.K8S.ReplicationControllers(*r.App().Name).Get(r.Name)
}

func (r *InstanceResource) waitForReplicationControllerReady() error {
	start := time.Now()
	maxWait := 5 * time.Minute
	for {
		if elapsed := time.Since(start); elapsed < maxWait {
			rc, err := r.replicationController()
			if err != nil {
				return err
			} else if rc.Status.Replicas == 1 { // TODO this may not assert pod running
				break
			}
		} else {
			return fmt.Errorf("Timed out waiting for RC '%s' to start", r.Name)
		}
	}
	return nil
}

func (r *InstanceResource) provisionReplicationController() error {
	if rc, err := r.replicationController(); err != nil {
		return err // some systemic error (err, along with rc, is nil when rc does not exist)
	} else if rc != nil {
		return nil // rc already exists
	}

	// We load them here because the repos may not exist, which needs to return error
	imagePullSecrets, err := r.Release().ImagePullSecrets()
	if err != nil {
		return err
	}

	rc := &guber.ReplicationController{
		Metadata: &guber.Metadata{
			Name: r.Name,
		},
		Spec: &guber.ReplicationControllerSpec{
			Selector: map[string]string{
				"instance": r.Name,
			},
			Replicas: 1,
			Template: &guber.PodTemplate{
				Metadata: &guber.Metadata{
					Name: r.Name, // pod base name is same as RC
					Labels: map[string]string{
						"service":  *r.Component().Name, // for Service
						"instance": r.Name,              // for RC (above)
					},
				},
				Spec: &guber.PodSpec{
					Volumes:                       r.kubeVolumes(),
					Containers:                    r.kubeContainers(),
					ImagePullSecrets:              imagePullSecrets,
					TerminationGracePeriodSeconds: r.Release().TerminationGracePeriod,
				},
			},
		},
	}
	if _, err = r.collection.core.K8S.ReplicationControllers(*r.App().Name).Create(rc); err != nil {
		return err
	}
	return r.waitForReplicationControllerReady()
}

func (r *InstanceResource) pod() (*guber.Pod, error) {
	q := &guber.QueryParams{
		LabelSelector: "instance=" + r.Name,
	}
	pods, err := r.collection.core.K8S.Pods(*r.App().Name).Query(q)
	if err != nil {
		return nil, err // Not sure what the error might be here
	}

	if len(pods.Items) > 1 {
		panic("More than 1 pod returned in query?")
	} else if len(pods.Items) == 1 {
		return pods.Items[0], nil
	}
	return nil, nil
}

func (r *InstanceResource) deleteReplicationControllerAndPod() error {
	// TODO we call r.collection.core.K8S.ReplicationControllers(r.App().Name) enough to warrant its own method -- confusing nomenclature awaits assuredly
	if _, err := r.collection.core.K8S.ReplicationControllers(*r.App().Name).Delete(r.Name); err != nil {
		return err
	}
	pod, err := r.pod()
	if err != nil {
		return err
	}
	if pod != nil {
		// _ is found bool, we don't care if it was found or not, just don't want an error
		if _, err := r.collection.core.K8S.Pods(*r.App().Name).Delete(pod.Metadata.Name); err != nil {
			return err
		}
	}
	return nil
}

// exposed for use in deploy_component.go
func (r *InstanceResource) DeleteVolumes() error {
	for _, vol := range r.Volumes() {
		if err := vol.Delete(); err != nil { // NOTE this should not be a "not found" error -- since Volumes() will naturally do an existence check
			return err
		}
	}
	return nil
}
