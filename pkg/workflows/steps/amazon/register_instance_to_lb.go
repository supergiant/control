package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const RegisterInstanceStepName = "register_instance"

type LoadBalancerRegister interface {
	RegisterInstancesWithLoadBalancerWithContext(aws.Context, *elb.RegisterInstancesWithLoadBalancerInput, ...request.Option) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
}

type RegisterInstanceStep struct {
	getLoadBalancerService func(cfg steps.AWSConfig) (LoadBalancerRegister, error)
}

//InitCreateMachine adds the step to the registry
func InitRegisterInstance(getELBFn GetELBFn) {
	steps.RegisterStep(RegisterInstanceStepName, NewRegisterInstanceStep(getELBFn))
}

func NewRegisterInstanceStep(getELBFn GetELBFn) *RegisterInstanceStep {
	return &RegisterInstanceStep{
		getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerRegister, error) {

			elbInstance, err := getELBFn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					RegisterInstanceStepName, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return elbInstance, nil
		},
	}
}

func (s *RegisterInstanceStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	svc, err := s.getLoadBalancerService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("error getting ELB service %v", err)
		return errors.Wrapf(err, "error getting ELB service %s",
			RegisterInstanceStepName)
	}

	logrus.Infof("Register instance Name: %s ID: %s to external load balancer: %s",
		cfg.Node.Name, cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName)
	_, err = svc.RegisterInstancesWithLoadBalancerWithContext(ctx, &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(cfg.AWSConfig.ExternalLoadBalancerName),
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(cfg.Node.ID),
			},
		},
	})

	if err != nil {
		logrus.Errorf("error registering instance %s to external loadbalancer %s %v", cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName, err)
		return errors.Wrapf(err, "registering instance %s to load balancer Load balancer %s %s",
			cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName,
			DeleteLoadBalancerStepName)
	}

	logrus.Infof("Register instance Name: %s ID: %s to internal load balancer: %s",
		cfg.Node.Name, cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName)
	_, err = svc.RegisterInstancesWithLoadBalancerWithContext(ctx, &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(cfg.AWSConfig.InternalLoadBalancerName),
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(cfg.Node.ID),
			},
		},
	})

	if err != nil {
		logrus.Errorf("error registering instance %s to internal loadbalancer %s %v", cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName, err)
		return errors.Wrapf(err, "registering instance %s to internal load balancer Load balancer %s %s",
			cfg.Node.ID, cfg.AWSConfig.ExternalLoadBalancerName,
			DeleteLoadBalancerStepName)
	}

	return nil
}

func (s *RegisterInstanceStep) Name() string {
	return RegisterInstanceStepName
}

func (s *RegisterInstanceStep) Description() string {
	return "Register node to external and internal Load balancers"
}

func (s *RegisterInstanceStep) Depends() []string {
	return nil
}

func (s *RegisterInstanceStep) Rollback(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	return nil
}
