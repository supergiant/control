package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"guber"
	"strconv"
	"strings"
	cloud "supergiant/core/aws" // TODO (name conflict with aws sdk)
	"supergiant/core/model"
	"supergiant/core/storage"
	"time"
)

type DeployComponentMessage struct {
	AppName       string
	ComponentName string
}

// DeployComponent implements job.Performable interface
type DeployComponent struct {
	db   *storage.Client
	kube *guber.Client
}

func (j DeployComponent) MaxAttempts() int {
	return 20
}

func (j DeployComponent) Perform(data []byte) error {
	message := new(DeployComponentMessage)
	if err := json.Unmarshal(data, message); err != nil {
		return err
	}

	component, err := j.db.ComponentStorage.Get(message.AppName, message.ComponentName)
	if err != nil {
		return err
	}

	// create namespace
	namespace, err := j.kube.Namespaces().Get(message.AppName)
	if err != nil {
		namespace = &guber.Namespace{
			Metadata: &guber.Metadata{
				Name: message.AppName,
			},
		}
		namespace, err = j.kube.Namespaces().Create(namespace)
		if err != nil {
			return err
		}
	}

	// get repo names from images
	var repoNames []string
	for _, container := range component.Containers {
		repoName := strings.Split(container.Image, "/")[0]
		repoNames = append(repoNames, repoName)
	}

	// create secrets for namespace
	// and create imagePullSecrets from repos
	var imagePullSecrets []*guber.ImagePullSecret
	for _, repoName := range repoNames {
		repo, err := j.db.ImageRepoStorage.Get(repoName)
		if err != nil {
			return err
		}

		secret, err := j.kube.Secrets(namespace.Metadata.Name).Get(repo.Name)
		if err != nil {
			secret = &guber.Secret{
				Metadata: &guber.Metadata{
					Name: repo.Name,
				},
				Type: "kubernetes.io/dockercfg",
				Data: map[string]string{
					".dockercfg": repo.Key,
				},
			}
			secret, err = j.kube.Secrets(namespace.Metadata.Name).Create(secret)
			if err != nil {
				return err
			}
		}

		imagePullSecret := &guber.ImagePullSecret{repo.Name}
		imagePullSecrets = append(imagePullSecrets, imagePullSecret)
	}

	// get all (uniq) ports from containers array
	// divide into private(ClusterIP) and public(NodePort) groups
	var externalPorts []*model.Port
	var internalPorts []*model.Port
	for _, container := range component.Containers {
		for _, port := range container.Ports {
			if port.Public == true {
				externalPorts = append(externalPorts, port)
			} else {
				internalPorts = append(internalPorts, port)
			}
		}
	}

	// create internal service
	if len(internalPorts) > 0 {
		internalServiceName := component.Name
		service, err := j.kube.Services(namespace.Metadata.Name).Get(internalServiceName)
		if err != nil {

			var servicePorts []*guber.ServicePort
			for _, port := range internalPorts {
				servicePort := &guber.ServicePort{
					Name: strconv.Itoa(port.Number),
					// Protocol: "TCP", // TODO (UDP supported as well)
					Port: port.Number,
				}
				servicePorts = append(servicePorts, servicePort)
			}

			service = &guber.Service{
				Metadata: &guber.Metadata{
					Name: internalServiceName,
				},
				Spec: &guber.ServiceSpec{
					Selector: map[string]string{
						"deployment": component.ActiveDeploymentID,
					},
					Ports: servicePorts,
				},
			}
			service, err = j.kube.Services(namespace.Metadata.Name).Create(service)
			if err != nil {
				return err
			}
		}
	}

	// create external service
	if len(externalPorts) > 0 {
		externalServiceName := fmt.Sprintf("%s-public", component.Name)
		service, err := j.kube.Services(namespace.Metadata.Name).Get(externalServiceName)
		if err != nil {

			var servicePorts []*guber.ServicePort
			for _, port := range externalPorts {
				servicePort := &guber.ServicePort{
					Name: strconv.Itoa(port.Number),
					// Protocol: "TCP",
					Port: port.Number,
					// NodePort:
				}
				servicePorts = append(servicePorts, servicePort)
			}

			service = &guber.Service{
				Metadata: &guber.Metadata{
					Name: externalServiceName,
				},
				Spec: &guber.ServiceSpec{
					Type: "NodePort",
					Selector: map[string]string{
						"deployment": component.ActiveDeploymentID,
					},
					Ports: servicePorts,
				},
			}
			service, err = j.kube.Services(namespace.Metadata.Name).Create(service)
			if err != nil {
				return err
			}
		}
	}

	volMgr := &cloud.VolumeManager{}

	// for each instance,
	//    create RC, and interpolate data into Env vars
	var rcs []*guber.ReplicationController
	for instanceID := 1; instanceID < component.Instances+1; instanceID++ {
		rcName := fmt.Sprintf("%s-%d", component.Name, instanceID)
		rc, err := j.kube.ReplicationControllers(namespace.Metadata.Name).Get(rcName)

		// if err == nil {
		// 	continue
		// }
		if err == nil {
			rcs = append(rcs, rc)
			continue
		}

		var volumes []*cloud.Volume
		outc := make(chan *cloud.Volume)
		errc := make(chan error)
		timeout := time.After(10 * time.Minute)
		for _, volDef := range component.Volumes {
			volName := fmt.Sprintf("%s-%d", volDef.Name, instanceID)
			go func() {
				volume, err := volMgr.Find(volName)
				if volume == nil {
					volume, err = volMgr.Create(volName, volDef.Type, volDef.Size)

					// TODO
					if volume != nil {
						volMgr.WaitForAvailable(*volume.VolumeId)
					}

				}
				if err != nil {
					errc <- err
				} else {
					outc <- &cloud.Volume{volume, volDef.Name, volName}
				}
			}()
		}
		for i := 0; i < len(component.Volumes); i++ {
			select {
			case volume := <-outc:
				volumes = append(volumes, volume)
			case err := <-errc:
				return err
			case <-timeout:
				return errors.New("Timed out waiting for volume creation")
			}
		}

		var podVols []*guber.Volume
		for _, volume := range volumes {
			podVol := &guber.Volume{
				Name: volume.BaseName, // use base name since it is localized to spec
				AwsElasticBlockStore: &guber.AwsElasticBlockStore{
					VolumeID: *volume.VolumeId,
					FSType:   "ext4",
				},
			}
			podVols = append(podVols, podVol)
		}

		var containers []*guber.Container
		for i, containerDef := range component.Containers {

			var volumeMounts []*guber.VolumeMount
			for _, mountDef := range containerDef.Mounts {
				volumeMount := &guber.VolumeMount{
					Name:      mountDef.Volume,
					MountPath: mountDef.Path,
				}
				volumeMounts = append(volumeMounts, volumeMount)
			}

			var containerPorts []*guber.ContainerPort
			for _, portDef := range containerDef.Ports {
				containerPort := &guber.ContainerPort{
					ContainerPort: portDef.Number,
				}
				containerPorts = append(containerPorts, containerPort)
			}

			var env []*guber.EnvVar
			for _, varDef := range containerDef.Env {
				envVar := &guber.EnvVar{varDef.Name, varDef.Value}
				env = append(env, envVar)
			}

			container := &guber.Container{

				// Name:  containerDef.Image,
				// TODO !
				Name: strconv.Itoa(i),

				Image: containerDef.Image,
				Env:   env,
				Resources: &guber.Resources{
					Requests: &guber.ResourceValues{
						Memory: fmt.Sprintf("%dMi", containerDef.RAM.Min),
						CPU:    fmt.Sprintf("%dm", containerDef.CPU.Min),
					},
					Limits: &guber.ResourceValues{
						Memory: fmt.Sprintf("%dMi", containerDef.RAM.Max),
						CPU:    fmt.Sprintf("%dm", containerDef.CPU.Max),
					},
				},
				VolumeMounts: volumeMounts,
				Ports:        containerPorts,
			}
			containers = append(containers, container)
		}

		rc = &guber.ReplicationController{
			Metadata: &guber.Metadata{
				Name: rcName,
			},
			Spec: &guber.ReplicationControllerSpec{
				Selector: map[string]string{
					"instance": rcName,
				},
				Replicas: 1,
				Template: &guber.PodTemplate{
					Metadata: &guber.Metadata{
						Labels: map[string]string{
							"deployment": component.ActiveDeploymentID, // for service
							"instance":   rcName,                       // for RC
						},
					},
					Spec: &guber.PodSpec{
						Volumes:                       podVols,
						Containers:                    containers,
						ImagePullSecrets:              imagePullSecrets,
						TerminationGracePeriodSeconds: component.TerminationGracePeriod,
					},
				},
			},
		}
		rc, err = j.kube.ReplicationControllers(namespace.Metadata.Name).Create(rc)
		if err != nil {
			return err
		}

		rcs = append(rcs, rc)
	}

	// wait for RCs to be ready
	start := time.Now()
	maxWait := 5 * time.Minute
	for {
		elapsed := time.Since(start)
		if elapsed < maxWait {

			ready := true

			for _, rc := range rcs {

				str, _ := json.Marshal(rc)
				fmt.Println(string(str))

				rc, err = j.kube.ReplicationControllers(namespace.Metadata.Name).Get(rc.Metadata.Name)
				if err != nil {
					return err
				}

				str, _ = json.Marshal(rc)
				fmt.Println(string(str))

				if rc.Status.Replicas < 1 { // TODO this may not assert pod running
					ready = false
				}
			}

			if ready == true {
				break
			}

		} else {
			return errors.New("Timed out waiting for RCs to start")
		}
	}

	return nil
}
