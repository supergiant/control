package openstack

import (
	"context"
	"io"
	"time"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/util"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateMachineStepName = "create_machine"
	Active                = "ACTIVE"
)

type CreateMachineStep struct {
	attemptCount int
	timeout      time.Duration
	getClient    func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateMachineStep() *CreateMachineStep {
	return &CreateMachineStep{
		attemptCount: 10,
		timeout:      time.Second * 30,
		getClient: func(config steps.OpenStackConfig) (client *gophercloud.ProviderClient, e error) {
			opts := gophercloud.AuthOptions{
				IdentityEndpoint: config.AuthURL,
				Username:         config.UserName,
				Password:         config.Password,
				TenantID:         config.TenantID,
				DomainID:         config.DomainID,
				DomainName:       config.DomainName,
			}

			client, err := openstack.AuthenticatedClient(opts)
			if err != nil {
				return nil, err
			}

			return client, nil
		},
	}
}

func (s *CreateMachineStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateNetworkStepName)
	}

	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", FindImageStepName)
	}

	serverCreateOpts := servers.CreateOpts{
		ServiceClient: computeClient,
		Name:          util.MakeNodeName(config.ClusterName, config.TaskID, config.IsMaster),
		FlavorName:    config.OpenStackConfig.FlavorName,
		ImageName:     config.OpenStackConfig.ImageName,
		Networks: []servers.Network{
			{
				UUID: config.OpenStackConfig.NetworkID,
			},
		},
		Metadata: map[string]string{
			"KubernetesCluster": config.ClusterName,
			"Role":              util.MakeRole(config.IsMaster),
			clouds.TagClusterID: config.ClusterID,
		},
	}

	server, err := servers.Create(computeClient, serverCreateOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create machine caused")
	}

	machine := &model.Machine{
		Name:  server.Name,
		State: model.MachineStateBuilding,
	}

	config.NodeChan() <- *machine

	for i := 0; i < s.attemptCount; i++ {
		server, _ := servers.Get(computeClient, server.ID).Extract()

		if server.Status == Active {
			break
		}

		time.Sleep(s.timeout)
	}

	var privateIP string

	items := server.Addresses[config.OpenStackConfig.NetworkName].([]interface{})
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		privateIP = itemMap["addr"].(string)
	}

	machine.PrivateIp = privateIP
	machine.State = model.MachineStateActive
	machine.ID = server.ID
	machine.Size = config.OpenStackConfig.FlavorName

	// Assign floating IP for master nodes
	if config.IsMaster {
		// Lets keep trying to create
		var floatIP *floatingips.FloatingIP

		networkClient, err := openstack.NewNetworkV2(client, gophercloud.EndpointOpts{
			Region: config.OpenStackConfig.Region,
		})

		if err != nil {
			return errors.Wrapf(err, "step %s get network client", CreateNetworkStepName)
		}

		opts := floatingips.CreateOpts{
			Pool: config.OpenStackConfig.RouterID,
		}
		floatIP, err = floatingips.Create(networkClient, opts).Extract()

		if err != nil {
			return errors.Wrapf(err, "create floating ip")
		}

		associateOpts := floatingips.AssociateOpts{
			FloatingIP: floatIP.IP,
		}
		err = floatingips.AssociateInstance(computeClient, machine.ID, associateOpts).ExtractErr()

		if err != nil {
			return errors.Wrapf(err, "associate instance %s with floating ip %s", machine.ID, floatIP.ID)
		}

		config.OpenStackConfig.FloatingIP = floatIP.IP
		config.OpenStackConfig.FloatingID = floatIP.ID
	}

	return nil
}

func (s *CreateMachineStep) Name() string {
	return CreateMachineStepName
}

func (s *CreateMachineStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateMachineStep) Description() string {
	return "Create machine"
}

func (s *CreateMachineStep) Depends() []string {
	return nil
}
