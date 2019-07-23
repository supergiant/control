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

const CreateHealthCheckStepName = "gce_create_health_check"

type CreateHealthCheck struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateHealthCheckStep() *CreateHealthCheck {
	return &CreateHealthCheck{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertHealthCheck: func(ctx context.Context, config steps.GCEConfig, check *compute.HealthCheck) (*compute.Operation, error) {
					return client.HealthChecks.Insert(config.ServiceAccount.ProjectID, check).Do()
				},
				getHealthCheck: func(ctx context.Context, config steps.GCEConfig, healthCheckName string) (*compute.HealthCheck, error) {
					return client.HealthChecks.Get(config.ProjectID, healthCheckName).Do()
				},
			}, nil
		},
	}
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
		Name:               fmt.Sprintf("hc-%s", config.Kube.ID),
		CheckIntervalSec:   10,
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
		Type:               "HTTPS",
		HttpsHealthCheck: &compute.HTTPSHealthCheck{
			Port:        config.Kube.APIServerPort,
			RequestPath: "/healthz",
		},
	}

	_, err = svc.insertHealthCheck(ctx, config.GCEConfig, healthCheck)

	if err != nil {
		logrus.Errorf("Error creating external health check %v", err)
		return errors.Wrapf(err, "%s creating external health check caused",
			CreateHealthCheckStepName)
	}

	hc, err := svc.getHealthCheck(ctx, config.GCEConfig, healthCheck.Name)

	if err != nil {
		return errors.Wrapf(err, "Error creating health check")
	}

	logrus.Debugf("Created health check link %s", hc.SelfLink)
	config.GCEConfig.HealthCheckName = hc.SelfLink
	time.Sleep(time.Minute * 1)

	return nil
}

func (s *CreateHealthCheck) Name() string {
	return CreateHealthCheckStepName
}

func (s *CreateHealthCheck) Depends() []string {
	return []string{CreateTargetPullStepName}
}

func (s *CreateHealthCheck) Description() string {
	return "Create health checks"
}

func (s *CreateHealthCheck) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
