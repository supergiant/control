package gce

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"google.golang.org/api/compute/v1"
	"io"

	"fmt"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateTargetPullStepName = "gce_create_target_pool"

type CreateTargetPoolStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateTargetPoolStep() *CreateTargetPoolStep {
	return &CreateTargetPoolStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPool *compute.TargetPool) (*compute.Operation, error) {
					return client.TargetPools.Insert(config.ServiceAccount.ProjectID, config.Region, targetPool).Do()
				},
				getTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPoolName string) (*compute.TargetPool, error) {
					return client.TargetPools.Get(config.ProjectID, config.Region, targetPoolName).Do()
				},
			}, nil
		},
	}
}

func (s *CreateTargetPoolStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateTargetPullStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateTargetPullStepName)
	}

	targetPoolName := fmt.Sprintf("ex-%s", config.ClusterID)
	targetPool := &compute.TargetPool{
		Description: "Target pool for balancing traffic",
		Name:        targetPoolName,
	}

	_, err = svc.insertTargetPool(ctx, config.GCEConfig, targetPool)

	if err != nil {
		logrus.Errorf("Error creating target pool %v", err)
		return errors.Wrapf(err, "%s creating target pool", CreateTargetPullStepName)
	}

	targetPool, err = svc.getTargetPool(ctx, config.GCEConfig, targetPoolName)

	if err != nil {
		logrus.Errorf("Error getting target pool %v", err)
		return errors.Wrapf(err, "%s getting target pool", CreateTargetPullStepName)
	}

	config.GCEConfig.TargetPoolName = targetPoolName
	config.GCEConfig.TargetPoolLink = targetPool.SelfLink
	logrus.Debugf("Created target pool name %s link %s",
		config.GCEConfig.TargetPoolName, config.GCEConfig.TargetPoolLink)
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
