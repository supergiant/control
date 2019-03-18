package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
)

const DeleteLoadBalancerStepName = "delete_load_balancer"

type LoadBalancerDeleter interface {
	DeleteLoadBalancer(*elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error)
}

type DeleteLoadBalancerStep struct {
	getLoadBalancerService func(cfg steps.AWSConfig) (LoadBalancerDeleter, error)
}

//InitCreateMachine adds the step to the registry
func InitDeleteLoadBalancer(getELBFn GetELBFn) {
	steps.RegisterStep(DeleteLoadBalancerStepName, NewDeleteLoadBalancerStep(getELBFn))
}

func NewDeleteLoadBalancerStep(getELBFn GetELBFn) *DeleteLoadBalancerStep {
	return &DeleteLoadBalancerStep{
		getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerDeleter, error) {

			elbInstance, err := getELBFn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					DeleteLoadBalancerStepName, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return elbInstance, nil
		},
	}
}

func (s *DeleteLoadBalancerStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	svc, err := s.getLoadBalancerService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("error getting ELB service %v", err)
		return errors.Wrapf(err, "error getting ELB service %s",
			DeleteLoadBalancerStepName)
	}

	_, err = svc.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(cfg.AWSConfig.LoadBalancerName),
	})

	if err != nil {
		logrus.Errorf("error deleting loadbalancer %s %v", cfg.AWSConfig.LoadBalancerName, err)
		return errors.Wrapf(err, "error deleteing Load balancer %s %s", cfg.AWSConfig.LoadBalancerName,
			DeleteLoadBalancerStepName)
	}

	return nil
}

func (s *DeleteLoadBalancerStep) Name() string {
	return StepCreateLoadBalancer
}

func (s *DeleteLoadBalancerStep) Description() string {
	return "Create ELB load balancer for master nodes"
}

func (s *DeleteLoadBalancerStep) Depends() []string {
	return nil
}

func (s *DeleteLoadBalancerStep) Rollback(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	return nil
}
