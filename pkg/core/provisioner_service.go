package core

import (
	"encoding/json"
	"errors"

	"github.com/supergiant/supergiant/pkg/model"
)

type ServiceProvisioner struct {
	Core *Core
}

type serviceProvisionerAssetActionType int

const (
	serviceProvisionerAssetDelete serviceProvisionerAssetActionType = iota // Delete is 0
	serviceProvisionerAssetCreate
	serviceProvisionerAssetKeep
	serviceProvisionerAssetReplace
)

type serviceProvisionerAsset struct {
	model         *model.EntrypointListener
	plannedAction serviceProvisionerAssetActionType
}

// We return a map here for simple find-by-name
func (p *ServiceProvisioner) existingAssets(r *model.KubeResource) (map[string]*serviceProvisionerAsset, error) {
	var listeners []*model.EntrypointListener
	if err := p.Core.DB.Find(&listeners, "kube_resource_id = ?", r.ID); err != nil {
		return nil, err
	}
	assets := make(map[string]*serviceProvisionerAsset)
	for _, listener := range listeners {
		assets[listener.Name] = &serviceProvisionerAsset{model: listener}
	}
	return assets, nil
}

//------------------------------------------------------------------------------

func (p *ServiceProvisioner) Provision(kubeResource *model.KubeResource) error {

	var templateMap map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Template, &templateMap); err != nil {
		return err
	}

	// Get the ports array from the copied Template
	spec, _ := templateMap["spec"].(map[string]interface{})
	ports, _ := spec["ports"].([]interface{})

	// Initialize asset action map by loading existing assets if they exist
	assets, err := p.existingAssets(kubeResource)
	if err != nil {
		return err
	}

	var newPortDefs []map[string]interface{}

	for _, p := range ports {

		port := p.(map[string]interface{})

		// If there is no entrypoint def on this port, then we leave it alone
		entrypointDef, ok := port["SUPERGIANT_ENTRYPOINT_LISTENER"].(map[string]interface{})
		if !ok {
			newPortDefs = append(newPortDefs, port) // Append so we preserve non-SG ports
			continue
		}

		// We've got to have a named port, because we need an identifier (since the
		// ports can change, they can't be used as an identifier).
		portName, _ := port["name"].(string)
		if portName == "" {
			return errors.New("Port must have a 'name' field with SUPERGIANT_ENTRYPOINT_LISTENER")
		}

		// Build up the EntrypointListener (asset)
		newEntrypointListener := &model.EntrypointListener{
			KubeResourceID: kubeResource.ID,
			Name:           portName,
		}

		entrypointPort := entrypointDef["entrypoint_port"].(float64)

		newEntrypointListener.NodeProtocol, _ = port["protocol"].(string)
		newEntrypointListener.EntrypointName, _ = entrypointDef["entrypoint_name"].(string)
		newEntrypointListener.EntrypointPort = int64(entrypointPort)
		newEntrypointListener.EntrypointProtocol, _ = entrypointDef["entrypoint_protocol"].(string)

		// If there isn't already an EntrypointListener, prepare to create one
		if existingAsset := assets[portName]; existingAsset == nil {
			assets[portName] = &serviceProvisionerAsset{
				model:         newEntrypointListener,
				plannedAction: serviceProvisionerAssetCreate,
			}

		} else {
			// If we have an existing Listener, then we received a nodePort assignment
			// from Kubernetes. We preserve that here.
			port["nodePort"] = int(existingAsset.model.NodePort)

			definitionChanged :=
				newEntrypointListener.NodeProtocol != existingAsset.model.NodeProtocol ||
					newEntrypointListener.EntrypointName != existingAsset.model.EntrypointName ||
					newEntrypointListener.EntrypointPort != existingAsset.model.EntrypointPort ||
					newEntrypointListener.EntrypointProtocol != existingAsset.model.EntrypointProtocol

			if definitionChanged {
				// If the definition changed, we have to replace.
				existingAsset.model.NodeProtocol = newEntrypointListener.NodeProtocol
				existingAsset.model.EntrypointName = newEntrypointListener.EntrypointName
				existingAsset.model.EntrypointPort = newEntrypointListener.EntrypointPort
				existingAsset.model.EntrypointProtocol = newEntrypointListener.EntrypointProtocol
				existingAsset.plannedAction = serviceProvisionerAssetReplace

			} else {
				// Otherwise, we're doing nothing.
				existingAsset.plannedAction = serviceProvisionerAssetKeep
			}
		}

		// Remove the Supergiant-specific definition
		delete(port, "SUPERGIANT_ENTRYPOINT_LISTENER")

		newPortDefs = append(newPortDefs, port)
	}

	// Replace old port definitions in the original spec
	if len(newPortDefs) > 0 {
		spec["ports"] = newPortDefs
	}

	// Delete all the EntrypointListeners we no longer need
	for _, asset := range assets {
		if asset.plannedAction == serviceProvisionerAssetDelete {
			if err := p.Core.EntrypointListeners.Delete(asset.model.ID, asset.model).Now(); err != nil {
				return err
			}
		}
	}

	// Serialize new Definition
	marshalledDef, err := json.Marshal(templateMap)
	if err != nil {
		return err
	}
	rawMsgDef := json.RawMessage(marshalledDef)
	kubeResource.Definition = &rawMsgDef

	// Create the Service (this is where new ports get nodePort assignments)
	if err := p.Core.DefaultProvisioner.Provision(kubeResource); err != nil {
		return err
	}

	var artifactMap map[string]interface{}
	if err := json.Unmarshal(*kubeResource.Artifact, &artifactMap); err != nil {
		return err
	}

	// Get ports array from the Artifact (which was saved by DefaultProvisioner)
	artifactSpec, _ := artifactMap["spec"].(map[string]interface{})
	artifactPorts, _ := artifactSpec["ports"].([]interface{})

	// Iterate through the ports in the Kubernetes response, and run relevant
	// actions we have planned.
	for _, ap := range artifactPorts {

		artifactPort := ap.(map[string]interface{})

		portName, _ := artifactPort["name"].(string)
		asset := assets[portName]

		// If there is no asset, this is assumed to be a port without
		// SUPERGIANT_ENTRYPOINT_LISTENER.
		if asset == nil {
			continue
		}

		// Capture the nodePort assignment.
		nodePort := artifactPort["nodePort"].(float64)
		asset.model.NodePort = int64(nodePort)

		// Run the planned action.
		switch asset.plannedAction {

		case serviceProvisionerAssetKeep:
			continue // we do nothing

		case serviceProvisionerAssetCreate:
			if err := p.Core.EntrypointListeners.Create(asset.model); err != nil {
				return err
			}

		case serviceProvisionerAssetReplace:
			// We don't want to overwrite asset.model, so we pass a dummy to render to
			// (we pass it the name for testing convenience)
			dummy := &model.EntrypointListener{Name: asset.model.Name}
			if err := p.Core.EntrypointListeners.Delete(asset.model.ID, dummy).Now(); err != nil {
				return err
			}
			if err := p.Core.EntrypointListeners.Create(asset.model); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ServiceProvisioner) Teardown(kubeResource *model.KubeResource) error {
	assets, err := p.existingAssets(kubeResource)
	if err != nil {
		return err
	}
	for _, asset := range assets {
		if err := p.Core.EntrypointListeners.Delete(asset.model.ID, asset.model).Now(); err != nil {
			return err
		}
	}
	return p.Core.DefaultProvisioner.Teardown(kubeResource)
}

func (p *ServiceProvisioner) IsRunning(kubeResource *model.KubeResource) (bool, error) {
	return p.Core.DefaultProvisioner.IsRunning(kubeResource)
}
