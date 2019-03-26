package gce

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateInstanceGroupStepName = "gce_create_instance_group"

type CreateInstanceGroupStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateInstanceGroupStep() (*CreateInstanceGroupStep, error) {
	return &CreateInstanceGroupStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertInstanceGroup: func(ctx context.Context, config steps.GCEConfig, group *compute.InstanceGroup) (*compute.Operation, error) {
					// TODO(stgleb): Create instance group for each AZ
					config.AvailabilityZone = "us-central1-a"
					return client.InstanceGroups.Insert(config.ServiceAccount.ProjectID, config.AvailabilityZone, group).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateInstanceGroupStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateInstanceGroupStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateInstanceGroupStepName)
	}

	instanceGroup := &compute.InstanceGroup{
		Name:        fmt.Sprintf("masters-%s", config.ClusterID),
		Description: "Instance group for master nodes",
	}

	_, err = svc.insertInstanceGroup(ctx, config.GCEConfig, instanceGroup)

	if err != nil {
		logrus.Errorf("Error creating instance group %v", err)
		return errors.Wrapf(err, "%s creating instance group caused", CreateInstanceGroupStepName)
	}

	config.GCEConfig.InstanceGroup = instanceGroup.Name

	return nil
}

func (s *CreateInstanceGroupStep) Name() string {
	return CreateInstanceGroupStepName
}

func (s *CreateInstanceGroupStep) Depends() []string {
	return nil
}

func (s *CreateInstanceGroupStep) Description() string {
	return "Create instance group for master nodes"
}

func (s *CreateInstanceGroupStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
