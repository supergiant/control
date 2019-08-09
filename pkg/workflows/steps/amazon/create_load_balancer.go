package amazon

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateLoadBalancer = "create_load_balancer"

type LoadBalancerCreater interface {
	CreateLoadBalancerWithContext(aws.Context, *elb.CreateLoadBalancerInput, ...request.Option) (*elb.CreateLoadBalancerOutput, error)
	ConfigureHealthCheck(*elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error)
}

var (
	healthyThreshold   int64 = 2
	unhealthyThreshold int64 = 10
	checkInternal      int64 = 10
	checkTimeout       int64 = 5
)

type CreateLoadBalancerStep struct {
	timeout                time.Duration
	attemptCount           int
	getLoadBalancerService func(cfg steps.AWSConfig) (LoadBalancerCreater, error)
}

//InitCreateMachine adds the step to the registry
func InitCreateLoadBalancer(getELBFn GetELBFn) {
	steps.RegisterStep(StepCreateLoadBalancer, NewCreateLoadBalancerStep(getELBFn))
}

func NewCreateLoadBalancerStep(getELBFn GetELBFn) *CreateLoadBalancerStep {
	return &CreateLoadBalancerStep{
		timeout:      time.Second * 10,
		attemptCount: 120,
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

	if cfg.AWSConfig.ExternalLoadBalancerName == "" {
		externalLoadBalancerName := aws.String(util.CreateLBName(cfg.Kube.ID, true))
		output, err := svc.CreateLoadBalancerWithContext(ctx, &elb.CreateLoadBalancerInput{
			Listeners: []*elb.Listener{
				{
					InstancePort:     aws.Int64(cfg.Kube.APIServerPort),
					LoadBalancerPort: aws.Int64(cfg.Kube.APIServerPort),
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
					Key:   aws.String(clouds.TagClusterID),
					Value: aws.String(cfg.Kube.ID),
				},
				{
					Key:   aws.String("ClusterName"),
					Value: aws.String(cfg.Kube.Name),
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

		cfg.Kube.ExternalDNSName = *output.DNSName
		cfg.AWSConfig.ExternalLoadBalancerName = *externalLoadBalancerName
	}

	if cfg.AWSConfig.InternalLoadBalancerName == "" {
		internalLoadBalancerName := aws.String(util.CreateLBName(cfg.Kube.ID, false))

		output, err := svc.CreateLoadBalancerWithContext(ctx, &elb.CreateLoadBalancerInput{
			Listeners: []*elb.Listener{
				{
					InstancePort:     aws.Int64(cfg.Kube.APIServerPort),
					LoadBalancerPort: aws.Int64(cfg.Kube.APIServerPort),
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
					Key:   aws.String(clouds.TagClusterID),
					Value: aws.String(cfg.Kube.ID),
				},
				{
					Key:   aws.String("ClusterName"),
					Value: aws.String(cfg.Kube.Name),
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

		logrus.Infof("Created load internal balancer %s with dns name %s", *internalLoadBalancerName, *output.DNSName)

		cfg.Kube.InternalDNSName = *output.DNSName
		cfg.AWSConfig.InternalLoadBalancerName = *internalLoadBalancerName
	}

	for i := 0; i < s.attemptCount; i++ {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return errors.Wrap(err, "context has finished with")
			}
			return nil
		default:
			_, err = net.LookupIP(cfg.Kube.InternalDNSName)

			if err == nil {
				break
			}
			time.Sleep(s.timeout)
			logrus.Debugf("connect to load balancer %s with %v", cfg.Kube.InternalDNSName, err)
		}
	}

	if err != nil {
		return errors.Wrap(err, "error waiting for load balancer to come up")
	}

	logrus.Debugf("Configure health check for %s", cfg.AWSConfig.ExternalLoadBalancerName)
	healthCheckInput := &elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(cfg.AWSConfig.ExternalLoadBalancerName),
		HealthCheck: &elb.HealthCheck{
			HealthyThreshold:   &healthyThreshold,
			UnhealthyThreshold: &unhealthyThreshold,
			Interval:           &checkInternal,
			Timeout:            &checkTimeout,
			Target:             aws.String(fmt.Sprintf("HTTPS:%d/healthz", cfg.Kube.APIServerPort)),
		},
	}

	if _, err := svc.ConfigureHealthCheck(healthCheckInput); err != nil {
		logrus.Errorf("error configuring health check for %v  %s", err, cfg.AWSConfig.ExternalLoadBalancerName)
	}

	logrus.Debugf("Configure health check for %s", cfg.AWSConfig.InternalLoadBalancerName)
	healthCheckInput = &elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(cfg.AWSConfig.InternalLoadBalancerName),
		HealthCheck: &elb.HealthCheck{
			HealthyThreshold:   &healthyThreshold,
			UnhealthyThreshold: &unhealthyThreshold,
			Interval:           &checkInternal,
			Timeout:            &checkTimeout,
			Target:             aws.String(fmt.Sprintf("HTTPS:%d/healthz", cfg.Kube.APIServerPort)),
		},
	}

	if _, err := svc.ConfigureHealthCheck(healthCheckInput); err != nil {
		logrus.Errorf("error configuring health check for %v %s", err, cfg.AWSConfig.InternalLoadBalancerName)
	}

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
