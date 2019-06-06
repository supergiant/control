package openstack

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/util"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
)

const CreateMachineStepName = "create_machine"

type CreateMachineStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateMachineStep() *CreateMachineStep {
	return &CreateMachineStep{
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
			{UUID: config.OpenStackConfig.NetworkID},
		},
		Metadata: map[string]string{
			"KubernetesCluster": config.ClusterName,
			"Role": util.MakeRole(config.IsMaster),
			clouds.TagClusterID: config.ClusterID,

		},
	}

	server, err := servers.Create(computeClient, serverCreateOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create machine caused")
	}

	// TODO(stgleb): wait for machine to become active
	machine := &model.Machine{
		Name: server.Name,
	}

	return nil
}

func (s *CreateMachineStep) Name() string {
	return CreateNetworkStepName
}

func (s *CreateMachineStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateMachineStep) Description() string {
	return "Create machine"
}

func (s *CreateMachineStep) Depends() []string {
	return []string{kubeadm.StepName}
}
