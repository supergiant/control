package amazon

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepNameCreateEC2Instance = "aws_create_instance"
const (
	IPAttempts             = 12
	SleepSecondsPerAttempt = 6
)

type StepCreateInstance struct {
}

//InitStepCreateInstance adds the step to the registry
func InitStepCreateInstance() {
	steps.RegisterStep(StepNameCreateEC2Instance, NewCreateInstance())
}

func NewCreateInstance() *StepCreateInstance {
	return &StepCreateInstance{}
}

func (s *StepCreateInstance) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - started", StepNameCreateEC2Instance)

	ec2Cfg := cfg.AWSConfig.EC2Config

	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	//If subnetID is nil, the default would be used
	var subnetID *string
	if ec2Cfg.SubnetID != "" {
		subnetID = &ec2Cfg.SubnetID
	}

	nodeName := util.MakeNodeName(cfg.ClusterName, cfg.TaskID, cfg.IsMaster)
	runInstanceInput := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(int64(ec2Cfg.VolumeSize)),
				},
			},
		},
		EbsOptimized: &ec2Cfg.EbsOptimized,
		ImageId:      &ec2Cfg.ImageID,
		InstanceType: &ec2Cfg.InstanceType,
		KeyName:      &cfg.AWSConfig.KeyPairName,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		//PrivateIpAddress:        nil,
		//TODO security groups
		SecurityGroupIds: nil,
		SecurityGroups:   nil,
		SubnetId:         subnetID,

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
				},
			},
		},
	}
	if ec2Cfg.HasPublicAddr {
		runInstanceInput.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				AssociatePublicIpAddress: aws.Bool(ec2Cfg.HasPublicAddr),
				DeleteOnTermination:      aws.Bool(true),

				//TODO security groups
				SubnetId: subnetID,
			},
		}
	}

	if cfg.AWSConfig.EC2Config.GPU {
		//TODO ADD GPU SUPPORT FOR AWS
	}

	role := node.RoleMaster
	if !cfg.IsMaster {
		role = node.RoleNode
	}

	cfg.Node = node.Node{
		TaskID: cfg.TaskID,
		Region:   cfg.AWSConfig.Region,
		Role:     role,
		Provider: clouds.AWS,
		State:    node.StateBuilding,
	}

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	res, err := sdk.EC2.RunInstances(runInstanceInput)
	if err != nil {
		cfg.Node.State = node.StateError
		cfg.NodeChan() <- cfg.Node

		log.Errorf("[%s] - failed to create ec2 instance: %v", StepNameCreateEC2Instance, err)
		return errors.Wrap(err, "aws: failed to connect")
	}
	if len(res.Instances) == 0 {
		cfg.Node.State = node.StateError
		cfg.NodeChan() <- cfg.Node

		return errors.Wrap(err, "aws: no instances created")
	}
	instance := res.Instances[0]

	cfg.Node.Region = cfg.AWSConfig.Region
	cfg.Node.CreatedAt = instance.LaunchTime.Unix()
	cfg.Node.ID = *instance.InstanceId

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	if ec2Cfg.HasPublicAddr {
		log.Infof("[%s] - waiting to obtain public IP", StepNameCreateEC2Instance)

		found := false
		//Waiting for AWS to assign public IP requires us to poll an describe ec2 endpoint several times
		for i := 0; i < IPAttempts; i++ {
			lookup := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("tag:Name"),
						Values: []*string{aws.String(nodeName)},
					},
					{
						Name:   aws.String("tag:KubernetesCluster"),
						Values: []*string{aws.String(cfg.ClusterName)},
					},
				},
			}
			out, err := sdk.EC2.DescribeInstancesWithContext(ctx, lookup)
			if err != nil {
				cfg.Node.State = node.StateError
				cfg.NodeChan() <- cfg.Node
				log.Errorf("[%s] - failed to obtain public IP for node %s: %v", StepNameCreateEC2Instance, nodeName, err)
				return errors.Wrap(err, "aws: failed to obtain public IP")
			}

			if len(out.Reservations) == 0 {
				log.Infof("[%s] - found 0 ec2 instances, attempt %d", StepNameCreateEC2Instance, i)
				time.Sleep(time.Duration(SleepSecondsPerAttempt) * time.Second)
				continue
			}

			if i := findInstanceWithPublicAddr(out.Reservations); i != nil {
				cfg.Node.PublicIp = *i.PublicIpAddress
				cfg.Node.PrivateIp = *i.PrivateIpAddress
				log.Infof("[%s] - found public ip - %s for node %s", StepNameCreateEC2Instance, cfg.Node.PublicIp, nodeName)
				found = true
				break
			}
		}
		if !found {
			cfg.Node.State = node.StateError
			cfg.NodeChan() <- cfg.Node
			log.Errorf("[%s] - failed to find public IP address after %d attempts", StepNameCreateEC2Instance,
				IPAttempts)
			return errors.New("aws: failed to obtain public IP")
		}
	}

	if cfg.IsMaster {
		cfg.AddMaster(&cfg.Node)
	}
	cfg.Node.State = node.StateProvisioning
	cfg.NodeChan() <- cfg.Node

	log.Infof("[%s] - success! Created node %s with instanceID %s",
		StepNameCreateEC2Instance, nodeName, cfg.Node.ID)
	logrus.Debugf("%v", *instance)

	return nil
}

func (s *StepCreateInstance) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - rollback initiated", s.Name())

	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	if cfg.Node.ID != "" {
		_, err := sdk.EC2.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(cfg.Node.ID),
			},
		})
		if err != nil {
			return err
		}
		log.Infof("[%s] - deleted ec2 instance %s", s.Name(), cfg.Node.ID)
		return nil
	}
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
