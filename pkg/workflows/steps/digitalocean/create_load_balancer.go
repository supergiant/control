package digitalocean

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type CreateLoadBalancerStep struct {
	Timeout     time.Duration
	Attempts    int
	getServices func(string) LoadBalancerService
}

func NewCreateLoadBalancerStep() *CreateLoadBalancerStep {
	return &CreateLoadBalancerStep{
		Timeout:  time.Second * 10,
		Attempts: 6,
		getServices: func(accessToken string) LoadBalancerService {
			client := digitaloceansdk.New(accessToken).GetClient()

			return client.LoadBalancers
		},
	}
}

func (s *CreateLoadBalancerStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	lbSvc := s.getServices(config.DigitalOceanConfig.AccessToken)

	req := &godo.LoadBalancerRequest{
		Name:   util.CreateLBName(config.ClusterID, true),
		Region: config.DigitalOceanConfig.Region,
		ForwardingRules: []godo.ForwardingRule{
			{
				EntryPort:      int(config.Kube.APIServerPort),
				EntryProtocol:  "TCP",
				TargetPort:     int(config.Kube.APIServerPort),
				TargetProtocol: "TCP",
			},
			{
				EntryPort:      2379,
				EntryProtocol:  "TCP",
				TargetPort:     2379,
				TargetProtocol: "TCP",
			},
			{
				EntryPort:      2380,
				EntryProtocol:  "TCP",
				TargetPort:     2380,
				TargetProtocol: "TCP",
			},
		},
		HealthCheck: &godo.HealthCheck{
			Protocol:               "TCP",
			Port:                   int(config.Kube.APIServerPort),
			CheckIntervalSeconds:   10,
			UnhealthyThreshold:     3,
			HealthyThreshold:       3,
			ResponseTimeoutSeconds: 10,
		},
		Tag: fmt.Sprintf("master-%s", config.ClusterID),
	}

	externalLoadBalancer, _, err := lbSvc.Create(ctx, req)

	if err != nil {
		logrus.Errorf("Error while creating external load balancer %v", err)
		return errors.Wrapf(err, "Error while creating external load balancer")
	}

	config.DigitalOceanConfig.ExternalLoadBalancerID = externalLoadBalancer.ID

	timeout := s.Timeout
	logrus.Infof("Wait until External load balancer %s become active", externalLoadBalancer.ID)
	for i := 0; i < s.Attempts; i++ {
		externalLoadBalancer, _, err = lbSvc.Get(ctx, config.DigitalOceanConfig.ExternalLoadBalancerID)

		if err == nil {
			logrus.Debugf("External Load balancer %s status %s",
				config.DigitalOceanConfig.ExternalLoadBalancerID, externalLoadBalancer.Status)
		}

		if err == nil && externalLoadBalancer.Status == StatusActive {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error while getting external load balancer %v", err)
		return errors.Wrapf(err, "Error while getting external load balancer")
	}

	if externalLoadBalancer.IP == "" {
		logrus.Errorf("External Load balancer IP must not be empty")
		return errors.New("External Load balancer IP must not be empty")
	}

	config.ExternalDNSName = externalLoadBalancer.IP

	req = &godo.LoadBalancerRequest{
		Name:   util.CreateLBName(config.ClusterID, false),
		Region: config.DigitalOceanConfig.Region,
		ForwardingRules: []godo.ForwardingRule{
			{
				EntryPort:      int(config.Kube.APIServerPort),
				EntryProtocol:  "TCP",
				TargetPort:     int(config.Kube.APIServerPort),
				TargetProtocol: "TCP",
			},
			{
				EntryPort:      2379,
				EntryProtocol:  "TCP",
				TargetPort:     2379,
				TargetProtocol: "TCP",
			},
			{
				EntryPort:      2380,
				EntryProtocol:  "TCP",
				TargetPort:     2380,
				TargetProtocol: "TCP",
			},
		},
		HealthCheck: &godo.HealthCheck{
			Protocol:               "TCP",
			Port:                   int(config.Kube.APIServerPort),
			CheckIntervalSeconds:   10,
			UnhealthyThreshold:     3,
			HealthyThreshold:       3,
			ResponseTimeoutSeconds: 10,
		},
		Tag: fmt.Sprintf("master-%s", config.ClusterID),
	}

	internalLoadBalancer, _, err := lbSvc.Create(ctx, req)

	if err != nil {
		logrus.Errorf("Error while creating internal load balancer %v", err)
		return errors.Wrapf(err, "Error while creating internal load balancer")
	}

	config.DigitalOceanConfig.InternalLoadBalancerID = internalLoadBalancer.ID
	logrus.Infof("Wait until Internal load balancer %s become active", internalLoadBalancer.ID)

	timeout = s.Timeout
	for i := 0; i < s.Attempts; i++ {
		internalLoadBalancer, _, err = lbSvc.Get(ctx, config.DigitalOceanConfig.InternalLoadBalancerID)

		if err == nil {
			logrus.Debugf("Internal Load balancer %s status %s",
				config.DigitalOceanConfig.InternalLoadBalancerID, internalLoadBalancer.Status)
		}

		if err == nil && internalLoadBalancer.Status == StatusActive {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error while getting internal load balancer %v", err)
		return errors.Wrapf(err, "Error while getting internal load balancer")
	}

	if internalLoadBalancer.IP == "" {
		logrus.Errorf("Internal Load balancer IP must not be empty")
		return errors.New("Internal Load balancer IP must not be empty")
	}

	config.InternalDNSName = internalLoadBalancer.IP

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
