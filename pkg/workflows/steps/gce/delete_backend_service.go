package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteBackendServicStepName = "gce_delete_backend_service"

type DeleteBackendServiceStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteBackendServiceStep() (*DeleteBackendServiceStep, error) {
	return &DeleteBackendServiceStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteBackendService: func(ctx context.Context, config steps.GCEConfig, backendServiceName string) (*compute.Operation, error) {
					return client.RegionBackendServices.Delete(config.ServiceAccount.ProjectID, config.Region, backendServiceName).Do()
				},
			}, nil
		},
	}, nil
}

func (s *DeleteBackendServiceStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteBackendServicStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteBackendServicStepName)
	}

	_, err = svc.deleteBackendService(ctx, config.GCEConfig, config.GCEConfig.BackendServiceName)

	if err != nil {
		logrus.Errorf("Error deleting backend service rule %v", err)
		return errors.Wrapf(err, "%s deleting backend service rule caused", DeleteBackendServicStepName)
	}

	return nil
}

func (s *DeleteBackendServiceStep) Name() string {
	return DeleteBackendServicStepName
}

func (s *DeleteBackendServiceStep) Depends() []string {
	return nil
}

func (s *DeleteBackendServiceStep) Description() string {
	return "Delete backend service"
}

func (s *DeleteBackendServiceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
