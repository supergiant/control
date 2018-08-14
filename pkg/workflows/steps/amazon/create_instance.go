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
	"github.com/supergiant/supergiant/pkg/clouds/amazonSDK"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "aws_create_instance"
const (
	IPAttempts             = 12
	SleepSecondsPerAttempt = 6
)

type StepCreateInstance struct {
	//InstanceID could be used in Rollback to delete node
	instanceID string
	sdk        *amazonSDK.SDK
}

func (s StepCreateInstance) GetInstanceID() string {
	return s.instanceID
}

func NewCreateInstance(awsSDK *amazonSDK.SDK) *StepCreateInstance {
	return &StepCreateInstance{
		sdk: awsSDK,
	}
}

func (s *StepCreateInstance) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - started", StepName)

	ec2Cfg := cfg.AWSConfig.EC2Config

	//If subnetID is nil, the default would be used
	var subnetID *string
	if ec2Cfg.SubnetID != "" {
		subnetID = &ec2Cfg.SubnetID
	}

	nodeName := util.MakeNodeName(cfg.ClusterName, cfg.IsMaster)

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

	res, err := s.sdk.EC2.RunInstances(runInstanceInput)
	if err != nil {
		log.Errorf("[%s] - failed to create ec2 instance: %v", StepName, err)
		return errors.Wrap(err, "aws: failed to connect")
	}
	if len(res.Instances) == 0 {
		return errors.Wrap(err, "aws: no instances created")
	}
	instance := res.Instances[0]

	//saving instance ID for rollback
	s.instanceID = *instance.InstanceId

	n := &node.Node{
		Region:    cfg.AWSConfig.Region,
		CreatedAt: time.Now().Unix(),
		Provider:  clouds.AWS,
		Id:        *instance.InstanceId,
	}

	if ec2Cfg.HasPublicAddr {
		log.Infof("[%s] - waiting to obtain public IP", StepName)

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
			out, err := s.sdk.EC2.DescribeInstancesWithContext(ctx, lookup)
			if err != nil {
				log.Errorf("[%s] - failed to obtain public IP for node %s: %v", StepName, nodeName, err)
				return errors.Wrap(err, "aws: failed to obtain public IP")
			}

			if l := len(out.Reservations); l == 0 {
				log.Infof("[%s] - found 0 ec2 instances, attempt %d", StepName, i)
				time.Sleep(time.Duration(SleepSecondsPerAttempt) * time.Second)
				continue
			}

			for _, r := range out.Reservations {
				for _, i := range r.Instances {
					if i.PublicIpAddress != nil {
						n.PublicIp = *i.PublicIpAddress
						n.PrivateIp = *i.PrivateIpAddress

						log.Info("[%s] - found public ip - %s for node %s", StepName, n.PublicIp, nodeName)
						goto writeResult
					}
				}
			}
			log.Errorf("[%s] - failed to find public IP address after %d attempts", StepName, i)
			return errors.New("aws: failed to obtain public IP")
		}
	}

writeResult:
	if cfg.IsMaster {
		cfg.AddMaster(n)
	} else {
		cfg.Node = *n
	}
	log.Infof("[%s] - success! Created node %s with instanceID %s",
		StepName, nodeName, n.Id)
	logrus.Debugf("%v", *instance)

	return nil
}

func (s *StepCreateInstance) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	log.Infof("[%s] - rollback initiated", s.Name())
	if s.instanceID != "" {
		_, err := s.sdk.EC2.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(s.instanceID),
			},
		})
		if err != nil {
			return err
		}
		log.Infof("[%s] - deleted ec2 instance %s", s.Name(), s.instanceID)
		return nil
	}
	return nil
}

func (*StepCreateInstance) Name() string {
	return StepName
}

func (*StepCreateInstance) Description() string {
	return "Create EC2 Instance"
}

func (*StepCreateInstance) Depends() []string {
	return nil
}
