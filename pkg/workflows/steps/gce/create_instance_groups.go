package gce

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateInstanceGroupsStepName = "gce_create_instance_group"

type CreateInstanceGroupsStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateInstanceGroupsStep() (*CreateInstanceGroupsStep, error) {
	return &CreateInstanceGroupsStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertInstanceGroup: func(ctx context.Context, config steps.GCEConfig, group *compute.InstanceGroup) (*compute.Operation, error) {
					return client.InstanceGroups.Insert(config.ServiceAccount.ProjectID, config.AvailabilityZone, group).Do()
				},
				addInstanceToInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string, request *compute.InstanceGroupsAddInstancesRequest) (*compute.Operation, error) {
					return client.InstanceGroups.AddInstances(config.ServiceAccount.ProjectID, config.AvailabilityZone, instanceGroupName, request).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.InstanceGroup, error) {
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

func (s *CreateInstanceGroupsStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", CreateInstanceGroupsStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateInstanceGroupsStepName)
	}

	config.GCEConfig.InstanceGroupLinks = make(map[string]string)
	config.GCEConfig.InstanceGroupNames = make(map[string]string)

	for az := range config.GCEConfig.Subnets {
		// Use naming convention az name + cluster so we do not need to store names of subnets for each subnet
		instanceGroup := &compute.InstanceGroup{
			Name:        fmt.Sprintf("%s-%s", az, config.ClusterID),
			Description: "Instance group for master nodes",
			Network:     config.GCEConfig.NetworkLink,
		}

		config.GCEConfig.AvailabilityZone = az
		_, err = svc.insertInstanceGroup(ctx, config.GCEConfig, instanceGroup)

		if err != nil {
			logrus.Errorf("Error creating instance group %v", err)
			return errors.Wrapf(err, "%s creating instance group caused", CreateInstanceGroupsStepName)
		}

		instanceGroup, err = svc.getInstanceGroup(ctx, config.GCEConfig, instanceGroup.Name)

		if err != nil {
			logrus.Errorf("Error getting instance group %v", err)
			return errors.Wrapf(err, "%s creating getting group caused", CreateInstanceGroupsStepName)
		}

		config.GCEConfig.InstanceGroupLinks[az] = instanceGroup.SelfLink
		config.GCEConfig.InstanceGroupNames[az] = instanceGroup.Name
	}

	return nil
}

func (s *CreateInstanceGroupsStep) Name() string {
	return CreateInstanceGroupsStepName
}

func (s *CreateInstanceGroupsStep) Depends() []string {
	return nil
}

func (s *CreateInstanceGroupsStep) Description() string {
	return "Create instance group for master nodes"
}

func (s *CreateInstanceGroupsStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
