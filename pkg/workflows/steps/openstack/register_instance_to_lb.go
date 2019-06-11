package openstack

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	RegisterInstancetoLBStepName = "register_instance_to_lb"
)

type RegisterInstancetoLBStep struct {
	attemptCount int
	timeout      time.Duration

	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewRegisterInstancetoLBStep() *RegisterInstancetoLBStep {
	return &RegisterInstancetoLBStep{
		attemptCount: 30,
		timeout:      time.Second * 10,
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

func (s *RegisterInstancetoLBStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", RegisterInstancetoLBStepName)
	}

	loadBalancerClient, err := openstack.NewLoadBalancerV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", RegisterInstancetoLBStepName)
	}

	memberOpts := pools.CreateMemberOpts{
		Address:      config.Node.PrivateIp,
		ProtocolPort: 443,
		SubnetID:     config.OpenStackConfig.SubnetID,
		Name:         fmt.Sprintf("member-%s", config.Node.ID),
	}

	_, err = pools.CreateMember(loadBalancerClient, config.OpenStackConfig.PoolID, memberOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create member")
	}

	return nil
}

func (s *RegisterInstancetoLBStep) Name() string {
	return RegisterInstancetoLBStepName
}

func (s *RegisterInstancetoLBStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *RegisterInstancetoLBStep) Description() string {
	return "Register instance to load balancer"
}

func (s *RegisterInstancetoLBStep) Depends() []string {
	return nil
}
