package openstack

import (
	"context"
	"io"
	"time"

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
				TenantName:       config.TenantName,
				DomainID:         config.DomainID,
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
		return errors.Wrapf(err, "step %s get compute client", CreateMachineStepName)
	}

	var securityGroupId string

	if config.IsMaster {
		securityGroupId = config.OpenStackConfig.MasterSecurityGroupId
	} else {
		securityGroupId = config.OpenStackConfig.WorkerSecurityGroupId
	}

	serverCreateOpts := servers.CreateOpts{
		ServiceClient: computeClient,
		Name:          util.MakeNodeName(config.Kube.Name, config.TaskID, config.IsMaster),
		FlavorName:    config.OpenStackConfig.FlavorName,
		ImageRef:      config.OpenStackConfig.ImageID,
		Networks: []servers.Network{
			{
				UUID: config.OpenStackConfig.NetworkID,
			},
		},
		SecurityGroups: []string{
			securityGroupId,
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
		server, _ = servers.Get(computeClient, server.ID).Extract()

		if server.Status == Active {
			break
		}

		time.Sleep(s.timeout)
	}

	var privateIP string

	if server.Addresses[config.OpenStackConfig.NetworkName] != nil {
		items := server.Addresses[config.OpenStackConfig.NetworkName].([]interface{})
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			privateIP = itemMap["addr"].(string)
		}
	}

	machine.PrivateIp = privateIP
	machine.State = model.MachineStateActive
	machine.ID = server.ID
	machine.Size = config.OpenStackConfig.FlavorName

	// Assign floating IP for master nodes
	if config.IsMaster {
		// Lets keep trying to create
		var floatIP *floatingips.FloatingIP

		computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
			Region: config.OpenStackConfig.Region,
		})

		if err != nil {
			return errors.Wrapf(err, "step %s get compute client", CreateNetworkStepName)
		}

		opts := floatingips.CreateOpts{
			Pool: config.OpenStackConfig.NetworkName,
		}
		floatIP, err = floatingips.Create(computeClient, opts).Extract()

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

		machine.PublicIp = floatIP.IP
	}

	config.Node = *machine

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
