package amazon

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	StepNameCreateEC2Instance = "aws_create_instance"
)

type instanceService interface {
	RunInstancesWithContext(aws.Context, *ec2.RunInstancesInput, ...request.Option) (*ec2.Reservation, error)
	DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error)
	WaitUntilInstanceRunningWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.WaiterOption) error
}

type StepCreateInstance struct {
	getSvc func(steps.AWSConfig) (instanceService, error)
}

//InitCreateMachine adds the step to the registry
func InitCreateMachine(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepNameCreateEC2Instance, NewCreateInstance(ec2fn))
}

func NewCreateInstance(ec2fn GetEC2Fn) *StepCreateInstance {
	return &StepCreateInstance{
		getSvc: func(config steps.AWSConfig) (instanceService, error) {
			EC2, err := ec2fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *StepCreateInstance) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	// TODO: reuse sessions
	ec2Svc, err := s.getSvc(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	role := model.RoleMaster
	if !cfg.IsMaster {
		role = model.RoleNode
	}

	nodeName := util.MakeNodeName(cfg.ClusterName, cfg.TaskID, cfg.IsMaster)

	cfg.Node = model.Machine{
		Name:     nodeName,
		TaskID:   cfg.TaskID,
		Region:   cfg.AWSConfig.Region,
		Role:     role,
		Size:     cfg.AWSConfig.InstanceType,
		Provider: clouds.AWS,
		State:    model.MachineStatePlanned,
	}

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	var secGroupID *string
	var instanceProfileName *string

	//Determining a sec group in AWS for EC2 instance to be spawned.
	if cfg.IsMaster {
		secGroupID = &cfg.AWSConfig.MastersSecurityGroupID
		instanceProfileName = &cfg.AWSConfig.MastersInstanceProfile
	} else {
		secGroupID = &cfg.AWSConfig.NodesSecurityGroupID
		instanceProfileName = &cfg.AWSConfig.NodesInstanceProfile
	}

	isEbs := false
	volumeSize, err := strconv.Atoi(cfg.AWSConfig.VolumeSize)

	runInstanceInput := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(int64(volumeSize)),
				},
			},
		},
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(cfg.AWSConfig.AvailabilityZone),
		},
		EbsOptimized: &isEbs,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: instanceProfileName,
		},
		ImageId:      &cfg.AWSConfig.ImageID,
		InstanceType: &cfg.AWSConfig.InstanceType,
		KeyName:      &cfg.AWSConfig.KeyPairName,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),

		//TODO add custom TAGS
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("KubernetesCluster"),
						Value: aws.String(cfg.ClusterName),
					},
					{
						Key:   aws.String("Name"),
						Value: aws.String(nodeName),
					},
					{
						Key:   aws.String("Role"),
						Value: aws.String(util.MakeRole(cfg.IsMaster)),
					},
					{
						Key:   aws.String(clouds.TagClusterID),
						Value: aws.String(cfg.ClusterID),
					},
				},
			},
		},
	}
	if cfg.AWSConfig.HasPublicAddr {
		runInstanceInput.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				AssociatePublicIpAddress: aws.Bool(true),
				DeleteOnTermination:      aws.Bool(true),
				SubnetId:                 aws.String(cfg.AWSConfig.Subnets[cfg.AWSConfig.AvailabilityZone]),
				Groups:                   []*string{secGroupID},
			},
		}
	}

	res, err := ec2Svc.RunInstancesWithContext(ctx, runInstanceInput)
	if err != nil {
		cfg.Node.State = model.MachineStateError
		cfg.NodeChan() <- cfg.Node

		log.Errorf("[%s] - failed to create ec2 instance: %v", StepNameCreateEC2Instance, err)
		return errors.Wrap(ErrCreateInstance, err.Error())
	}

	cfg.Node = model.Machine{
		Name:     nodeName,
		TaskID:   cfg.TaskID,
		Region:   cfg.AWSConfig.Region,
		Role:     role,
		Provider: clouds.AWS,
		Size:     cfg.AWSConfig.InstanceType,
		State:    model.MachineStateBuilding,
	}

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	if len(res.Instances) == 0 {
		cfg.Node.State = model.MachineStateError
		cfg.NodeChan() <- cfg.Node

		return errors.Wrap(ErrCreateInstance, "no instances created")
	}

	instance := res.Instances[0]

	if cfg.AWSConfig.HasPublicAddr {
		log.Infof("[%s] - waiting to obtain public IP...", s.Name())

		lookup := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String(nodeName)},
				},
				{
					Name:   aws.String(fmt.Sprintf("tag:%s", clouds.TagClusterID)),
					Values: []*string{aws.String(cfg.ClusterID)},
				},
			},
		}
		logrus.Debugf("Wait until instance %s running", nodeName)
		err = ec2Svc.WaitUntilInstanceRunningWithContext(ctx, lookup)

		if err != nil {
			logrus.Errorf("Error waiting instance %s cluster %s running %v",
				nodeName, cfg.ClusterID, err)
			return errors.Wrapf(err, "Error waiting instance %s cluster-id %s",
				nodeName, cfg.ClusterID)
		}

		logrus.Debugf("Instance running %s", nodeName)

		out, err := ec2Svc.DescribeInstancesWithContext(ctx, lookup)
		if err != nil {
			cfg.Node.State = model.MachineStateError
			cfg.NodeChan() <- cfg.Node
			log.Errorf("[%s] - failed to obtain public IP for node %s: %v", s.Name(), nodeName, err)
			return errors.Wrap(ErrNoPublicIP, err.Error())
		}

		if i := findInstanceWithPublicAddr(out.Reservations); i != nil {
			cfg.Node.PublicIp = *i.PublicIpAddress
			cfg.Node.PrivateIp = *i.PrivateIpAddress
			log.Infof("[%s] - found public ip - %s for node %s", s.Name(), cfg.Node.PublicIp, nodeName)
		} else {
			log.Errorf("[%s] - failed to find public IP address", s.Name())
			cfg.Node.State = model.MachineStateError
			cfg.NodeChan() <- cfg.Node
			return ErrNoPublicIP
		}
	}

	cfg.Node.Region = cfg.AWSConfig.Region
	cfg.Node.CreatedAt = instance.LaunchTime.Unix()
	cfg.Node.ID = *instance.InstanceId
	cfg.Node.State = model.MachineStateProvisioning

	logrus.Infof("Machine created %v", cfg.Node)
	cfg.NodeChan() <- cfg.Node
	if cfg.IsMaster {
		cfg.AddMaster(&cfg.Node)
	} else {
		cfg.AddNode(&cfg.Node)
	}

	log.Infof("[%s] - success! Created node %s with instanceID %s ",
		s.Name(), nodeName, cfg.Node.ID)

	return nil
}

func (s *StepCreateInstance) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func findInstanceWithPublicAddr(reservations []*ec2.Reservation) *ec2.Instance {
	for _, r := range reservations {
		for _, i := range r.Instances {
			if i.PublicIpAddress != nil {
				return i
			}
		}
	}
	return nil
}

func (*StepCreateInstance) Name() string {
	return StepNameCreateEC2Instance
}

func (*StepCreateInstance) Description() string {
	return "Create EC2 Instance"
}

func (*StepCreateInstance) Depends() []string {
	return nil
}
