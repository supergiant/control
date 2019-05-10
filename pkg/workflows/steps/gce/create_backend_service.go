package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateBackendServiceStepName = "gce_create_backend_service"

type CreateBackendServiceStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateBackendServiceStep() (*CreateBackendServiceStep, error) {
	return &CreateBackendServiceStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertBackendService: func(ctx context.Context, config steps.GCEConfig, service *compute.BackendService) (*compute.Operation, error) {
					return client.RegionBackendServices.Insert(config.ServiceAccount.ProjectID, config.Region, service).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instaceGroupName string) (*compute.InstanceGroup, error) {
					return client.InstanceGroups.Get(config.ProjectID, config.AvailabilityZone, instaceGroupName).Do()
				},
				getBackendService: func(ctx context.Context, config steps.GCEConfig, backenServiceName string) (*compute.BackendService, error) {
					return client.RegionBackendServices.Get(config.ProjectID, config.Region, backenServiceName).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateBackendServiceStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	if !config.IsBootstrap {
		logrus.Debugf("Skip step %s for bootstrap node %s", CreateBackendServiceStepName, config.Node.Name)
		return nil
	}
	logrus.Debugf("Step %s", CreateBackendServiceStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateBackendServiceStepName)
	}

	backends := make([]*compute.Backend, 0)

	for _, instanceGroup := range config.GCEConfig.InstanceGroupLinks {
		backends = append(backends, &compute.Backend{
			Group: instanceGroup,
		})
	}

	backendService := &compute.BackendService{
		Name:                fmt.Sprintf("bs-%s", config.ClusterID),
		Description:         "Backend service for internal traffic",
		LoadBalancingScheme: "INTERNAL",
		Protocol:            "TCP",
		Region:              config.GCEConfig.Region,
		Backends:            backends,
		HealthChecks:        []string{config.GCEConfig.HealthCheckName},
	}

	_, err = svc.insertBackendService(ctx, config.GCEConfig, backendService)

	if err != nil {
		logrus.Errorf("error creating backend service %v", err)
		return errors.Wrapf(err, "error creating backend service")
	}

	backendService, err = svc.getBackendService(ctx, config.GCEConfig, backendService.Name)

	if err != nil {
		logrus.Errorf("error getting backend service %v", err)
		return errors.Wrapf(err, "error getting backend service")
	}

	config.GCEConfig.BackendServiceName = backendService.Name
	config.GCEConfig.BackendServiceLink = backendService.SelfLink

	logrus.Debugf("Created backend service name %s link %s",
		config.GCEConfig.BackendServiceName,
		config.GCEConfig.BackendServiceLink)
	// NOTE(stgleb): There is no field that signals for backend service is available so we simple sleep.
	time.Sleep(time.Minute * 1)
	return nil
}

func (s *CreateBackendServiceStep) Name() string {
	return CreateBackendServiceStepName
}

func (s *CreateBackendServiceStep) Depends() []string {
	return []string{CreateInstanceGroupsStepName, CreateNetworksStepName}
}

func (s *CreateBackendServiceStep) Description() string {
	return "Create backend service"
}

func (s *CreateBackendServiceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
