package gce

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
	"fmt"
)

const CreateBackendServiceStepName = "gce_create_backend_service"

type CreateBackendServiceStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateBackendServiceStep() (*CreateAddressStep, error) {
	return &CreateAddressStep{
		Timeout:      time.Second * 10,
		AttemptCount: 6,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config.ClientEmail,
				config.PrivateKey, config.TokenURI)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertBackendService: func(ctx context.Context, config steps.GCEConfig, service *compute.BackendService) (*compute.Operation, error) {
					return client.BackendServices.Insert(config.ProjectID, service).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateBackendServiceStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateBackendServiceStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateBackendServiceStepName)
	}


	backendService := &compute.BackendService{
		Name: fmt.Sprintf("backendService-%s", config.ClusterID),
		Description: "Backend service for internal traffic",
		LoadBalancingScheme: "INTERNAL",
		Protocol: "TCP",
		Backends: []*compute.Backend{
			{
				Group: config.GCEConfig.InstanceGroup,
			},
		},
	}

	_, err = svc.insertBackendService(ctx, config.GCEConfig, backendService)

	if err != nil {
		logrus.Error("error creating backend service %v", err)
		return errors.Wrapf(err, "error creating backend service")
	}

	config.GCEConfig.BackendServiceName = backendService.Name

	return nil
}

func (s *CreateBackendServiceStep) Name() string {
	return CreateBackendServiceStepName
}

func (s *CreateBackendServiceStep) Depends() []string {
	return []string{CreateInstanceGroupStepName}
}

func (s *CreateBackendServiceStep) Description() string {
	return "Create backend service"
}

func (s *CreateBackendServiceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
