package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
)

const DeleteTargetPoolStepName = "gce_delete_target_pool"

type DeleteTargetPoolStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteTargetPoolStep() *DeleteTargetPoolStep {
	return &DeleteTargetPoolStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPoolName string) (*compute.Operation, error) {
					return client.TargetPools.Delete(config.ServiceAccount.ProjectID, config.Region, targetPoolName).Do()
				},
			}, nil
		},
	}
}

func (s *DeleteTargetPoolStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteTargetPoolStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteTargetPoolStepName)
	}

	_, err = svc.deleteTargetPool(ctx, config.GCEConfig, config.GCEConfig.TargetPoolName)

	if err != nil {
		logrus.Errorf("Error deleting target pool %v", err)
	}

	return nil
}

func (s *DeleteTargetPoolStep) Name() string {
	return DeleteTargetPoolStepName
}

func (s *DeleteTargetPoolStep) Depends() []string {
	return nil
}

func (s *DeleteTargetPoolStep) Description() string {
	return "Delete target pool master nodes"
}

func (s *DeleteTargetPoolStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
