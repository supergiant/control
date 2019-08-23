package openstack

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
)

const CreateNetworkStepName = "create_network"

var (
	isExternal = true
)

type CreateNetworkStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateNetworkStep() *CreateNetworkStep {
	return &CreateNetworkStep{
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

func (s *CreateNetworkStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateNetworkStepName)
	}

	networkClient, err := openstack.NewNetworkV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get network client", CreateNetworkStepName)
	}

	net, err := networks.Create(networkClient, external.CreateOptsExt{
		CreateOptsBuilder: networks.CreateOpts{
			Name:         fmt.Sprintf("network-%s", config.Kube.ID),
			AdminStateUp: gophercloud.Enabled,
		},
		External: &isExternal,
	}).Extract()

	if err != nil {
		return errors.Wrapf(err, "create network error step %s", CreateNetworkStepName)
	}

	// Save network ID
	config.OpenStackConfig.NetworkID = net.ID
	config.OpenStackConfig.NetworkName = net.Name

	return nil
}

func (s *CreateNetworkStep) Name() string {
	return CreateNetworkStepName
}

func (s *CreateNetworkStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateNetworkStep) Description() string {
	return "Create network"
}

func (s *CreateNetworkStep) Depends() []string {
	return []string{kubeadm.StepName}
}
