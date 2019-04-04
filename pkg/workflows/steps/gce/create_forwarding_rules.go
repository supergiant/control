package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateForwardingRulesStepName = "gce_create_forwarding_rules"

type CreateForwardingRules struct {
	timeout time.Duration
	attemptCount int

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateForwardingRulesStep() *CreateForwardingRules {
	return &CreateForwardingRules{
		timeout: time.Second * 10,
		attemptCount: 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertForwardingRule: func(ctx context.Context, config steps.GCEConfig, rule *compute.ForwardingRule) (*compute.Operation, error) {
					return client.ForwardingRules.Insert(config.ServiceAccount.ProjectID, config.Region, rule).Do()
				},
			}, nil
		},
	}
}

func (s *CreateForwardingRules) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	// Skip this step for the rest of nodes
	if !config.IsBootstrap {
		return nil
	}

	logrus.Debugf("Step %s", CreateForwardingRulesStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateForwardingRulesStepName)
	}

	exName := fmt.Sprintf("exrule-%s", config.ClusterID)
	externalForwardingRule := &compute.ForwardingRule{
		Name:                exName,
		IPAddress:           config.GCEConfig.ExternalIPAddressLink,
		LoadBalancingScheme: "EXTERNAL",
		Description:         "External forwarding rule to target pool",
		IPProtocol:          "TCP",
		Target:              config.GCEConfig.TargetPoolLink,
	}


	timeout := s.timeout

	for i := 0;i < s.attemptCount; i++ {
		_, err = svc.insertForwardingRule(ctx, config.GCEConfig, externalForwardingRule)

		if err == nil {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error creating external forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating external forwarding rule caused", CreateForwardingRulesStepName)
	}

	logrus.Debugf("Created external forwarding rule %s", exName)
	config.GCEConfig.ForwardingRuleName = exName

	return nil
}

func (s *CreateForwardingRules) Name() string {
	return CreateForwardingRulesStepName
}

func (s *CreateForwardingRules) Depends() []string {
	return []string{CreateTargetPullStepName, CreateIPAddressStepName}
}

func (s *CreateForwardingRules) Description() string {
	return "Create forwarding rules to pass traffic to nodes"
}

func (s *CreateForwardingRules) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
