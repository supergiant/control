package openstack

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	floatingip "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
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

// CreateKube creates a new kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {

	// Initialize steps
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
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

	// Default master count to 1
	if m.KubeMasterCount == 0 {
		m.KubeMasterCount = 1
	}

	// provision an etcd token
	url, err := etcdToken(strconv.Itoa(m.KubeMasterCount))
	if err != nil {
		return err
	}

	m.CustomFiles = fmt.Sprintf(`
  - path: /etc/kubernetes/ssl/cloud.conf
    permissions: 0600
    content: |
      [Global]
      auth-url = %s
      username = %s
      password = %s
      tenant-id = %s
      domain-id = %s
      domain-name = %s
      region = %s
      [LoadBalancer]
      subnet-id = %s`, m.CloudAccount.Credentials["identity_endpoint"],
		m.CloudAccount.Credentials["username"],
		m.CloudAccount.Credentials["password"],
		m.CloudAccount.Credentials["tenant_id"],
		m.CloudAccount.Credentials["domain_id"],
		m.CloudAccount.Credentials["domain_name"],
		m.OpenStackConfig.Region,
		m.OpenStackConfig.SubnetID,
	)

	m.KubeProviderString = `
         --cloud-provider=openstack \
         --cloud-config=/etc/kubernetes/ssl/cloud.conf \`

	m.ProviderString = `
          - --cloud-provider=openstack
          - --cloud-config=/etc/kubernetes/ssl/cloud.conf`

	err = p.Core.DB.Save(m)
	if err != nil {
		return err
	}
	// save the token
	m.ETCDDiscoveryURL = url

	// Procedures
	// Check for image.
	procedure.AddStep("Checking that a CoreOS image exists...", func() error {
		_, err = images.IDFromName(computeClient, "CoreOS")
		if err != nil {
			return err
		}
		return nil
	})

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
		err = p.Core.DB.Save(m)
		if err != nil {
			return err
		}
		return nil
	})

	// Router
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

	// Create Security Groups
	procedure.AddStep("Creating master security group...", func() error {
		secGroupOpts := secgroups.CreateOpts{
			Name:        m.Name + "-master-security-group",
			Description: "Security group for masters in Supergiant cluster " + m.Name,
		}
		group, err := secgroups.Create(computeClient, secGroupOpts).Extract()
		if err != nil {
			return err
		}
		m.OpenStackConfig.MasterSecurityGroupID = group.ID

		// TODO: It would be better if there was a single set of rules that each provider could use without excessive duplication!
		ruleOpts := secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      22,
			ToPort:        22,
			IPProtocol:    "TCP",
			CIDR:          "0.0.0.0/0",
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		ruleOpts = secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      443,
			ToPort:        443,
			IPProtocol:    "TCP",
			CIDR:          "0.0.0.0/0",
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		return nil
	})

	procedure.AddStep("Creating node security group...", func() error {
		// TODO: Refactor out duplication...
		secGroupOpts := secgroups.CreateOpts{
			Name:        m.Name + "-node-security-group",
			Description: "Security group for nodes in Supergiant cluster " + m.Name,
		}
		group, err := secgroups.Create(computeClient, secGroupOpts).Extract()
		if err != nil {
			return err
		}
		m.OpenStackConfig.NodeSecurityGroupID = group.ID

		// TODO: It would be better if there was a single set of rules that each provider could use without excessive duplication!
		ruleOpts := secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      22,
			ToPort:        22,
			IPProtocol:    "TCP",
			CIDR:          m.OpenStackConfig.PrivateSubnetRange,
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		ruleOpts = secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      443,
			ToPort:        443,
			IPProtocol:    "TCP",
			CIDR:          m.OpenStackConfig.PrivateSubnetRange,
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		ruleOpts = secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      10250,
			ToPort:        10250,
			IPProtocol:    "TCP",
			CIDR:          m.OpenStackConfig.PrivateSubnetRange,
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		ruleOpts = secgroups.CreateRuleOpts{
			ParentGroupID: group.ID,
			FromPort:      30000,
			ToPort:        40000,
			IPProtocol:    "TCP",
			CIDR:          m.OpenStackConfig.PrivateSubnetRange,
		}
		_, err = secgroups.CreateRule(computeClient, ruleOpts).Extract()
		if err != nil {
			return err
		}
		return nil
	})

	for i := 1; i <= m.KubeMasterCount; i++ {
		// Create master(s)
		count := strconv.Itoa(i)
		// Master
		procedure.AddStep("Creating Kubernetes Master Node "+count+"...", func() error {
			err := err
			// Build template

			// Master name
			name := m.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))

			m.MasterName = name

			mversion := strings.Split(m.KubernetesVersion, ".")
			masterUserdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/master.yaml")
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
			serverCreateOpts := servers.CreateOpts{
				ServiceClient: computeClient,
				Name:          name,
				FlavorName:    m.MasterNodeSize,
				ImageName:     m.OpenStackConfig.ImageName,
				UserData:      masterUserdata.Bytes(),
				Networks: []servers.Network{
					{UUID: m.OpenStackConfig.NetworkID},
				},
				SecurityGroups: []string{m.OpenStackConfig.MasterSecurityGroupID},
				Metadata:       map[string]string{"kubernetes-cluster": m.Name, "Role": "master"},
			}
			p.Core.Log.Debug(m.OpenStackConfig.ImageName)
			masterServer, err := servers.Create(computeClient, serverCreateOpts).Extract()
			if err != nil {
				return err
			}

			// Save serverID
			m.MasterID = masterServer.ID
			m.MasterNodes = append(m.MasterNodes, masterServer.ID)
			p.Core.DB.Save(m)
			// Wait for IP to be assigned.
			pNetwork := m.Name + "-network"
			duration := 2 * time.Minute
			interval := 10 * time.Second

			// TODO(stgleb): Context should be inherited from higher level context
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			waitErr := util.WaitFor(ctx, "Kubernetes Master IP asssign...", interval, func() (bool, error) {
				server, _ := servers.Get(computeClient, masterServer.ID).Extract()
				if server.Addresses[pNetwork] == nil {
					return false, nil
				}
				items := server.Addresses[pNetwork].([]interface{})
				for _, item := range items {
					itemMap := item.(map[string]interface{})
					m.MasterPrivateIP = itemMap["addr"].(string)
				}
				return true, nil
			})
			if waitErr != nil {
				return waitErr
			}

			return nil
		})
	}

	// Setup floading IP for master api
	if m.OpenStackConfig.PublicGatwayID != publicDisabled {

		procedure.AddStep("Waiting for Kubernetes Floating IP to create...", func() error {
			err := err
			// Lets keep trying to create
			var floatIP *floatingips.FloatingIP
			duration := 5 * time.Minute
			interval := 10 * time.Second

			// TODO(stgleb): Context should be inherited from higher level context
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			waitErr := util.WaitFor(ctx, "OpenStack floating IP creation", interval, func() (bool, error) {
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
			err = floatingip.AssociateInstance(computeClient, m.MasterID, associateOpts).ExtractErr()
			if err != nil {
				return err
			}

			m.MasterPublicIP = floatIP.FloatingIP
			return nil
		})
	}
	// Minion
	procedure.AddStep("Creating Kubernetes Worker Node...", func() error {
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

	procedure.AddStep("waiting for Kubernetes", func() error {
		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, 3*time.Second, func() (bool, error) {
			k8s := p.Core.K8S(m)
			k8sNodes, err := k8s.ListNodes("")
			if err != nil {
				return false, nil
			}
			return len(k8sNodes) > 0, nil
		})
	})

	return procedure.Run()
}
