package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
)

const (
	ImportClusterMachinesStepName = "import_cluster_machines_aws"
	running                       = "running"
)

type InstanceDescriber interface {
	DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error)
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
}

type ImportClusterStep struct {
	getSvc func(config steps.AWSConfig) (InstanceDescriber, error)
}

func NewImportClusterStep(fn GetEC2Fn) *ImportClusterStep {
	return &ImportClusterStep{
		getSvc: func(config steps.AWSConfig) (describer InstanceDescriber, e error) {
			EC2, err := fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func InitImportClusterStep(fn GetEC2Fn) {
	steps.RegisterStep(ImportClusterMachinesStepName, NewImportClusterStep(fn))
}

func (s ImportClusterStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	logrus.Info(ImportClusterMachinesStepName)
	ec2Svc, err := s.getSvc(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	err = s.importMachines(ctx, model.RoleMaster, ec2Svc, cfg.GetMasters(), cfg)

	if err != nil {
		return errors.Wrapf(err, "error importing masters")
	}

	if master := cfg.GetMaster(); master != nil {
		cfg.Node = *master
	}

	err = s.importMachines(ctx, model.RoleNode, ec2Svc, cfg.GetNodes(), cfg)

	if err != nil {
		return errors.Wrapf(err, "error importing workers")
	}

	secGroupInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(cfg.AWSConfig.MastersSecurityGroupID)},
	}

	output, err := ec2Svc.DescribeSecurityGroups(secGroupInput)

	if err != nil {
		return errors.Wrapf(err, "master security group not found %s", cfg.AWSConfig.MastersSecurityGroupID)
	}

	masterSecGroup := output.SecurityGroups[0]
	cfg.AWSConfig.VPCID = *masterSecGroup.VpcId

	return nil
}

func (s ImportClusterStep) Name() string {
	return ImportClusterMachinesStepName
}

func (s ImportClusterStep) Description() string {
	return ImportClusterMachinesStepName
}

func (s ImportClusterStep) Depends() []string {
	return nil
}

func (s ImportClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func findInstanceWithPrivateIPAddr(reservations []*ec2.Reservation) *ec2.Instance {
	for _, r := range reservations {
		for _, i := range r.Instances {
			if i.PrivateIpAddress != nil {
				return i
			}
		}
	}
	return nil
}

func (s *ImportClusterStep) importMachines(ctx context.Context, role model.Role, ec2Svc InstanceDescriber, machines map[string]*model.Machine, cfg *steps.Config) error {
	for _, machine := range machines {
		describeReq := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("network-interface.addresses.private-ip-address"),
					Values: []*string{aws.String(machine.PrivateIp)},
				},
			},
		}

		output, err := ec2Svc.DescribeInstancesWithContext(ctx, describeReq)

		if err != nil {
			return errors.Wrapf(err, "error importing node %v", machine)
		}

		if len(output.Reservations) == 0 {
			return errors.Wrapf(sgerrors.ErrNotFound, "instance %v not found", machine)
		}


		instance := findInstanceWithPrivateIPAddr(output.Reservations)
		machine.ID = *instance.InstanceId
		machine.Size = *instance.InstanceType
		machine.CreatedAt = instance.LaunchTime.Unix()
		machine.AvailabilityZone = *instance.Placement.AvailabilityZone
		machine.Region = azToRegion(*instance.Placement.AvailabilityZone)
		machine.Provider = cfg.Provider
		machine.PrivateIp = *instance.PrivateIpAddress
		machine.PublicIp = *instance.PublicIpAddress
		machine.State = instanceStateToMachineState(*instance.State.Name)

		cfg.AWSConfig.ImageID = *instance.ImageId
		cfg.AWSConfig.Region = azToRegion(*instance.Placement.AvailabilityZone)
		cfg.AWSConfig.KeyPairName = *instance.KeyName

		if len(instance.SecurityGroups) == 0 {
			return errors.Wrapf(sgerrors.ErrNotFound, "no security groups found for %v", machine)
		}

		if role == model.RoleMaster {
			cfg.AWSConfig.MastersSecurityGroupID = *instance.SecurityGroups[0].GroupId
			machine.Role = roleMaster
			cfg.AddMaster(machine)
		} else {
			cfg.AWSConfig.NodesSecurityGroupID = *instance.SecurityGroups[0].GroupId
			machine.Role = model.RoleNode
			cfg.AddNode(machine)
		}

		tags := &ec2.CreateTagsInput{
			Resources: []*string{instance.InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("KubernetesCluster"),
					Value: aws.String(cfg.ClusterName),
				},
				{
					Key:   aws.String("Role"),
					Value: aws.String(string(role)),
				},
				{
					Key:   aws.String(clouds.TagClusterID),
					Value: aws.String(cfg.ClusterID),
				},
			},
		}

		_, err = ec2Svc.CreateTags(tags)

		if err != nil {
			return errors.Wrapf(err, "error tag instance %v", machine)
		}
	}

	return nil
}

func azToRegion(az string) string {
	if len(az) == 0 {
		return az
	}

	return az[:len(az)-1]
}

func instanceStateToMachineState(instanceState string) model.MachineState {
	if instanceState == running {
		return model.MachineStateActive
	}

	return model.MachineStateError
}
