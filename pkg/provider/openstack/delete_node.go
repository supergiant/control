package openstack

import (
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// DeleteNode deletes a node on a kubernetes cluster.
func (p *Provider) DeleteNode(m *model.Node, action *core.Action) error {
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewComputeV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.Kube.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}

	err = servers.Delete(computeClient, m.ProviderID).ExtractErr()
	if err != nil {
		if ignoreErrors(err) {
			return nil
		}
		return err
	}

	return action.CancellableWaitFor("Wait for node delete", 20*time.Minute, 3*time.Second, func() (bool, error) {
		_, err := servers.Get(computeClient, m.ProviderID).Extract()
		if err == nil {
			return false, nil
		}
		return true, nil
	})
}
