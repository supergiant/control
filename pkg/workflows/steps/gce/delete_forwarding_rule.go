package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteForwardingRulesStepName = "gce_delete_forwarding_rules"

type DeleteForwardingRulesStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteForwardingRulesStep() *DeleteForwardingRulesStep {
	return &DeleteForwardingRulesStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteForwardingRule: func(ctx context.Context, config steps.GCEConfig, forwardingRuleName string) (*compute.Operation, error) {
					return client.ForwardingRules.Delete(config.ServiceAccount.ProjectID, config.Region, forwardingRuleName).Do()
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

	_, err = svc.deleteForwardingRule(ctx, config.GCEConfig, config.GCEConfig.ForwardingRuleName)

	if err != nil {
		logrus.Errorf("Error deleting external forwarding rule %v", err)
		return errors.Wrapf(err, "%s deleting external forwarding %s rule caused",
			config.GCEConfig.ForwardingRuleName, DeleteForwardingRulesStepName)
	}

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
