package digitalocean

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type DeleteLoadBalancerStep struct {
	getServices func(string) LoadBalancerService
}

func NewDeleteLoadBalancerStep() *DeleteLoadBalancerStep {
	return &DeleteLoadBalancerStep{
		getServices: func(accessToken string) LoadBalancerService {
			client := digitaloceansdk.New(accessToken).GetClient()

			return client.LoadBalancers
		},
	}
}

func (s *DeleteLoadBalancerStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	lbSvc := s.getServices(config.DigitalOceanConfig.AccessToken)

	_, err := lbSvc.Delete(ctx, config.DigitalOceanConfig.ExternalLoadBalancerID)

	if err != nil {
		logrus.Errorf("Error deleting external load balancer %s %v", config.DigitalOceanConfig.ExternalLoadBalancerID, err)
		return errors.Wrapf(err, "Error deleting external load balancer %s", config.DigitalOceanConfig.ExternalLoadBalancerID)
	}

	_, err = lbSvc.Delete(ctx, config.DigitalOceanConfig.InternalLoadBalancerID)

	if err != nil {
		logrus.Errorf("Error deleting internal load balancer %s %v", config.DigitalOceanConfig.InternalLoadBalancerID, err)
		return errors.Wrapf(err, "Errordeleting internal load balancer %s", config.DigitalOceanConfig.InternalLoadBalancerID)
	}

	return nil
}

func (s *DeleteLoadBalancerStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteLoadBalancerStep) Name() string {
	return DeleteLoadBalancerStepName
}

func (s *DeleteLoadBalancerStep) Depends() []string {
	return nil
}

func (s *DeleteLoadBalancerStep) Description() string {
	return "Delete external and internal load balancers in Digital Ocean"
}
