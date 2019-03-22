package gce

import (
	"context"
	"io"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"


	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/pkg/errors"
	"fmt"
)

const CreateTargetPullStepName = "gce_create_target_pool"

type CreateTargetPoolStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateTargetPoolStep() (*CreateInstanceStep, error) {
	return &CreateInstanceStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config.ClientEmail,
				config.PrivateKey, config.TokenURI)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertAddress: func(ctx context.Context, config steps.GCEConfig, address *compute.Address) (*compute.Operation, error) {
					return client.Addresses.Insert(config.ProjectID, config.Region, address).Do()
				},
				insertForwardingRule: func(ctx context.Context, config steps.GCEConfig, forwardingRule *compute.ForwardingRule) (*compute.Operation, error) {
					return client.ForwardingRules.Insert(config.ProjectID, config.Region, forwardingRule).Do()
				},
				insertTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPool *compute.TargetPool) (*compute.Operation, error) {
					return client.TargetPools.Insert(config.ProjectID, config.Region, targetPool).Do()
				},
				addInstanceToTargetGroup: func(ctx context.Context, config steps.GCEConfig, targetPoolName string, req *compute.TargetPoolsAddInstanceRequest) (*compute.Operation, error) {
					return client.TargetPools.AddInstance(config.ProjectID, config.Region, targetPoolName, req).Do()
				},
				insertHealthCheck: func(ctx context.Context, config steps.GCEConfig, check *compute.HealthCheck) (*compute.Operation, error) {
					return client.HealthChecks.Insert(config.ProjectID, check).Do()
				},
				addHealthCheckToTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPool string, request *compute.TargetPoolsAddHealthCheckRequest) (*compute.Operation, error) {
					return client.TargetPools.AddHealthCheck(config.ProjectID, config.Region, targetPool, request).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateTargetPoolStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateTargetPullStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateTargetPullStepName)
	}

	externalTargetPoolName := fmt.Sprintf("ex-%s", config.ClusterID)
	externalTargetPool := &compute.TargetPool{
		Description: "Target pool for balancing external traffic",
		Name: externalTargetPoolName,
	}

	_, err = svc.insertTargetPool(ctx, config.GCEConfig, externalTargetPool)

	if err != nil {
		logrus.Errorf("Error creating external target pool %v", err)
		return errors.Wrapf(err, "%s creating external target pool", CreateTargetPullStepName)
	}

	config.GCEConfig.ExternalTargetPoolName = externalTargetPoolName

	internalTargetPoolName := fmt.Sprintf("ex-%s", config.ClusterID)
	internalTargetPool := &compute.TargetPool{
		Description: "Target pool for balancing internal traffic",
		Name: internalTargetPoolName,
	}

	_, err = svc.insertTargetPool(ctx, config.GCEConfig, internalTargetPool)

	if err != nil {
		logrus.Errorf("Error creating internal target pool %v", err)
		return errors.Wrapf(err, "%s creating internal target pool", CreateTargetPullStepName)
	}

	config.GCEConfig.InternalTargetPoolName = internalTargetPoolName

	return nil
}

func (s *CreateTargetPoolStep) Name() string {
	return CreateTargetPullStepName
}

func (s *CreateTargetPoolStep) Depends() []string {
	return nil
}

func (s *CreateTargetPoolStep) Description() string {
	return "Create target pool"
}

func (s *CreateTargetPoolStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
