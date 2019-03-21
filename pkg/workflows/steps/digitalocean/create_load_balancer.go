package digitalocean

import (
	"context"
	"io"
	"time"

	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/util"
	"fmt"
)

type CreateLoadBalancerStep struct {
	Timeout time.Duration

	getServices func(string) LoadBalancerService
}

func NewCreateLoadBalancerStep() *CreateLoadBalancerStep {
	return &CreateLoadBalancerStep{
		getServices: func(accessToken string) LoadBalancerService {
			client := digitaloceansdk.New(accessToken).GetClient()

			return client.LoadBalancers
		},
	}
}

func (s *CreateLoadBalancerStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	lbSvc := s.getServices(config.DigitalOceanConfig.AccessToken)

	req := &godo.LoadBalancerRequest{
		Name: util.CreateLBName(config.ClusterID, true),
		Region: config.DigitalOceanConfig.Region,
		ForwardingRules: []godo.ForwardingRule{
			{
				EntryPort: 443,
				EntryProtocol: "HTTPS",
				TargetPort: 443,
				TargetProtocol: "HTTPS",
				// NOTE(stgleb): Sticky sessions won't work with TLS passthrough
				// https://www.digitalocean.com/docs/networking/load-balancers/how-to/ssl-passthrough/
				TlsPassthrough: true,
			},
		},
		HealthCheck: &godo.HealthCheck{
			Protocol: "HTTP",
			Port: 443,
			Path: "/healthz",
			CheckIntervalSeconds: 60,
			UnhealthyThreshold: 3,
			ResponseTimeoutSeconds: 30,
		},
		Tag: fmt.Sprintf("master-%s", config.ClusterID),
	}

	externalLoadBalancer, _, err := lbSvc.Create(ctx, req)

	if err != nil {
		logrus.Errorf("Error while creating external load balancer %v", err)
		return errors.Wrapf(err, "Error while creating external load balancer")
	}

	config.ExternalDNSName = externalLoadBalancer.IP
	config.DigitalOceanConfig.ExternalLoadBalancerID = externalLoadBalancer.ID

	req = &godo.LoadBalancerRequest{
		Name: util.CreateLBName(config.ClusterID, false),
		Region: config.DigitalOceanConfig.Region,
		ForwardingRules: []godo.ForwardingRule{
			{
				EntryPort: 443,
				EntryProtocol: "HTTPS",
				TargetPort: 443,
				TargetProtocol: "HTTPS",
				// NOTE(stgleb): Sticky sessions won't work with TLS passthrough
				// https://www.digitalocean.com/docs/networking/load-balancers/how-to/ssl-passthrough/
				TlsPassthrough: true,
			},
		},
		HealthCheck: &godo.HealthCheck{
			Protocol: "HTTP",
			Port: 443,
			Path: "/healthz",
			CheckIntervalSeconds: 60,
			UnhealthyThreshold: 3,
			ResponseTimeoutSeconds: 30,
		},
		Tag: config.ClusterID,
	}

	internalLoadBalancer, _, err := lbSvc.Create(ctx, req)

	if err != nil {
		logrus.Errorf("Error while creating internal load balancer %v", err)
		return errors.Wrapf(err, "Error while creating internal load balancer")
	}

	config.InternalDNSName = internalLoadBalancer.IP
	config.DigitalOceanConfig.InternalLoadBalancerID = internalLoadBalancer.ID

	return nil
}

func (s *CreateLoadBalancerStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateLoadBalancerStep) Name() string {
	return CreateLoadBalancerStepName
}

func (s *CreateLoadBalancerStep) Depends() []string {
	return nil
}

func (s *CreateLoadBalancerStep) Description() string {
	return "Create load balancer in Digital Ocean"
}
