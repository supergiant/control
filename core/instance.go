package core

import (
	"fmt"
	"time"

	"github.com/supergiant/guber"
)

// Instance is really just a Kubernetes Pod (with a better name)
type Instance struct {
	c          *Client
	Deployment *Deployment `json:"-"`
	ID         int         `json:"id"` // actually just the number (starting w/ 1) of the instance order in the deployment
	// Active bool `json:"active"`
}

func (m *Instance) Name() string {
	return fmt.Sprintf("%s-%s-%d", m.component().Name, m.Deployment.ID, m.ID)
}

func (m *Instance) component() *Component {
	return m.Deployment.Component()
}

func (m *Instance) appName() string {
	return m.component().App().Name
}

func (m *Instance) volumes() (vols []*AwsVolume) {
	for _, blueprint := range m.component().Blueprint.Volumes {
		vol := &AwsVolume{
			c:         m.c,
			Blueprint: blueprint,
			Instance:  m,
		}
		vols = append(vols, vol)
	}
	return vols
}

func (m *Instance) provisionVolumes() error {
	for _, vol := range m.volumes() {
		if err := vol.Provision(); err != nil {
			return err
		}
	}
	for _, vol := range m.volumes() {
		if err := vol.WaitForAvailable(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Instance) kubeVolumes() (vols []*guber.Volume) {
	for _, vol := range m.volumes() {
		vols = append(vols, vol.AsKubeVolume())
	}
	return vols
}

func (m *Instance) kubeContainers() (containers []*guber.Container) {
	for _, blueprint := range m.component().Blueprint.Containers {
		containers = append(containers, blueprint.AsKubeContainer(m))
	}
	return containers
}

func (m *Instance) replicationController() (*guber.ReplicationController, error) {
	return m.c.K8S.ReplicationControllers(m.appName()).Get(m.Name())
}

func (m *Instance) waitForReplicationControllerReady() error {
	start := time.Now()
	maxWait := 5 * time.Minute
	for {
		if elapsed := time.Since(start); elapsed < maxWait {
			rc, err := m.replicationController()
			if err != nil {
				return err
			} else if rc.Status.Replicas == 1 { // TODO this may not assert pod running
				break
			}
		} else {
			return fmt.Errorf("Timed out waiting for RC '%s' to start", m.Name())
		}
	}
	return nil
}

func (m *Instance) provisionReplicationController() error {
	if rc, err := m.replicationController(); err != nil {
		return err // some systemic error (err, along with rc, is nil when rc does not exist)
	} else if rc != nil {
		return nil // rc already exists
	}

	// We load them here because the repos may not exist, which needs to return error
	imagePullSecrets, err := m.component().ImagePullSecrets()
	if err != nil {
		return err
	}

	rc := &guber.ReplicationController{
		Metadata: &guber.Metadata{
			Name: m.Name(),
		},
		Spec: &guber.ReplicationControllerSpec{
			Selector: map[string]string{
				"instance": m.Name(),
			},
			Replicas: 1,
			Template: &guber.PodTemplate{
				Metadata: &guber.Metadata{
					Name: m.Name(), // pod base name is same as RC
					Labels: map[string]string{
						"deployment": m.Deployment.ID, // for Service
						"instance":   m.Name(),        // for RC
					},
				},
				Spec: &guber.PodSpec{
					Volumes:                       m.kubeVolumes(),
					Containers:                    m.kubeContainers(),
					ImagePullSecrets:              imagePullSecrets,
					TerminationGracePeriodSeconds: m.component().Blueprint.TerminationGracePeriod,
				},
			},
		},
	}
	if _, err = m.c.K8S.ReplicationControllers(m.appName()).Create(rc); err != nil {
		return err
	}
	return m.waitForReplicationControllerReady()
}

func (m *Instance) pod() (*guber.Pod, error) {
	q := &guber.QueryParams{
		LabelSelector: "instance=" + m.Name(),
	}
	pods, err := m.c.K8S.Pods(m.appName()).List(q)
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

func (m *Instance) destroyReplicationControllerAndPod() error {
	// TODO we call m.c.K8S.ReplicationControllers(m.appName()) enough to warrant its own method -- confusing nomenclature awaits assuredly
	if _, err := m.c.K8S.ReplicationControllers(m.appName()).Delete(m.Name()); err != nil {
		return err
	}
	pod, err := m.pod()
	if err != nil {
		return err
	}
	if pod != nil {
		// _ is found bool, we don't care if it was found or not, just don't want an error
		if _, err := m.c.K8S.Pods(m.appName()).Delete(pod.Metadata.Name); err != nil {
			return err
		}
	}
	return nil
}

func (m *Instance) destroyVolumes() error {
	for _, vol := range m.volumes() {
		if err := vol.Destroy(); err != nil { // NOTE this should not be a "not found" error -- since volumes() will naturally do an existence check
			return err
		}
	}
	return nil
}

// Provision is the method extended down to the user for deploy control
func (m *Instance) Provision() (err error) {
	if err = m.provisionVolumes(); err != nil {
		return err
	}
	if err = m.provisionReplicationController(); err != nil {
		return err
	}
	return nil
}

// Destroy tears down the instance
func (m *Instance) Destroy() (err error) {
	if err = m.destroyReplicationControllerAndPod(); err != nil {
		return err
	}
	if err = m.destroyVolumes(); err != nil {
		return err
	}
	return nil
}
