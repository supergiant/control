package gce

import (
	"fmt"
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateForwardingRulesStepName = "gce_create_forwarding_rules"

type CreateForwardingRules struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateForwardingRulesStep() (*CreateForwardingRules, error) {
	return &CreateForwardingRules{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config.ClientEmail,
				config.PrivateKey, config.TokenURI)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertForwardingRule: func(ctx context.Context, config steps.GCEConfig, rule *compute.ForwardingRule) (*compute.Operation, error) {
					return client.ForwardingRules.Insert(config.ProjectID, config.Region, rule).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateForwardingRules) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateTargetPullStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateTargetPullStepName)
	}

	externalForwardingRule := &compute.ForwardingRule{
		Name: fmt.Sprintf("ex-rule-%s", config.ClusterID),
		IPAddress: config.GCEConfig.ExternalIPAddress,
		LoadBalancingScheme:  "EXTERNAL",
		Description: "External forwarding rule to target pool",
		IpVersion: "IPV4",
		IPProtocol: "TCP",
		Target: config.GCEConfig.TargetPoolName,
		Ports: []string{"443"},
	}

	_, err = svc.insertForwardingRule(ctx, config.GCEConfig, externalForwardingRule)

	if err != nil {
		logrus.Errorf("Error creating external forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating external forwarding rule caused", CreateTargetPullStepName)
	}

	internalForwardingRule := &compute.ForwardingRule{
		Name: fmt.Sprintf("in-rule-%s", config.ClusterID),
		IPAddress: config.GCEConfig.InternalIPAddress,
		LoadBalancingScheme:  "INTERNAL",
		Description: "Internal forwarding rule to target pool",
		IpVersion: "IPV4",
		IPProtocol: "TCP",
		Ports: []string{"443"},
		// TODO(stgleb): Create backend service and backend
		BackendService: config.GCEConfig.BackendServiceName,
	}

	_, err = svc.insertForwardingRule(ctx, config.GCEConfig, internalForwardingRule)

	if err != nil {
		logrus.Errorf("Error creating internal forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating internal forwarding rule caused", CreateTargetPullStepName)
	}

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
