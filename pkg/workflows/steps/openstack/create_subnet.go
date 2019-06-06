package openstack

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateSubnetStepName = "create_subnet"

type CreateSubnetStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateSubnetStep() *CreateSubnetStep {
	return &CreateSubnetStep{
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

func (s *CreateSubnetStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
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

	sub, err := subnets.Create(networkClient, subnets.CreateOpts{
		NetworkID:      config.OpenStackConfig.NetworkID,
		CIDR:           config.OpenStackConfig.SubnetIPRange,
		IPVersion:      gophercloud.IPv4,
		Name:           fmt.Sprintf("subnet-%s", config.ClusterID),
		DNSNameservers: []string{"8.8.8.8"},
	}).Extract()

	if err != nil {
		return errors.Wrapf(err, "create subnet caused")
	}

	// Save result
	config.OpenStackConfig.SubnetID = sub.ID

	return nil
}

func (s *CreateSubnetStep) Name() string {
	return CreateSubnetStepName
}

func (s *CreateSubnetStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateSubnetStep) Description() string {
	return "Create subnet"
}

func (s *CreateSubnetStep) Depends() []string {
	return []string{}
}
