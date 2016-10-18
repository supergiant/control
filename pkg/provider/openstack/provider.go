package openstack

import (
	"bytes"
	"errors"
	"strings"
	"text/template"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) (*gophercloud.ProviderClient, error)
}

var (
	publicDisabled = "disabled"
)

// ValidateAccount Valitades Open Stack account info.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	_, err := p.Client(&model.Kube{CloudAccount: m})
	if err != nil {
		return err
	}
	return nil
}

// CreateKube creates a new DO kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {

	// Initialize steps
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	// Method vars
	masterName := m.Name + "-master"

	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewComputeV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}

	// Fetch network client.
	networkClient, err := openstack.NewNetworkV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}

	// Proceedures
	// Network
	procedure.AddStep("Creating Kubernetes Network...", func() error {
		err := err
		// Create network
		net, err := networks.Create(networkClient, networks.CreateOpts{
			Name:         m.Name + "-network",
			AdminStateUp: networks.Up,
		}).Extract()
		if err != nil {
			return err
		}
		// Save result
		m.OpenStackConfig.NetworkID = net.ID
		return nil
	})

	// Subnet
	procedure.AddStep("Creating Kubernetes Subnet...", func() error {
		err := err
		// Create subnet
		sub, err := subnets.Create(networkClient, subnets.CreateOpts{
			NetworkID:      m.OpenStackConfig.NetworkID,
			CIDR:           m.OpenStackConfig.PrivateSubnetRange,
			IPVersion:      subnets.IPv4,
			Name:           m.Name + "-subnet",
			DNSNameservers: []string{"8.8.8.8"},
		}).Extract()
		if err != nil {
			return err
		}
		// Save result
		m.OpenStackConfig.SubnetID = sub.ID
		return nil
	})

	// Network
	procedure.AddStep("Creating Kubernetes Router...", func() error {
		err := err
		// Create Router
		var opts routers.CreateOpts
		if m.OpenStackConfig.PublicGatwayID != publicDisabled {
			opts = routers.CreateOpts{
				Name:         m.Name + "-router",
				AdminStateUp: networks.Up,
				GatewayInfo: &routers.GatewayInfo{
					NetworkID: m.OpenStackConfig.PublicGatwayID,
				},
			}
		} else {
			opts = routers.CreateOpts{
				Name:         m.Name + "-router",
				AdminStateUp: networks.Up,
			}
		}
		router, err := routers.Create(networkClient, opts).Extract()
		if err != nil {
			return err
		}

		// interface our subnet to the new router.
		routers.AddInterface(networkClient, router.ID, routers.InterfaceOpts{
			SubnetID: m.OpenStackConfig.SubnetID,
		})
		m.OpenStackConfig.RouterID = router.ID
		return nil
	})

	// Master
	procedure.AddStep("Creating Kubernetes Master...", func() error {
		err := err
		// Build template
		masterUserdataTemplate, err := bindata.Asset("config/providers/openstack/master.yaml")
		if err != nil {
			return err
		}
		masterTemplate, err := template.New("master_template").Parse(string(masterUserdataTemplate))
		if err != nil {
			return err
		}
		var masterUserdata bytes.Buffer
		if err = masterTemplate.Execute(&masterUserdata, m); err != nil {
			return err
		}

		// Create Server
		masterServer, err := servers.Create(computeClient, servers.CreateOpts{
			Name:       masterName,
			FlavorName: m.MasterNodeSize,
			ImageName:  "CoreOS",
			UserData:   masterUserdata.Bytes(),
			Networks: []servers.Network{
				servers.Network{UUID: m.OpenStackConfig.NetworkID},
			},
			Metadata: map[string]string{"kubernetes-cluster": m.Name, "Role": "master"},
		}).Extract()
		if err != nil {
			return err
		}

		// Save serverID
		m.OpenStackConfig.MasterID = masterServer.ID

		// Wait for IP to be assigned.
		pNetwork := m.Name + "-network"
		duration := 2 * time.Minute
		interval := 10 * time.Second
		waitErr := util.WaitFor("Kubernetes Master IP asssign...", duration, interval, func() (bool, error) {
			server, _ := servers.Get(computeClient, masterServer.ID).Extract()
			if server.Addresses[pNetwork] == nil {
				return false, nil
			}
			items := server.Addresses[pNetwork].([]interface{})
			for _, item := range items {
				itemMap := item.(map[string]interface{})
				m.OpenStackConfig.MasterPrivateIP = itemMap["addr"].(string)
			}
			return true, nil
		})
		if waitErr != nil {
			return waitErr
		}

		return nil
	})

	// Setup floading IP for master api
	if m.OpenStackConfig.PublicGatwayID != publicDisabled {

		procedure.AddStep("Waiting for Kubernetes Floating IP to create...", func() error {
			err := err
			// Lets keep trying to create
			var floatIP *floatingips.FloatingIP
			duration := 5 * time.Minute
			interval := 10 * time.Second
			waitErr := util.WaitFor("OpenStack floating IP creation", duration, interval, func() (bool, error) {
				opts := floatingips.CreateOpts{
					FloatingNetworkID: m.OpenStackConfig.PublicGatwayID,
				}
				floatIP, err = floatingips.Create(networkClient, opts).Extract()
				if err != nil {
					if strings.Contains(err.Error(), "Quota exceeded for resources") {
						// Don't return error, just return false to indicate we should retry.
						return false, nil
					}
					// Else this is another more badder type of error
					return false, err
				}
				return true, nil
			})
			if waitErr != nil {
				return waitErr
			}
			// save results
			m.OpenStackConfig.FloatingIpID = floatIP.ID
			// Associate with master
			err = floatingip.AssociateInstance(computeClient, floatingip.AssociateOpts{
				ServerID:   m.OpenStackConfig.MasterID,
				FloatingIP: floatIP.FloatingIP,
			}).ExtractErr()
			if err != nil {
				return err
			}

			m.MasterPublicIP = floatIP.FloatingIP
			return nil
		})
	}
	// Minion
	procedure.AddStep("Creating Kubernetes Minion...", func() error {
		// Load Nodes to see if we've already created a minion
		// TODO -- I think we can get rid of a lot of this do-unless behavior if we
		// modify Procedure to save progess on Action (which is easy to implement).
		if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		}

		node := &model.Node{
			KubeName: m.Name,
			Kube:     m,
			Size:     m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})
	return procedure.Run()
}

// DeleteKube deletes a DO kubernetes cluster.
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
		return err
	}
	// Fetch compute client.
	computeClient, err := openstack.NewComputeV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}

	// Fetch network client.
	networkClient, err := openstack.NewNetworkV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}
	if m.OpenStackConfig.PublicGatwayID != publicDisabled {
		procedure.AddStep("Destroying kubernetes Floating IP...", func() error {
			err := err
			floatIP, err := floatingip.Get(computeClient, m.OpenStackConfig.FloatingIpID).Extract()
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					// it does not exist,
					return nil
				}
				return err
			}
			// Disassociate Instance from floating IP
			err = floatingip.DisassociateInstance(computeClient, floatingip.AssociateOpts{
				ServerID:   m.OpenStackConfig.MasterID,
				FloatingIP: floatIP.IP,
			}).ExtractErr()
			if err != nil {
				if strings.Contains(err.Error(), "field missing") {
					// it does not exist,
					return nil
				}
				return err
			}
			// Delete the floating IP
			err = floatingips.Delete(networkClient, floatIP.ID).ExtractErr()
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					// it does not exist,
					return nil
				}
				return err
			}
			return nil
		})
	}

	procedure.AddStep("Destroying kubernetes nodes...", func() error {
		err := err
		err = servers.Delete(computeClient, m.OpenStackConfig.MasterID).ExtractErr()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				// it does not exist,
				return nil
			}
			return err
		}
		return nil
	})

	procedure.AddStep("Destroying kubernetes Router...", func() error {
		// Remove router interface
		_, err = routers.RemoveInterface(networkClient, m.OpenStackConfig.RouterID, routers.InterfaceOpts{
			SubnetID: m.OpenStackConfig.SubnetID,
		}).Extract()
		if err != nil {
			if strings.Contains(err.Error(), "Expected HTTP") {
				// it does not exist,
				return nil
			}
			return err
		}
		// Delete router
		result := routers.Delete(networkClient, m.OpenStackConfig.RouterID)
		err = result.ExtractErr()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				// it does not exist,
				return nil
			}
			return err
		}

		return nil
	})

	procedure.AddStep("Destroying kubernetes network...", func() error {
		// Delete network
		result := networks.Delete(networkClient, m.OpenStackConfig.NetworkID)
		err = result.ExtractErr()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				// it does not exist,
				return nil
			}
			return err
		}
		return nil
	})

	return procedure.Run()
}

// CreateNode creates a new minion on DO kubernetes cluster.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
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
	m.Name = m.Kube.Name + "-minion-" + util.RandomString(5)
	// Build template
	minionUserdataTemplate, err := bindata.Asset("config/providers/openstack/minion.yaml")
	if err != nil {
		return err
	}
	minionTemplate, err := template.New("minion_template").Parse(string(minionUserdataTemplate))
	if err != nil {
		return err
	}
	var minionUserdata bytes.Buffer
	if err = minionTemplate.Execute(&minionUserdata, m); err != nil {
		return err
	}

	// Create server
	server, err := servers.Create(computeClient, servers.CreateOpts{
		Name:       m.Name,
		FlavorName: m.Size, // <- Do we need a minion node size? This will work for now.
		ImageName:  "CoreOS",
		UserData:   minionUserdata.Bytes(),
		Networks: []servers.Network{
			servers.Network{UUID: m.Kube.OpenStackConfig.NetworkID},
		},
		Metadata: map[string]string{"kubernetes-cluster": m.Kube.Name, "Role": "minion"},
	}).Extract()
	if err != nil {
		return err
	}
	// Save data
	m.ProviderID = server.Name
	m.ProviderCreationTimestamp = time.Now()

	// Wait for IP to be assigned.
	pNetwork := m.Kube.Name + "-network"
	duration := 2 * time.Minute
	interval := 10 * time.Second
	waitErr := util.WaitFor("Kubernetes Minion IP asssign...", duration, interval, func() (bool, error) {
		serverObj, _ := servers.Get(computeClient, server.ID).Extract()
		if serverObj.Addresses[pNetwork] == nil {
			return false, nil
		}
		items := serverObj.Addresses[pNetwork].([]interface{})
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			m.Name = itemMap["addr"].(string)
		}
		return true, nil
	})
	if waitErr != nil {
		return waitErr
	}

	return p.Core.DB.Save(m)
}

// DeleteNode deletes a minsion on a DO kubernetes cluster.
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
		if strings.Contains(err.Error(), "404") {
			// it does not exist,
			return nil
		}
		return err
	}
	return nil
}

// CreateVolume createss a Volume on DO for Kubernetes
func (p *Provider) CreateVolume(m *model.Volume, action *core.Action) error {
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewBlockStorageV1(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.Kube.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}
	// Create Volume
	vol, err := volumes.Create(computeClient, volumes.CreateOpts{
		Size:       m.Size,
		Name:       m.Name,
		VolumeType: m.Type,
	}).Extract()
	if err != nil {
		return err
	}

	m.ProviderID = vol.ID
	return p.Core.DB.Save(m)
}

func (p *Provider) KubernetesVolumeDefinition(m *model.Volume) *kubernetes.Volume {
	return &kubernetes.Volume{
		Name: m.Name,
		FlexVolume: &kubernetes.FlexVolume{
			Driver: "supergiant.io/digitalocean",
			FSType: "ext4",
			Options: map[string]string{
				"volumeID": m.ProviderID,
				"name":     m.Name,
			},
		},
	}
}

// ResizeVolume re-sizes volume on DO kubernetes cluster.
func (p *Provider) ResizeVolume(m *model.Volume, action *core.Action) error {
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewBlockStorageV1(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.Kube.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}
	// Make snapshot of old volume
	snap, err := snapshots.Create(computeClient, snapshots.CreateOpts{
		Name:     "resize-snapshot",
		VolumeID: m.ProviderID,
	}).Extract()
	if err != nil {
		return err
	}
	// Wait for snapshot to complete.
	duration := 30 * time.Minute
	interval := 10 * time.Second
	waitErr := util.WaitFor("Volume snapshot in progress", duration, interval, func() (bool, error) {
		err := err
		snapDeet, err := snapshots.Get(computeClient, snap.ID).Extract()
		if err != nil {
			return false, err
		}

		if snapDeet.Status == "error" {
			return false, errors.New("Snapshot creation failed... Aborting resize.")
		}
		if snapDeet.Status == "available" {
			return true, nil
		}
		return false, nil
	})

	if waitErr != nil {
		err = snapshots.Delete(computeClient, snap.ID).ExtractErr()
		if err != nil {
			return err
		}
		return waitErr
	}

	// Create Volume
	vol, err := volumes.Create(computeClient, volumes.CreateOpts{
		Size:       m.Size,
		Name:       m.Name,
		VolumeType: m.Type,
		SnapshotID: snap.ID,
	}).Extract()
	if err != nil {
		return err
	}
	m.ProviderID = vol.ID

	return nil
}

// WaitForVolumeAvailable waits for DO volume to become available.
func (p *Provider) WaitForVolumeAvailable(m *model.Volume, action *core.Action) error {
	return nil
}

// DeleteVolume deletes a DO volume.
func (p *Provider) DeleteVolume(m *model.Volume, action *core.Action) error {
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewBlockStorageV1(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.Kube.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}
	// Delete Volume
	err = volumes.Delete(computeClient, m.ProviderID).ExtractErr()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// it does not exist,
			return nil
		}
		return err
	}
	return nil
}

// CreateEntrypoint creates a new Load Balancer for Kubernetes in DO
func (p *Provider) CreateEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return nil
}

// DeleteEntrypoint deletes load balancer from DO.
func (p *Provider) DeleteEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return nil
}

func (p *Provider) CreateEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	return nil
}

func (p *Provider) DeleteEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

// Client creates the client for the provider.
func Client(kube *model.Kube) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: kube.CloudAccount.Credentials["identity_endpoint"],
		Username:         kube.CloudAccount.Credentials["username"],
		Password:         kube.CloudAccount.Credentials["password"],
		TenantID:         kube.CloudAccount.Credentials["tenant_id"],
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
