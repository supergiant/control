package openstack

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateHealthCheckStepName = "create_health_check"
)

type CreateHealthCheckStep struct {
	attemptCount int
	timeout      time.Duration

	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateHealthCheckStep() *CreateHealthCheckStep {
	return &CreateHealthCheckStep{
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

func (s *CreateHealthCheckStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateHealthCheckStepName)
	}

	loadBalancerClient, err := openstack.NewLoadBalancerV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", CreateHealthCheckStepName)
	}

	healthOpts := monitors.CreateOpts{
		Type:          monitors.TypeHTTPS,
		HTTPMethod:    http.MethodGet,
		ExpectedCodes: "200-202",
		MaxRetries:    3,
		Delay:         20,
		Timeout:       10,
		URLPath:       "/healthz",
		PoolID:        config.OpenStackConfig.PoolID,
	}

	healthCheck, err := monitors.Create(loadBalancerClient, healthOpts).Extract()

	if err != nil {
		return errors.Wrapf(err, "create health check")
	}

	config.OpenStackConfig.HealthCheckID = healthCheck.ID

	return nil
}

func (s *CreateHealthCheckStep) Name() string {
	return CreateLoadBalancerStepName
}

func (s *CreateHealthCheckStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateHealthCheckStep) Description() string {
	return "Create health check"
}

func (s *CreateHealthCheckStep) Depends() []string {
	return nil
}
