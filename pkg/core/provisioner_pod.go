package core

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/supergiant/supergiant/pkg/model"
)

type PodProvisioner struct {
	Core *Core
}

func (p *PodProvisioner) Provision(kubeResource *model.KubeResource) error {
	var templateMap map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Template, &templateMap); err != nil {
		return err
	}

	if templateMap["spec"] == nil {
		return errors.New("Missing spec field on Pod")
	}
	spec := templateMap["spec"].(map[string]interface{})

	// We only use this custom Provisioner to create Supergiant Volumes, so if
	// there are no volumes defined, or none of them are the special Supergiant
	// type, then we will just use the DefaultProvsiioner.
	if spec["volumes"] == nil {
		return p.Core.DefaultProvisioner.Provision(kubeResource)
	}
	volumeDefs := spec["volumes"].([]interface{})

	// Load up existing volumes. We do this so this method can be re-ran on error,
	// and volumes won't be recreated.
	volumes, err := p.existingVolumes(kubeResource)
	if err != nil {
		return err
	}

	var newVolumeDefs []interface{}

	for _, vd := range volumeDefs {

		volumeDef := vd.(map[string]interface{})

		sgVolDef, ok := volumeDef["SUPERGIANT_EXTERNAL_VOLUME"].(map[string]interface{})
		if !ok {
			newVolumeDefs = append(newVolumeDefs, volumeDef) // Append to array so we preserve non-SG volumes
			continue
		}

		name, okName := volumeDef["name"].(string)
		if !okName {
			return errors.New("Missing or malformed 'name' field in Pod Volume")
		}

		volume := volumes[name]

		if volume == nil {
			volume = &model.Volume{
				Name:           name,
				KubeName:       kubeResource.KubeName,
				KubeResourceID: kubeResource.ID,
			}

			if volType, ok := sgVolDef["type"].(string); ok {
				volume.Type = volType
			}

			if size, ok := sgVolDef["size"].(float64); ok {
				volume.Size = int(size)
			} else {
				return errors.New("Missing or malformed 'size' field in SUPERGIANT_EXTERNAL_VOLUME")
			}

			if volErr := p.Core.Volumes.Create(volume); volErr != nil {
				return volErr
			}
		}

		newVolumeDef := p.Core.CloudAccounts.provider(kubeResource.Kube.CloudAccount).KubernetesVolumeDefinition(volume)
		newVolumeDefs = append(newVolumeDefs, newVolumeDef)
	}

	// Replace old volume definitions in the original spec
	spec["volumes"] = newVolumeDefs

	// Serialize new Definition
	marshalledDef, err := json.Marshal(templateMap)
	if err != nil {
		return err
	}
	rawMsgDef := json.RawMessage(marshalledDef)
	kubeResource.Definition = &rawMsgDef

	// Run default provisioning procedure
	return p.Core.DefaultProvisioner.Provision(kubeResource)
}

func (p *PodProvisioner) Teardown(kubeResource *model.KubeResource) error {
	volumes, err := p.existingVolumes(kubeResource)
	if err != nil {
		return err
	}
	if err := p.Core.DefaultProvisioner.Teardown(kubeResource); err != nil {
		return err
	}
	for _, volume := range volumes {
		if err := p.Core.Volumes.Delete(volume.ID, volume).Now(); err != nil {
			return err
		}
	}
	return nil
}

func (p *PodProvisioner) IsRunning(kubeResource *model.KubeResource) (bool, error) {
	err := p.Core.K8S(kubeResource.Kube).GetResource(kubeResource.Kind, kubeResource.Namespace, kubeResource.Name, kubeResource.Artifact)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}

	var artifactMap map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Artifact, &artifactMap); err != nil {
		return false, err
	}

	status, _ := artifactMap["status"].(map[string]interface{})
	conditions, _ := status["conditions"].([]interface{})

	for _, condition := range conditions {
		cond := condition.(map[string]interface{})
		if cond["type"] == "Ready" && cond["status"] == "True" {
			return true, nil
		}
	}

	return false, nil
}

// Private

// We return a map here for simple find-by-name
func (p *PodProvisioner) existingVolumes(kubeResource *model.KubeResource) (map[string]*model.Volume, error) {
	var volumes []*model.Volume
	if err := p.Core.DB.Find(&volumes, "kube_resource_id = ?", kubeResource.ID); err != nil {
		return nil, err
	}
	volumeMap := make(map[string]*model.Volume)
	for _, volume := range volumes {
		volumeMap[volume.Name] = volume
	}
	return volumeMap, nil
}
