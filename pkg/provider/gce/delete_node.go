package gce

import (
	"strings"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// DeleteNode deletes a minsion on a DO kubernetes cluster.
func (p *Provider) DeleteNode(m *model.Node, action *core.Action) error {

	// setup provider steps.
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Delete Node",
		Model:  m,
		Action: action,
	}

	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	procedure.AddStep("Destroying Kubernetes node...", func() error {
		_, err := client.Instances.Delete(m.Kube.CloudAccount.Credentials["project_id"], m.Kube.GCEConfig.Zone, convInstanceURLtoString(m.Name)).Do()
		if err != nil && !strings.Contains(err.Error(), "was not found") && !strings.Contains(err.Error(), "Values must match") {
			return err
		}
		return nil
	})
	return procedure.Run()
}
