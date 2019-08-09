package gce

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateInstanceGroupsStepName = "gce_create_instance_group"

type CreateInstanceGroupsStep struct {
	accountGetter accountGetter

	getComputeSvc     func(context.Context, steps.GCEConfig) (*computeService, error)
	zoneGetterFactory func(context.Context, accountGetter, *steps.Config) (account.ZonesGetter, error)
}

func NewCreateInstanceGroupsStep(getter accountGetter) (*CreateInstanceGroupsStep, error) {
	return &CreateInstanceGroupsStep{
		accountGetter: getter,
		zoneGetterFactory: func(ctx context.Context, accountGetter accountGetter,
			cfg *steps.Config) (account.ZonesGetter, error) {
			acc, err := accountGetter.Get(ctx, cfg.CloudAccountName)

			if err != nil {
				logrus.Errorf("Get cloud account %s caused error %v",
					cfg.CloudAccountName, err)
				return nil, errors.Wrapf(err, "Get cloud account")
			}

			zoneGetter, err := account.NewZonesGetter(acc, cfg)

			return zoneGetter, err
		},
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertInstanceGroup: func(ctx context.Context, config steps.GCEConfig, group *compute.InstanceGroup) (*compute.Operation, error) {
					return client.InstanceGroups.Insert(config.ServiceAccount.ProjectID, config.AvailabilityZone, group).Do()
				},
				getInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroupName string) (*compute.InstanceGroup, error) {
					return client.InstanceGroups.Get(config.ProjectID, config.AvailabilityZone, instanceGroupName).Do()
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

	zoneGetter, err := s.zoneGetterFactory(ctx, s.accountGetter, config)

	if err != nil {
		logrus.Errorf("Create zone getter caused error %v", err)
		return errors.Wrap(err, "Create zone getter caused")
	}

	azs, err := zoneGetter.GetZones(ctx, *config)

	if err != nil {
		logrus.Errorf("get availability zones %v", err)
		return errors.Wrap(err, "get availability zones")
	}

	config.GCEConfig.AZs = make(map[string]string)

	for _, az := range azs {
		config.GCEConfig.AZs[az] = "dummy"
	}

	for az := range config.GCEConfig.AZs {
		// Use naming convention az name + cluster so we do not need to store names of subnets for each subnet
		instanceGroup := &compute.InstanceGroup{
			Name:        fmt.Sprintf("%s-%s", az, config.Kube.ID),
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

		logrus.Debugf("Created instance group for az %s name %s link %s and network %s",
			az, config.GCEConfig.InstanceGroupNames[az],
			config.GCEConfig.InstanceGroupLinks[az],
			instanceGroup.Network)
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
