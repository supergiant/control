package openstack

import (
	"context"
	"fmt"

	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateGateway = "create_gateway"

type CreateGatewayStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateGatewayStep() *CreateSubnetStep {
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

func (s *CreateGatewayStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
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


	var opts routers.CreateOpts
	opts = routers.CreateOpts{
		Name:         fmt.Sprintf("router-%s", config.ClusterID),
		AdminStateUp: gophercloud.Enabled,
		GatewayInfo: &routers.GatewayInfo{
			NetworkID: config.OpenStackConfig.NetworkID,
		},
	}

	router, err := routers.Create(networkClient, opts).Extract()
	if err != nil {
		return err
	}

	// interface our subnet to the new router.
	routers.AddInterface(networkClient, router.ID, routers.AddInterfaceOpts{
		SubnetID: config.OpenStackConfig.SubnetID,
	})
	config.OpenStackConfig.RouterID = router.ID

	return nil
}

func (s *CreateGatewayStep) Name() string {
	return CreateSubnetStepName
}

func (s *CreateGatewayStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateGatewayStep) Description() string {
	return "Create subnet"
}

func (s *CreateGatewayStep) Depends() []string {
	return []string{}
}
