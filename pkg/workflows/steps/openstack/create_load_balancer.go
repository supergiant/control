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
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateLoadBalancerStepName = "openstack_create_load_balancer"
)

type CreateLoadBalancer struct {
	attemptCount int
	timeout      time.Duration

	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateLoadBalancer() *CreateLoadBalancer {
	return &CreateLoadBalancer{
		attemptCount: 60,
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

func (s *CreateLoadBalancer) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateLoadBalancerStepName)
	}

	loadBalancerClient, err := openstack.NewLoadBalancerV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", CreateLoadBalancerStepName)
	}

	lbOpts := loadbalancers.CreateOpts{
		Name:         fmt.Sprintf("load-balancer-%s", config.Kube.ID),
		AdminStateUp: gophercloud.Enabled,
		VipNetworkID: config.OpenStackConfig.NetworkID,
		VipSubnetID:  config.OpenStackConfig.SubnetID,
		Tags:         []string{config.Kube.ID, config.Kube.Name},
	}

	loadBalancer, err := loadbalancers.Create(loadBalancerClient, lbOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create load balancer")
	}

	logrus.Debugf("Wait for load balancer %s to become active", loadBalancer.ID)
	for i := 0; i < s.attemptCount; i++ {
		loadBalancer, err = loadbalancers.Get(loadBalancerClient, loadBalancer.ID).Extract()

		if err == nil && loadBalancer.OperatingStatus == StatusOnline {
			break
		}

		time.Sleep(s.timeout)
	}

	if err != nil {
		return errors.Wrapf(err, "error while getting load balancer online")
	}

	if loadBalancer.OperatingStatus == StatusOffline {
		return errors.Wrapf(err, "load balancer still offline")
	}

	config.OpenStackConfig.LoadBalancerID = loadBalancer.ID
	config.OpenStackConfig.LoadBalancerName = loadBalancer.Name

	// TODO(stgleb): Move it to separate step
	listenerOpts := listeners.CreateOpts{
		LoadbalancerID: loadBalancer.ID,
		Name:           fmt.Sprintf("listener-%s", config.Kube.ID),
		Protocol:       listeners.ProtocolHTTPS,
		ProtocolPort:   config.Kube.APIServerPort,
	}

	logrus.Debugf("create listener on load balancer %s", loadBalancer.ID)
	listener, err := listeners.Create(loadBalancerClient, listenerOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create listener")
	}

	logrus.Debugf("wait until listener %s becomes active", listener.ID)
	for i := 0; i < s.attemptCount; i++ {
		listener, err = listeners.Get(loadBalancerClient, listener.ID).Extract()

		logrus.Debugf("listener %s provisioning status  %s",
			listener.ID, listener.ProvisioningStatus)
		if err == nil && listener.ProvisioningStatus == StatusActive {
			break
		}

		time.Sleep(s.timeout)
	}

	if err != nil {
		return errors.Wrapf(err, "error while getting listener active")
	}

	if listener.ProvisioningStatus == StatusActive {
		return errors.Wrapf(err, "listener still is not active")
	}

	config.OpenStackConfig.ListenerID = listener.ID

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
