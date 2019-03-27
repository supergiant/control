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
				addInstanceToInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string, request *compute.InstanceGroupsAddInstancesRequest) (*compute.Operation, error) {
					// TODO(stgleb): Create instance group for each AZ
					config.AvailabilityZone = "us-central1-a"

					return client.InstanceGroups.AddInstances(config.ServiceAccount.ProjectID, config.AvailabilityZone, instanceGroupName, request).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.InstanceGroup, error) {
					// TODO(stgleb): Create instance group for each AZ
					config.AvailabilityZone = "us-central1-a"

					return client.InstanceGroups.Get(config.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
				},
				getInstance: func(ctx context.Context,
					config steps.GCEConfig, name string) (*compute.Instance, error) {
					return client.Instances.Get(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, name).Do()
				},
			}, nil
		},
	}, nil
}

func (s *CreateInstanceGroupStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	// Skip this step for the rest of nodes
	if !config.KubeadmConfig.IsBootstrap {
		return nil
	}

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

	instanceGroup, err = svc.getInstanceGroup(ctx, config.GCEConfig, instanceGroup.Name)

	if err != nil {
		logrus.Errorf("Error getting instance group %v", err)
		return errors.Wrapf(err, "%s creating getting group caused", CreateInstanceGroupStepName)
	}

	instance, err := svc.getInstance(ctx, config.GCEConfig, config.Node.Name)

	if err != nil {
		logrus.Errorf("getting instance caused %v", err)
		return errors.Wrapf(err, "%s getting instance",
			CreateInstanceStepName)
	}

	config.GCEConfig.InstanceGroupName = instanceGroup.Name
	config.GCEConfig.InstanceGroupLink = instanceGroup.SelfLink

	svc.addInstanceToInstanceGroup(ctx, config.GCEConfig, instanceGroup.Name, &compute.InstanceGroupsAddInstancesRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: instance.SelfLink,
			},
		},
	})

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
