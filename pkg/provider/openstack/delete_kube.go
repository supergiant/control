package openstack

import (
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	floatingip "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// DeleteKube deletes a kubernetes cluster.
func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
	// Initialize steps
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Delete Kube",
		Model:  m,
		Action: action,
	}
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m)
	if err != nil {
		if ignoreErrors(err) {
			return nil
		}
		return err
	}
	// Fetch compute client.
	computeClient, err := openstack.NewComputeV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		if ignoreErrors(err) {
			return nil
		}
		return err
	}

	// Fetch network client.
	networkClient, err := openstack.NewNetworkV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		if ignoreErrors(err) {
			return nil
		}
		return err
	}

	if m.OpenStackConfig.PublicGatwayID != publicDisabled {
		procedure.AddStep("Destroying kubernetes Floating IP...", func() error {
			err := err
			floatIP, err := floatingip.Get(computeClient, m.OpenStackConfig.FloatingIPID).Extract()
			if err != nil {
				if ignoreErrors(err) {
					return nil
				}
				return err
			}
			// Disassociate Instance from floating IP
			disassociateOpts := floatingip.DisassociateOpts{
				FloatingIP: floatIP.IP,
			}

			err = floatingip.DisassociateInstance(computeClient, floatIP.InstanceID, disassociateOpts).ExtractErr()
			if err != nil {
				if ignoreErrors(err) {
					return nil
				}
				return err
			}
			// Delete the floating IP
			err = floatingip.Delete(networkClient, floatIP.ID).ExtractErr()
			if err != nil {
				if ignoreErrors(err) {
					return nil
				}
				return err
			}
			return nil
		})
	}

	procedure.AddStep("Destroying kubernetes master nodes...", func() error {
		err := err
		for _, node := range m.MasterNodes {
			err = servers.Delete(computeClient, node).ExtractErr()
			if err != nil {
				if ignoreErrors(err) {
					return nil
				}
				return err
			}

			return action.CancellableWaitFor("Wait for master delete", 20*time.Minute, 3*time.Second, func() (bool, error) {
				_, err := servers.Get(computeClient, node).Extract()
				if err == nil {
					return false, nil
				}
				return true, nil
			})

		}
		return nil
	})

	procedure.AddStep("Destroying keypair...", func() error {
		err := err
		err = keypairs.Delete(computeClient, m.OpenStackConfig.KeyPair).ExtractErr()
		if err != nil {
			if ignoreErrors(err) {
				return nil
			}
			return err
		}
		return nil
	})

	procedure.AddStep("Destroying kubernetes Router...", func() error {
		// Remove router interface
		_, err = routers.RemoveInterface(networkClient, m.OpenStackConfig.RouterID, routers.RemoveInterfaceOpts{
			SubnetID: m.OpenStackConfig.SubnetID,
		}).Extract()
		if err != nil {
			if ignoreErrors(err) {
				return nil
			}
			return err
		}
		// Delete router
		result := routers.Delete(networkClient, m.OpenStackConfig.RouterID)
		err = result.ExtractErr()
		if err != nil {
			if ignoreErrors(err) {
				return nil
			}
			return err
		}

		return nil
	})

	procedure.AddStep("Destroying kubernetes network...", func() error {
		// Delete network
		err := networks.Delete(networkClient, m.OpenStackConfig.NetworkID).ExtractErr()
		if err != nil {
			if ignoreErrors(err) {
				return nil
			}
			return err
		}
		return action.CancellableWaitFor("Wait for network delete", 20*time.Minute, 3*time.Second, func() (bool, error) {
			_, err := networks.Get(computeClient, m.OpenStackConfig.NetworkID).Extract()
			if err == nil {
				return false, nil
			}
			networks.Delete(networkClient, m.OpenStackConfig.NetworkID).ExtractErr()
			return true, nil
		})
	})

	return procedure.Run()
}
