package gce

import (
	"context"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
	"time"
)

const DeleteForwardingRulesStepName = "gce_delete_forwarding_rules"

type DeleteForwardingRulesStep struct {
	attemptCount int
	timeout time.Duration

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteForwardingRulesStep() *DeleteForwardingRulesStep {
	return &DeleteForwardingRulesStep{
		attemptCount: 6,
		timeout: time.Second * 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteForwardingRule: func(ctx context.Context, config steps.GCEConfig, forwardingRuleName string) (*compute.Operation, error) {
					return client.ForwardingRules.Delete(config.ServiceAccount.ProjectID, config.Region, forwardingRuleName).Do()
				},
				getForwardingRule: func(ctx context.Context, config steps.GCEConfig, name string) (*compute.ForwardingRule, error) {
					return client.ForwardingRules.Get(config.ServiceAccount.ProjectID, config.Region, name).Do()
				},
			}, nil
		},
	}
}

func (s *DeleteForwardingRulesStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", DeleteForwardingRulesStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteForwardingRulesStepName)
	}

	_, err = svc.deleteForwardingRule(ctx, config.GCEConfig, config.GCEConfig.ExternalForwardingRuleName)

	if err != nil {
		logrus.Errorf("Error deleting external forwarding rule  %s %v", config.GCEConfig.ExternalForwardingRuleName, err)
	}

	for i := 0;i < s.attemptCount; i++ {
		_, err = svc.getForwardingRule(ctx, config.GCEConfig, config.GCEConfig.ExternalForwardingRuleName)

		if isNotFound(err) {
			break
		}

		logrus.Debugf("Forwarding rule %s still exists retry in %v",
			config.GCEConfig.ExternalForwardingRuleName, s.timeout)
		time.Sleep(s.timeout)
	}

	logrus.Debugf("Forwarding rule %s has been deleted",
		config.GCEConfig.ExternalForwardingRuleName)

	_, err = svc.deleteForwardingRule(ctx, config.GCEConfig, config.GCEConfig.InternalForwardingRuleName)

	if err != nil {
		logrus.Errorf("Error deleting internal forwarding rule %s rule %v", config.GCEConfig.InternalForwardingRuleName, err)
	}

	for i := 0;i < s.attemptCount; i++ {
		_, err = svc.getForwardingRule(ctx, config.GCEConfig, config.GCEConfig.InternalForwardingRuleName)

		if isNotFound(err) {
			break
		}

		logrus.Debugf("Forwarding rule %s still exists retry in %v",
			config.GCEConfig.InternalForwardingRuleName, s.timeout)
		time.Sleep(s.timeout)
	}

	logrus.Debugf("Forwarding rule %s has been deleted",
		config.GCEConfig.InternalForwardingRuleName)

	return nil
}

func (s *DeleteForwardingRulesStep) Name() string {
	return DeleteForwardingRulesStepName
}

func (s *DeleteForwardingRulesStep) Depends() []string {
	return nil
}

func (s *DeleteForwardingRulesStep) Description() string {
	return "Delete forwarding rules"
}

func (s *DeleteForwardingRulesStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
