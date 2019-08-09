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

const CreateRouterStepName = "create_router"

type CreateRouterStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateRouterStep() *CreateSubnetStep {
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

func (s *CreateRouterStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateRouterStepName)
	}

	networkClient, err := openstack.NewNetworkV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get network client", CreateRouterStepName)
	}

	var opts routers.CreateOpts
	opts = routers.CreateOpts{
		Name:         fmt.Sprintf("router-%s", config.Kube.ID),
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

func (s *CreateRouterStep) Name() string {
	return CreateSubnetStepName
}

func (s *CreateRouterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateRouterStep) Description() string {
	return "Create router"
}

func (s *CreateRouterStep) Depends() []string {
	return []string{}
}
