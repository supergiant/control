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
	"github.com/supergiant/control/pkg/util"
)

const StepCreateLoadBalancer = "create_load_balancer"

type LoadBalancerCreater interface {
	CreateLoadBalancerWithContext(aws.Context, *elb.CreateLoadBalancerInput, ...request.Option) (*elb.CreateLoadBalancerOutput, error)
}

type CreateLoadBalancerStep struct {
	getLoadBalancerService func(cfg steps.AWSConfig) (LoadBalancerCreater, error)
}

//InitCreateMachine adds the step to the registry
func InitCreateLoadBalancer(getELBFn GetELBFn) {
	steps.RegisterStep(StepCreateLoadBalancer, NewCreateLoadBalancerStep(getELBFn))
}

func NewCreateLoadBalancerStep(getELBFn GetELBFn) *CreateLoadBalancerStep {
	return &CreateLoadBalancerStep{
		getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerCreater, error) {

			elbInstance, err := getELBFn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepCreateLoadBalancer, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return elbInstance, nil
		},
	}
}

func (s *CreateLoadBalancerStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	svc, err := s.getLoadBalancerService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("error getting ELB service %v", err)
		return errors.Wrapf(err, "error getting ELB service %s",
			StepCreateLoadBalancer)
	}

	subnetsSlice := make([]*string, 0, len(cfg.AWSConfig.Subnets))

	for az := range cfg.AWSConfig.Subnets {
		subnet := cfg.AWSConfig.Subnets[az]
		subnetsSlice = append(subnetsSlice, aws.String(subnet))
	}

	externalLoadBalancerName := aws.String(util.CreateLBName(cfg.ClusterID, true))

	output, err := svc.CreateLoadBalancerWithContext(ctx, &elb.CreateLoadBalancerInput{
		Listeners: []*elb.Listener{
			{
				InstancePort:     aws.Int64(443),
				LoadBalancerPort: aws.Int64(443),
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: externalLoadBalancerName,
		Scheme:           aws.String("internet-facing"),
		SecurityGroups: []*string{
			aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		},
		Subnets: subnetsSlice,
		Tags: []*elb.Tag{
			{
				Key:   aws.String("ClusterID"),
				Value: aws.String(cfg.ClusterID),
			},
			{
				Key:   aws.String("ClusterName"),
				Value: aws.String(cfg.ClusterName),
			},
			{
				Key:   aws.String("Type"),
				Value: aws.String("external"),
			},
		},
	})

	if err != nil {
		logrus.Debugf("create external load balancer %v",
			err)
		return errors.Wrapf(err, "create load balancer %s", StepCreateLoadBalancer)
	}

	logrus.Infof("Created load external balancer %s with dns name %s", *externalLoadBalancerName, *output.DNSName)

	cfg.ExternalDNSName = *output.DNSName
	cfg.AWSConfig.ExternalLoadBalancerName = *externalLoadBalancerName

	internalLoadBalancerName := aws.String(util.CreateLBName(cfg.ClusterID, false))

	output, err = svc.CreateLoadBalancerWithContext(ctx, &elb.CreateLoadBalancerInput{
		Listeners: []*elb.Listener{
			{
				InstancePort:     aws.Int64(443),
				LoadBalancerPort: aws.Int64(443),
				Protocol:         aws.String("TCP"),
			},
			{
				InstancePort:     aws.Int64(2379),
				LoadBalancerPort: aws.Int64(2379),
				Protocol:         aws.String("TCP"),
			},
			{
				InstancePort:     aws.Int64(2380),
				LoadBalancerPort: aws.Int64(2380),
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: internalLoadBalancerName,
		Scheme:           aws.String("internal"),
		SecurityGroups: []*string{
			aws.String(cfg.AWSConfig.MastersSecurityGroupID),
			aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		},
		Subnets: subnetsSlice,
		Tags: []*elb.Tag{
			{
				Key:   aws.String("ClusterID"),
				Value: aws.String(cfg.ClusterID),
			},
			{
				Key:   aws.String("ClusterName"),
				Value: aws.String(cfg.ClusterName),
			},
			{
				Key:   aws.String("Type"),
				Value: aws.String("internal"),
			},
		},
	})

	if err != nil {
		logrus.Debugf("create internal load balancer %v",
			err)
		return errors.Wrapf(err, "create internal load balancer %s", StepCreateLoadBalancer)
	}

	logrus.Infof("Created load internal balancer %s with dns name %s", *externalLoadBalancerName, *output.DNSName)

	cfg.InternalDNSName = *output.DNSName
	cfg.AWSConfig.InternalLoadBalancerName = *internalLoadBalancerName

	return nil
}

func (s *CreateLoadBalancerStep) Name() string {
	return StepCreateLoadBalancer
}

func (s *CreateLoadBalancerStep) Description() string {
	return "Create ELB external and internal load balancers for master nodes"
}

func (s *CreateLoadBalancerStep) Depends() []string {
	return []string{StepCreateSubnets, StepCreateSecurityGroups}
}

func (s *CreateLoadBalancerStep) Rollback(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	return nil
}
