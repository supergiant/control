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
	if cfg.AWSConfig.LoadBalancerName != "" {
		logrus.Debugf("use load balancer %s",
			cfg.AWSConfig.LoadBalancerName)
		return nil
	} else {
		svc, err := s.getLoadBalancerService(cfg.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "error getting ELB service %s",
				StepCreateLoadBalancer)
		}

		subnetsSlice := make([]*string, 0, len(cfg.AWSConfig.Subnets))

		for _, subnet := range cfg.AWSConfig.Subnets {
			subnetsSlice = append(subnetsSlice, &subnet)
		}

		output, err := svc.CreateLoadBalancerWithContext(ctx, &elb.CreateLoadBalancerInput{
			Listeners: []*elb.Listener{
				{
					InstancePort:     aws.Int64(443),
					LoadBalancerPort: aws.Int64(443),
					Protocol:         aws.String("TCP"),
				},
			},
			LoadBalancerName: aws.String("load-balancer-" + cfg.ClusterID),
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
			},
		})

		if err != nil {
			logrus.Debugf("create load balancer %v",
				err)
			return nil
		}

		cfg.AWSConfig.ELBDNSName = *output.DNSName
	}

	return nil
}

func (s *CreateLoadBalancerStep) Name() string {
	return StepCreateLoadBalancer
}

func (s *CreateLoadBalancerStep) Description() string {
	return "Create ELB load balancer for master nodes"
}

func (s *CreateLoadBalancerStep) Depends() []string {
	return []string{StepCreateSubnets, StepCreateSecurityGroups}
}

func (s *CreateLoadBalancerStep) Rollback(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	return nil
}
