package gce

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
	"io"
)

const CreateInstanceGroupsStepName = "gce_create_instance_group"

type CreateInstanceGroupsStep struct {
	accountGetter accountGetter

	getComputeSvc     func(context.Context, steps.GCEConfig) (*computeService, error)
	zoneGetterFactory func(context.Context, accountGetter, *steps.Config) (account.ZonesGetter, error)
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
					return client.InstanceGroups.AddInstances(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, instanceGroupName, request).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.InstanceGroup, error) {
					return client.InstanceGroups.Get(config.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
				},
				getNetwork: func(ctx context.Context, config steps.GCEConfig, networkName string) (*compute.Network, error) {
					return client.Networks.Get(config.ProjectID, networkName).Do()
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
	if !config.IsMaster {
		logrus.Debugf("Skip step %s for non-master node %s", CreateInstanceGroupsStepName, config.Node.Name)
		return nil
	}

	logrus.Debugf("Step %s", CreateInstanceGroupsStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateInstanceGroupsStepName)
	}

	// Use naming convention az name + cluster so we do not need to store names of subnets for each subnet
	instanceGroup := &compute.InstanceGroup{
		Name:        fmt.Sprintf("%s-%s", config.GCEConfig.AvailabilityZone, config.ClusterID),
		Description: "Instance group for master nodes",
		Network:     config.GCEConfig.NetworkLink,
	}

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

	config.GCEConfig.InstanceGroupLinks[config.GCEConfig.AvailabilityZone] = instanceGroup.SelfLink
	config.GCEConfig.InstanceGroupNames[config.GCEConfig.AvailabilityZone] = instanceGroup.Name

	logrus.Debugf("Created instance group for az %s name %s link %s",
		config.GCEConfig.AvailabilityZone, config.GCEConfig.InstanceGroupNames[config.GCEConfig.AvailabilityZone],
		config.GCEConfig.InstanceGroupLinks[config.GCEConfig.AvailabilityZone])

	req := &compute.InstanceGroupsAddInstancesRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: config.Node.SelfLink,
			},
		},
	}

	logrus.Debugf("Add instance %s to instance group %s", config.Node.Name,
		config.GCEConfig.InstanceGroupNames[config.GCEConfig.AvailabilityZone])
	_, err = svc.addInstanceToInstanceGroup(ctx, config.GCEConfig, config.GCEConfig.AvailabilityZone, req)

	if err != nil {
		logrus.Errorf("error adding instance %s URL %s to instance group %s",
			config.Node.Name, config.Node.SelfLink, config.GCEConfig.InstanceGroupLinks[config.GCEConfig.AvailabilityZone])
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
