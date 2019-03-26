package gce

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateHealthCheckStepName = "gce_create_health_check"

type CreateHealthCheck struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateHealthCheckStep() (*CreateAddressStep, error) {
	return &CreateAddressStep{
		Timeout:      time.Second * 10,
		AttemptCount: 6,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertHealthCheck: func(ctx context.Context, config steps.GCEConfig, check *compute.HealthCheck) (*compute.Operation, error) {
					return client.HealthChecks.Insert(config.ServiceAccount.ProjectID, check).Do()
				},
				addHealthCheckToTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPool string, request *compute.TargetPoolsAddHealthCheckRequest) (*compute.Operation, error) {
					return client.TargetPools.AddHealthCheck(config.ServiceAccount.ProjectID, config.Region, targetPool, request).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateHealthCheck) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateHealthCheckStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateHealthCheckStepName)
	}

	healthCheck := &compute.HealthCheck{
		Name:               "healthCheck",
		CheckIntervalSec:   10,
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
		Type:               "HTTPS",
		HttpsHealthCheck: &compute.HTTPSHealthCheck{
			Port:        443,
			RequestPath: "/healthz",
		},
	}

	_, err = svc.insertHealthCheck(ctx, config.GCEConfig, healthCheck)

	if err != nil {
		logrus.Errorf("Error creating external health check %v", err)
		return errors.Wrapf(err, "%s creating external health check caused",
			CreateHealthCheckStepName)
	}

	addHealthCheckRequest := &compute.TargetPoolsAddHealthCheckRequest{
		HealthChecks: []*compute.HealthCheckReference{
			{
				HealthCheck: healthCheck.Name,
			},
		},
	}

	_, err = svc.addHealthCheckToTargetPool(ctx, config.GCEConfig,
		config.GCEConfig.TargetPoolName, addHealthCheckRequest)

	if err != nil {
		logrus.Errorf("Error adding health check %s to target pool %s %v",
			healthCheck.Name, config.GCEConfig.TargetPoolName, err)
		return errors.Wrapf(err, "%s adding health check %s to target pool %s caused",
			healthCheck.Name, config.GCEConfig.TargetPoolName, CreateHealthCheckStepName)
	}

	_, err = svc.addHealthCheckToTargetPool(ctx, config.GCEConfig,
		config.GCEConfig.BackendServiceName, addHealthCheckRequest)

	if err != nil {
		logrus.Errorf("Error adding health check %s to target pool %s %v",
			healthCheck.Name, config.GCEConfig.BackendServiceName, err)
		return errors.Wrapf(err, "%s adding health check %s to target pool %s caused",
			healthCheck.Name, config.GCEConfig.BackendServiceName, CreateHealthCheckStepName)
	}

	return nil
}

func (s *CreateHealthCheck) Name() string {
	return CreateHealthCheckStepName
}

func (s *CreateHealthCheck) Depends() []string {
	return []string{CreateTargetPullStepName, CreateBackendServiceStepName}
}

func (s *CreateHealthCheck) Description() string {
	return "Create health checks"
}

func (s *CreateHealthCheck) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
