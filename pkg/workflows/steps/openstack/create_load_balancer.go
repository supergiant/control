package openstack

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateLoadBalancerStepName = "create_load_balancer"
)

type CreateLoadBalancer struct {
	getClient    func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateLoadBalancer() *CreateLoadBalancer {
	return &CreateLoadBalancer{
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

func (s *CreateLoadBalancer) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateNetworkStepName)
	}

	loadBalancerClient, err := openstack.NewLoadBalancerV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", FindImageStepName)
	}

	opts := loadbalancers.CreateOpts{
		VipNetworkID: config.OpenStackConfig.NetworkID,
	}

	loadBalancer, err := loadbalancers.Create(loadBalancerClient, opts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create load balancer")
	}

	// TODO(stgleb): Wait for load balancer to become active
	loadBalancer.OperatingStatus

	listenOpts := listeners.CreateOpts{}
	listeners.Create(loadBalancerClient, listenOpts)

	return nil
}

func (s *CreateLoadBalancer) Name() string {
	return CreateLoadBalancerStepName
}

func (s *CreateLoadBalancer) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateLoadBalancer) Description() string {
	return "Create load balancer"
}

func (s *CreateLoadBalancer) Depends() []string {
	return nil
}
