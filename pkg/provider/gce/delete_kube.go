package gce

import (
	"strings"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// DeleteKube deletes a GCE kubernetes cluster.
func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
	// setup provider steps.
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Delete Kube",
		Model:  m,
		Action: action,
	}
	// fetch client.
	client, err := p.Client(m)
	if err != nil {
		return err
	}
	// Delete all master nodes.
	procedure.AddStep("Destroying Kubernetes Master(s)...", func() error {
		for _, master := range m.GCEConfig.MasterNodes {
			_, err := client.Instances.Delete(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, convInstanceURLtoString(master)).Do()
			if err != nil {
				if strings.Contains(err.Error(), "notFound") {
					// it does not exist,
					return nil
				}
				return err
			}
		}
		return nil
	})

	procedure.AddStep("Destroying Kubernetes Minions...", func() error {
		for _, node := range m.Nodes {
			_, err := client.Instances.Delete(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, convInstanceURLtoString(node.Name)).Do()
			if err != nil && !strings.Contains(err.Error(), "was not found") && !strings.Contains(err.Error(), "Values must match") {
				return err
			}
		}
		return nil
	})

	procedure.AddStep("Destroying Kubernetes master instance group...", func() error {
		_, err := client.InstanceGroups.Delete(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, m.Name+"-kubernetes-masters").Do()
		if err != nil && !strings.Contains(err.Error(), "was not found") && !strings.Contains(err.Error(), "Unknown zone") {
			return err
		}
		return nil
	})

	procedure.AddStep("Destroying Kubernetes minion instance group...", func() error {
		_, err := client.InstanceGroups.Delete(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, m.Name+"-kubernetes-minions").Do()
		if err != nil && !strings.Contains(err.Error(), "was not found") && !strings.Contains(err.Error(), "Unknown zone") {
			return err
		}
		return nil
	})
	// Initialize steps
	return procedure.Run()
}
