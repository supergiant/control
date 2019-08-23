package openstack

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreatePoolStepName = "create_pool"
)

type CreatePoolStep struct {
	attemptCount int
	timeout      time.Duration

	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreatePoolStep() *CreatePoolStep {
	return &CreatePoolStep{
		attemptCount: 30,
		timeout:      time.Second * 10,
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

func (s *CreatePoolStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreatePoolStepName)
	}

	loadBalancerClient, err := openstack.NewLoadBalancerV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", CreatePoolStepName)
	}

	poolOpts := pools.CreateOpts{
		Name:           fmt.Sprintf("pool-%s", config.Kube.ID),
		Protocol:       pools.ProtocolHTTPS,
		ListenerID:     config.OpenStackConfig.ListenerID,
		LBMethod:       pools.LBMethodSourceIp,
	}

	pool, err := pools.Create(loadBalancerClient, poolOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create pool")
	}

	logrus.Debugf("Wait until pool become active")
	for i := 0; i < s.attemptCount; i++ {
		pool, err = pools.Get(loadBalancerClient, pool.ID).Extract()

		logrus.Debugf("Pool %s provisioning status %s",
			pool.ID, pool.OperatingStatus)

		if err == nil && pool.ProvisioningStatus == StatusActive {
			break
		}

		time.Sleep(s.timeout)
	}

	if err != nil {
		return errors.Wrapf(err, "error while getting pool active")
	}

	if pool.ProvisioningStatus != StatusActive {
		return errors.Wrapf(err, "pool still is not active")
	}

	config.OpenStackConfig.PoolID = pool.ID

	return nil
}

func (s *CreatePoolStep) Name() string {
	return CreatePoolStepName
}

func (s *CreatePoolStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreatePoolStep) Description() string {
	return "Create pool"
}

func (s *CreatePoolStep) Depends() []string {
	return nil
}
