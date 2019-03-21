package digitalocean

import (
	"context"
	"io"
	"time"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type RegisterInstanceToLBStep struct {
	Timeout time.Duration

	getServices func(string) LoadBalancerService
}

func NewRegisterInstanceToLBStep() *RegisterInstanceToLBStep {
	return &RegisterInstanceToLBStep{
		getServices: func(accessToken string) LoadBalancerService {
			client := digitaloceansdk.New(accessToken).GetClient()

			return client.LoadBalancers
		},
	}
}

func (s *RegisterInstanceToLBStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	lbSvc := s.getServices(config.DigitalOceanConfig.AccessToken)

	instanceID, err := strconv.Atoi(config.Node.ID)

	if err != nil {
		logrus.Errorf("Error converting ID to int %v", err)
		return errors.Wrapf(err, "error converting ID to int")
	}

	_, err = lbSvc.AddDroplets(ctx, config.DigitalOceanConfig.ExternalLoadBalancerID, instanceID)

	if err != nil {
		logrus.Errorf("Error adding droplet %d to external load balancer %s %v",
			instanceID, config.DigitalOceanConfig.ExternalLoadBalancerID, err)
		return errors.Wrapf(err, "Error adding droplet %d to external load balancer %s",
			instanceID, config.DigitalOceanConfig.ExternalLoadBalancerID)
	}

	_, err = lbSvc.AddDroplets(ctx, config.DigitalOceanConfig.InternalLoadBalancerID, instanceID)

	if err != nil {
		logrus.Errorf("Error adding droplet %d to internal load balancer %s %v",
			instanceID, config.DigitalOceanConfig.InternalLoadBalancerID, err)
		return errors.Wrapf(err, "Error adding droplet %d to internal load balancer %s",
			instanceID, config.DigitalOceanConfig.InternalLoadBalancerID)
	}

	return nil
}

func (s *RegisterInstanceToLBStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *RegisterInstanceToLBStep) Name() string {
	return RegisterInstanceToLB
}

func (s *RegisterInstanceToLBStep) Depends() []string {
	return nil
}

func (s *RegisterInstanceToLBStep) Description() string {
	return "Register instance to load balancers in Digital Ocean"
}
