package openstack

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	floatingip "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
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
			AdminStateUp: gophercloud.Enabled,
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
			IPVersion:      gophercloud.IPv4,
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
				AdminStateUp: gophercloud.Enabled,
				GatewayInfo: &routers.GatewayInfo{
					NetworkID: m.OpenStackConfig.PublicGatwayID,
				},
			}
		} else {
			opts = routers.CreateOpts{
				Name:         m.Name + "-router",
				AdminStateUp: gophercloud.Enabled,
			}
		}
		router, err := routers.Create(networkClient, opts).Extract()
		if err != nil {
			return err
		}

		// interface our subnet to the new router.
		routers.AddInterface(networkClient, router.ID, routers.AddInterfaceOpts{
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
		//
		serverCreateOpts := servers.CreateOpts{
			ServiceClient: computeClient,
			Name:          masterName,
			FlavorName:    m.MasterNodeSize,
			ImageName:     m.OpenStackConfig.ImageName,
			UserData:      masterUserdata.Bytes(),
			Networks: []servers.Network{
				servers.Network{UUID: m.OpenStackConfig.NetworkID},
			},
			Metadata: map[string]string{"kubernetes-cluster": m.Name, "Role": "master"},
		}
		p.Core.Log.Debug(m.OpenStackConfig.ImageName)
		masterServer, err := servers.Create(computeClient, serverCreateOpts).Extract()
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
			m.OpenStackConfig.FloatingIPID = floatIP.ID
			// Associate with master
			associateOpts := floatingip.AssociateOpts{
				FloatingIP: floatIP.FloatingIP,
			}
			err = floatingip.AssociateInstance(computeClient, m.OpenStackConfig.MasterID, associateOpts).ExtractErr()
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
			floatIP, err := floatingip.Get(computeClient, m.OpenStackConfig.FloatingIPID).Extract()
			if err != nil {
				if strings.Contains(err.Error(), "404") {
					// it does not exist,
					return nil
				}
				return err
			}
			// Disassociate Instance from floating IP
			disassociateOpts := floatingip.DisassociateOpts{
				FloatingIP: floatIP.IP,
			}
			err = floatingip.DisassociateInstance(computeClient, floatIP.ID, disassociateOpts).ExtractErr()
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
		_, err = routers.RemoveInterface(networkClient, m.OpenStackConfig.RouterID, routers.RemoveInterfaceOpts{
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

	serverCreateOpts := servers.CreateOpts{
		ServiceClient: computeClient,
		Name:          m.Name,
		FlavorName:    m.Size,
		ImageName:     m.Kube.OpenStackConfig.ImageName,
		UserData:      minionUserdata.Bytes(),
		Networks: []servers.Network{
			servers.Network{UUID: m.Kube.OpenStackConfig.NetworkID},
		},
		Metadata: map[string]string{"kubernetes-cluster": m.Kube.Name, "Role": "minion"},
	}

	// Create server
	server, err := servers.Create(computeClient, serverCreateOpts).Extract()
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

func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.CreateLoadBalancer(m, action)
}

func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.UpdateLoadBalancer(m, action)
}

func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.DeleteLoadBalancer(m, action)
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
		DomainID:         kube.CloudAccount.Credentials["domain_id"],
		DomainName:       kube.CloudAccount.Credentials["domain_name"],
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
