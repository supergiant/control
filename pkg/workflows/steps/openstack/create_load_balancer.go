package openstack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas/monitors"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateLoadBalancerStepName = "create_load_balancer"
	StatusOnline               = "ONLINE"
	StatusOffline              = "OFFLINE"
	StatusActive               = "ACTIVE"
)

type CreateLoadBalancer struct {
	attemptCount int
	timeout      time.Duration

	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateLoadBalancer() *CreateLoadBalancer {
	return &CreateLoadBalancer{
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

	lbOpts := loadbalancers.CreateOpts{
		VipNetworkID: config.OpenStackConfig.NetworkID,
		Flavor:       config.OpenStackConfig.FlavorName,
	}

	loadBalancer, err := loadbalancers.Create(loadBalancerClient, lbOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create load balancer")
	}

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

	listenerOpts := listeners.CreateOpts{
		Name:         fmt.Sprintf("listener-%s", config.ClusterID),
		Protocol:     "HTTP",
		ProtocolPort: 443,
	}

	listener, err := listeners.Create(loadBalancerClient, listenerOpts).Extract()

	if err != nil {
		return errors.Wrapf(err,"create listener")
	}

	for i := 0; i < s.attemptCount; i++ {
		listener, err = listeners.Get(loadBalancerClient, listener.ID).Extract()

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

	poolOpts := pools.CreateOpts{
		Name: fmt.Sprintf("pool-%s", config.ClusterID),
		Protocol: "HTTP",
		LoadbalancerID: config.OpenStackConfig.LoadBalancerID,
		ListenerID:     config.OpenStackConfig.ListenerID,
	}

	pool, err := pools.Create(loadBalancerClient, poolOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create pool")
	}

	for i := 0; i < s.attemptCount; i++ {
		pool, err = pools.Get(loadBalancerClient, pool.ID).Extract()

		if err == nil && pool.ProvisioningStatus == StatusActive {
			break
		}

		time.Sleep(s.timeout)
	}

	if err != nil {
		return errors.Wrapf(err, "error while getting pool active")
	}

	if pool.ProvisioningStatus == StatusActive {
		return errors.Wrapf(err, "pool still is not active")
	}

	config.OpenStackConfig.PoolID = pool.ID

	healthOpts := monitors.CreateOpts{
		Type: monitors.TypeHTTPS,
		HTTPMethod: http.MethodGet,
		ExpectedCodes: "200-202",
		MaxRetries: 3,
		Delay: 20,
		Timeout: 10,
		URLPath: "/healthz",
	}

	healthCheck, err := monitors.Create(loadBalancerClient, healthOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create health check")
	}

	config.OpenStackConfig.HealthCheckID = healthCheck.ID

	memberOpts := pools.CreateMemberOpts{
		Address: config.Node.PrivateIp,
		ProtocolPort: 443,
		SubnetID: config.OpenStackConfig.SubnetID,
		Name: fmt.Sprintf("member-%s", config.Node.ID),
	}

	_, err = pools.CreateMember(loadBalancerClient,  config.OpenStackConfig.PoolID, memberOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create member")
	}

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
