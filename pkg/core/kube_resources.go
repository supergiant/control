package core

import (
	"encoding/json"
	"fmt"
	"time"

	"context"

	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

type KubeResourcesInterface interface {
	Populate() error
	Create(*model.KubeResource) error
	Get(*int64, model.Model) error
	GetWithIncludes(*int64, model.Model, []string) error
	Update(*int64, *model.KubeResource, *model.KubeResource) error
	Delete(*int64, *model.KubeResource) ActionInterface
	Start(*int64, *model.KubeResource) ActionInterface
	Stop(*int64, *model.KubeResource) ActionInterface
	Refresh(*model.KubeResource) error
}

type KubeResources struct {
	Collection
}

func (c *KubeResources) Populate() error {
	var kubes []*model.Kube
	if err := c.Core.DB.Preload("KubeResources").Where("ready = ?", true).Find(&kubes); err != nil {
		return err
	}

	for _, kube := range kubes {
		newResources, err := getKubeResources(c.Core, kube)
		if err != nil {
			return err
		}

		oldResources := kube.KubeResources

		for _, newResource := range newResources {
			var oldResource *model.KubeResource
			oldIndex := 0

			for i, resource := range oldResources {
				if resource.Namespace == newResource.Namespace && resource.Kind == newResource.Kind && resource.Name == newResource.Name {
					oldResource = resource
					oldIndex = i
					break
				}
			}

			if oldResource != nil {
				// remove from oldResources
				oldResources = append(oldResources[:oldIndex], oldResources[oldIndex+1:]...)
				// Update

				newResource.ID = oldResource.ID
				newResource.UUID = oldResource.UUID

				// TODO this is a bit.. freestyle right now.
				// We set i there and let the Save in Refresh() persist it.
				// Don't care about errors from Heapster.
				if newResource.Kind == "Pod" {
					cpuMetrics, _ := c.Core.K8S(kube).ListPodHeapsterCPUUsageMetrics(newResource.Namespace, newResource.Name)
					ramMetrics, _ := c.Core.K8S(kube).ListPodHeapsterRAMUsageMetrics(newResource.Namespace, newResource.Name)
					if cpuMetrics != nil && ramMetrics != nil {
						newResource.ExtraData = map[string]interface{}{
							"metrics": map[string][]*kubernetes.HeapsterMetric{
								"cpu_usage": cpuMetrics,
								"ram_usage": ramMetrics,
							},
						}
					}
				}
				if err := c.Core.KubeResources.Refresh(newResource); err != nil {
					return err
				}
			} else {
				newResource.Started, err = kubeResourceIsRunning(newResource)
				if err != nil {
					return err
				}
				if err := c.Collection.Create(newResource); err != nil {
					return err
				}
			}
		}

		for _, oldResource := range oldResources {
			if err := c.Core.DB.Delete(oldResource); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *KubeResources) Create(m *model.KubeResource) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	// NOTE we call this from core to get the interface
	return c.Core.KubeResources.Start(m.ID, m).Async()
}

// TODO
func (c *KubeResources) Update(id *int64, oldM *model.KubeResource, m *model.KubeResource) error {
	return c.Collection.Update(id, oldM, m)
}

func (c *KubeResources) Delete(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("Kube.CloudAccount"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(_ *Action) error {
			if err := c.provisioner(m).Teardown(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}

func (c *KubeResources) Start(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "starting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			k8s := c.Core.K8S(m.Kube)
			if err := k8s.EnsureNamespace(m.Namespace); err != nil {
				return err
			}
			if err := c.provisioner(m).Provision(m); err != nil {
				return err
			}
			// Wait for Resource to be ready
			desc := fmt.Sprintf("%s '%s' in Namespace '%s' to start", m.Kind, m.Name, m.Namespace)

			// TODO(stgleb): Context should be inherited from higher level context
			ctx, _ := context.WithTimeout(context.Background(), c.Core.KubeResourceStartTimeout)

			waitErr := util.WaitFor(ctx, desc, 3*time.Second, func() (bool, error) {
				return c.provisioner(m).IsRunning(m)
			})
			if waitErr != nil {
				return waitErr
			}
			return c.Core.DB.Model(m).Update("started", true)
		},
	}
}

func (c *KubeResources) Stop(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "stopping",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("Kube.CloudAccount"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(a *Action) error {
			if err := c.provisioner(m).Teardown(m); err != nil {
				return err
			}
			return c.Core.DB.Model(m).Update("started", false)
		},
	}
}

// TODO TODO TODO
func (c *KubeResources) Refresh(m *model.KubeResource) (err error) {
	m.Started, err = kubeResourceIsRunning(m)
	if err != nil {
		return err
	}
	return c.Core.DB.Save(m)
}

// Private

func (c *KubeResources) provisioner(m *model.KubeResource) Provisioner {
	return c.Core.DefaultProvisioner
}

//------------------------------------------------------------------------------

func getKubeResources(c *Core, kube *model.Kube) (kr []*model.KubeResource, err error) {

	// namespaces, err := c.K8S(kube).ListNamespaces("")
	// if err != nil {
	// 	return nil, err
	// }
	//
	// for _, namespace := range namespaces {

	// Pods
	pods, err := c.K8S(kube).ListPods("")
	if err != nil {
		return nil, err
	}
	for _, pod := range pods {

		// TODO see below
		podJSON, _ := json.Marshal(pod)
		template := json.RawMessage(podJSON)

		kubeResource := &model.KubeResource{
			KubeName:  kube.Name,
			Namespace: pod.Metadata.Namespace,
			Kind:      "Pod",
			Name:      pod.Metadata.Name,
			Resource:  &template,
		}
		kr = append(kr, kubeResource)
	}

	// Services
	services, err := c.K8S(kube).ListServices("")
	if err != nil {
		return nil, err
	}
	for _, service := range services {

		// TODO see below
		serviceJSON, _ := json.Marshal(service)
		template := json.RawMessage(serviceJSON)

		kubeResource := &model.KubeResource{
			KubeName:  kube.Name,
			Namespace: service.Metadata.Namespace,
			Kind:      "Service",
			Name:      service.Metadata.Name,
			Resource:  &template,
		}
		kr = append(kr, kubeResource)
	}

	// Volumes
	volumes, err := c.K8S(kube).ListPersistentVolumes("")
	if err != nil {
		return nil, err
	}
	for _, volume := range volumes {

		// TODO see below
		volumeJSON, _ := json.Marshal(volume)
		template := json.RawMessage(volumeJSON)

		kubeResource := &model.KubeResource{
			KubeName: kube.Name,
			Kind:     "PersistentVolume",
			Name:     volume.Metadata.Name,
			Resource: &template,
		}
		kr = append(kr, kubeResource)
	}

	return
}
