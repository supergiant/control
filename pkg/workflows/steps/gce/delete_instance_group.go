package gce

import (
	"context"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteInstanceGroupStepName = "gce_delete_instance_group"

type DeleteInstanceGroupStep struct {
	Timeout       time.Duration
	AttemptCount  int
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteInstanceGroupStep() (*DeleteInstanceGroupStep, error) {
	return &DeleteInstanceGroupStep{
		Timeout:      time.Second * 10,
		AttemptCount: 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.Operation, error) {
					return client.InstanceGroups.Delete(config.ServiceAccount.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.InstanceGroup, error) {
					return client.InstanceGroups.Get(config.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
				},
			}, nil
		},
	}, nil
}

func (s *DeleteInstanceGroupStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteInstanceGroupStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteInstanceGroupStepName)
	}

	timeout := s.Timeout

	for az, instanceGroupName := range config.GCEConfig.InstanceGroupNames {
		config.GCEConfig.AvailabilityZone = az

		for i := 0; i < s.AttemptCount; i++ {
			_, err = svc.getInstanceGroup(ctx, config.GCEConfig, instanceGroupName)

			if isNotFound(err) {
				break
			}

			_, err = svc.deleteInstanceGroup(ctx, config.GCEConfig, instanceGroupName)

			if err != nil {
				logrus.Debugf("Error deleting instance group %s %v", instanceGroupName, err)
			} else {
				break
			}

			time.Sleep(timeout)
			timeout = timeout * 2
		}
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
