package core

import (
	"encoding/json"
	"strings"

	"github.com/supergiant/supergiant/pkg/model"
)

type DefaultProvisioner struct {
	Core *Core
}

func (p *DefaultProvisioner) Provision(kubeResource *model.KubeResource) error {
	// If this is called directly, as opposed to by one of the non-default
	// Provisioners, we will need to make sure Template is copied to Definition.
	if kubeResource.Definition == nil || len(*kubeResource.Definition) == 0 {
		defRawMsg := make(json.RawMessage, len(*kubeResource.Template))
		kubeResource.Definition = &defRawMsg
		copy(*kubeResource.Definition, *kubeResource.Template)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Definition, &resource); err != nil {
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

	artifact := make(json.RawMessage, 0)
	kubeResource.Artifact = &artifact
	if err := k8s.CreateResource(kubeResource.Kind, kubeResource.Namespace, resource, kubeResource.Artifact); err != nil {
		return err
	}

	// Save since we just set Artifact
	return p.Core.DB.Save(kubeResource)
}

func (p *DefaultProvisioner) Teardown(kubeResource *model.KubeResource) error {
	k8s := p.Core.K8S(kubeResource.Kube)
	err := k8s.DeleteResource(kubeResource.Kind, kubeResource.Namespace, kubeResource.Name)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	return nil
}

func (p *DefaultProvisioner) IsRunning(kubeResource *model.KubeResource) (bool, error) {
	err := p.Core.K8S(kubeResource.Kube).GetResource(kubeResource.Kind, kubeResource.Namespace, kubeResource.Name, kubeResource.Artifact)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
