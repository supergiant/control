package core

import (
	"encoding/json"
	"strings"

	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type DefaultProvisioner struct {
	Core *Core
}

func (p *DefaultProvisioner) Provision(kubeResource *model.KubeResource) error {
	var resource map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Resource, &resource); err != nil {
		return err
	}

	if resource["apiVersion"] == nil {
		resource["apiVersion"] = "v1"
	}

	resource["kind"] = kubeResource.Kind

	metadata := make(map[string]interface{})
	if resource["metadata"] != nil {
		metadata = resource["metadata"].(map[string]interface{})
	}

	metadata["namespace"] = kubeResource.Namespace
	metadata["name"] = kubeResource.Name

	resource["metadata"] = metadata

	k8s := p.Core.K8S(kubeResource.Kube)

	if err := k8s.CreateResource("api/v1", kubeResource.Kind, kubeResource.Namespace, resource, kubeResource.Resource); err != nil {
		return err
	}

	// Save since we just set Resource
	return p.Core.DB.Save(kubeResource)
}

func (p *DefaultProvisioner) Teardown(kubeResource *model.KubeResource) error {
	k8s := p.Core.K8S(kubeResource.Kube)
	err := k8s.DeleteResource("api/v1", kubeResource.Kind, kubeResource.Namespace, kubeResource.Name)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	return nil
}

func (p *DefaultProvisioner) IsRunning(kubeResource *model.KubeResource) (bool, error) {
	err := p.Core.K8S(kubeResource.Kube).GetResource("api/v1", kubeResource.Kind, kubeResource.Namespace, kubeResource.Name, kubeResource.Resource)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return kubeResourceIsRunning(kubeResource)
}

//------------------------------------------------------------------------------
// TODO

func kubeResourceIsRunning(kubeResource *model.KubeResource) (bool, error) {
	if kubeResource.Kind == "Pod" {
		pod := new(kubernetes.Pod)
		if err := json.Unmarshal(*kubeResource.Resource, pod); err != nil {
			return false, err
		}

		for _, cond := range pod.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				return true, nil
			}
		}

		return false, nil

	} else if kubeResource.Kind == "PersistentVolume" {
		volume := new(kubernetes.PersistentVolume)
		if err := json.Unmarshal(*kubeResource.Resource, volume); err != nil {
			return false, err
		}

		return volume.Status.Phase == "Bound", nil
	}

	return true, nil
}
