package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateForwardingRulesStepName = "gce_create_forwarding_rules"

type CreateForwardingRules struct {
	timeout      time.Duration
	attemptCount int

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateForwardingRulesStep() *CreateForwardingRules {
	return &CreateForwardingRules{
		timeout:      time.Second * 10,
		attemptCount: 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

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
	logrus.Debugf("Step %s", CreateForwardingRulesStepName)
	if !config.IsBootstrap {
		logrus.Debugf("Skip step %s for bootstrap node %s", CreateForwardingRulesStepName, config.Node.Name)
		return nil
	}

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

	for i := 0; i < s.attemptCount; i++ {
		_, err = svc.insertForwardingRule(ctx, config.GCEConfig, externalForwardingRule)

		if err == nil {
			break
		}

		logrus.Debugf("Error external forwarding rule %v sleep for %v", err, timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error creating external forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating external forwarding rule caused", CreateForwardingRulesStepName)
	}

	logrus.Debugf("Created external forwarding rule %s link %s", exName, externalForwardingRule.SelfLink)
	config.GCEConfig.ExternalForwardingRuleName = exName

	inName := fmt.Sprintf("inrule-%s", config.ClusterID)
	internalForwardingRule := &compute.ForwardingRule{
		Name:                inName,
		IPAddress:           config.GCEConfig.InternalIPAddressLink,
		LoadBalancingScheme: "INTERNAL",
		Description:         "Internal forwarding rule to target pool",
		IPProtocol:          "TCP",
		Ports:               []string{"443"},
		BackendService:      config.GCEConfig.BackendServiceLink,
		Network:             config.GCEConfig.NetworkLink,
		Subnetwork:          config.GCEConfig.SubnetLink,
	}

	timeout = s.timeout

	for i := 0; i < s.attemptCount; i++ {
		_, err = svc.insertForwardingRule(ctx, config.GCEConfig, internalForwardingRule)

		if err == nil {
			break
		}

		logrus.Debugf("Error internal forwarding rule error %v sleep for %v", err, timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error creating internal forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating internal forwarding rule caused", CreateForwardingRulesStepName)
	}

	logrus.Debugf("Created internal forwarding rule %s", inName)
	config.GCEConfig.InternalForwardingRuleName = inName

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
