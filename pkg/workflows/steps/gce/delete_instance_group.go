package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteInstanceGroupStepName = "gce_delete_instance_group"

type DeleteInstanceGroupStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteInstanceGroupStep() (*DeleteInstanceGroupStep, error) {
	return &DeleteInstanceGroupStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.Operation, error) {
					config.AvailabilityZone = "us-central1-a"
					return client.InstanceGroups.Delete(config.ServiceAccount.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
				},
			}, nil
		},
	}, nil
}

func (s *DeleteInstanceGroupStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteForwardingRulesStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateTargetPullStepName)
	}

	_, err = svc.deleteInstanceGroup(ctx, config.GCEConfig, config.GCEConfig.InstanceGroupName)

	if err != nil {
		logrus.Errorf("Error deleting instance group %v", err)
		return errors.Wrapf(err, "%s deleting instance group  caused", CreateTargetPullStepName)
	}

	return nil
}

func (s *DeleteInstanceGroupStep) Name() string {
	return DeleteInstanceGroupStepName
}

func (s *DeleteInstanceGroupStep) Depends() []string {
	return nil
}

func (s *DeleteInstanceGroupStep) Description() string {
	return "Delete instance group for master nodes"
}

func (s *DeleteInstanceGroupStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
