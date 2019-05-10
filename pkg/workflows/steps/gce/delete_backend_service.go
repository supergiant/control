package gce

import (
	"context"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteBackendServicStepName = "gce_delete_backend_service"

type DeleteBackendServiceStep struct {
	Timeout      time.Duration
	AttemptCount int

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteBackendServiceStep() (*DeleteBackendServiceStep, error) {
	return &DeleteBackendServiceStep{
		Timeout:      time.Second * 10,
		AttemptCount: 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

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

	var err error
	logrus.Debugf("Step %s", DeleteBackendServicStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteBackendServicStepName)
	}

	timeout := s.Timeout

	for i := 0; i < s.AttemptCount; i++ {
		_, err = svc.deleteBackendService(ctx, config.GCEConfig, config.GCEConfig.BackendServiceName)

		if err != nil {
			logrus.Errorf("Error deleting backend service rule %v", err)
		} else {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
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
