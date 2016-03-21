package model

import (
	"fmt"
	"guber"
	"time"
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

func (m *Instance) loadReplicationController() (*guber.ReplicationController, error) {
	return m.c.K8S.ReplicationControllers(m.appName()).Get(m.Name())
}

func (m *Instance) waitForReplicationControllerReady() error {
	start := time.Now()
	maxWait := 5 * time.Minute
	for {
		if elapsed := time.Since(start); elapsed < maxWait {
			rc, err := m.loadReplicationController()
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
	if rc, _ := m.loadReplicationController(); rc != nil {
		return nil
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
