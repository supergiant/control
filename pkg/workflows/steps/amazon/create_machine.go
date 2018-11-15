package amazon

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const (
	StepNameCreateEC2Instance = "aws_create_instance"
	IPAttempts                = 10
	SleepSecondsPerAttempt    = 6
)

type StepCreateInstance struct {
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitCreateMachine(fn func(steps.AWSConfig) (ec2iface.EC2API, error)) {
	steps.RegisterStep(StepNameCreateEC2Instance, NewCreateInstance(fn))
}

func NewCreateInstance(fn GetEC2Fn) *StepCreateInstance {
	return &StepCreateInstance{
		GetEC2: fn,
	}
}

func (s *StepCreateInstance) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	role := node.RoleMaster
	if !cfg.IsMaster {
		role = node.RoleNode
	}

	nodeName := util.MakeNodeName(cfg.ClusterName, cfg.TaskID, cfg.IsMaster)

	cfg.Node = node.Node{
		Name:     nodeName,
		TaskID:   cfg.TaskID,
		Region:   cfg.AWSConfig.Region,
		Role:     role,
		Provider: clouds.AWS,
		State:    node.StatePlanned,
	}

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	var secGroupID *string

	//Determining a sec group in AWS for EC2 instance to be spawned.
	if cfg.IsMaster {
		secGroupID = &cfg.AWSConfig.MastersSecurityGroupID
	} else {
		secGroupID = &cfg.AWSConfig.NodesSecurityGroupID
	}

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	amiID, err := s.FindAMI(ctx, w, EC2)
	if err != nil {
		logrus.Errorf("[%s] - failed to find AMI for Ubuntu: %v", s.Name(), err)
		return errors.Wrap(err, "failed to find AMI")
	}

	isEbs, err := strconv.ParseBool(cfg.AWSConfig.EbsOptimized)
	volumeSize, err := strconv.Atoi(cfg.AWSConfig.VolumeSize)
	hasPublicAddress, err := strconv.ParseBool(cfg.AWSConfig.HasPublicAddr)

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
		ImageId:      &amiID,
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
						Key:   aws.String(clouds.ClusterIDTag),
						Value: aws.String(cfg.ClusterID),
					},
				},
			},
		},
	}
	if hasPublicAddress {
		runInstanceInput.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				AssociatePublicIpAddress: aws.Bool(true),
				DeleteOnTermination:      aws.Bool(true),
				SubnetId:                 aws.String(cfg.AWSConfig.SubnetID),
				Groups:                   []*string{secGroupID},
			},
		}
	}

	res, err := EC2.RunInstancesWithContext(ctx, runInstanceInput)
	if err != nil {
		cfg.Node.State = node.StateError
		cfg.NodeChan() <- cfg.Node

		log.Errorf("[%s] - failed to create ec2 instance: %v", StepNameCreateEC2Instance, err)
		return errors.Wrap(ErrCreateInstance, err.Error())
	}

	cfg.Node = node.Node{
		Name:     nodeName,
		TaskID:   cfg.TaskID,
		Region:   cfg.AWSConfig.Region,
		Role:     role,
		Provider: clouds.AWS,
		State:    node.StateBuilding,
	}

	// Update node state in cluster
	cfg.NodeChan() <- cfg.Node

	if len(res.Instances) == 0 {
		cfg.Node.State = node.StateError
		cfg.NodeChan() <- cfg.Node

		return errors.Wrap(ErrCreateInstance, "no instances created")
	}

	instance := res.Instances[0]

	if hasPublicAddress {
		log.Infof("[%s] - waiting to obtain public IP...", s.Name())

		//Waiting for AWS to assign public IP requires to poll an describe ec2 endpoint several times
		found := false
		for i := 0; i < IPAttempts; i++ {
			lookup := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("tag:Name"),
						Values: []*string{aws.String(nodeName)},
					},
					{
						Name:   aws.String(fmt.Sprintf("tag:%s", clouds.ClusterIDTag)),
						Values: []*string{aws.String(cfg.ClusterID)},
					},
				},
			}
			out, err := EC2.DescribeInstancesWithContext(ctx, lookup)
			if err != nil {
				cfg.Node.State = node.StateError
				cfg.NodeChan() <- cfg.Node
				log.Errorf("[%s] - failed to obtain public IP for node %s: %v", s.Name(), nodeName, err)
				return errors.Wrap(ErrNoPublicIP, err.Error())
			}

			if len(out.Reservations) == 0 {
				log.Infof("[%s] - found 0 ec2 instances, attempt %d", s.Name(), i)
				time.Sleep(time.Duration(SleepSecondsPerAttempt) * time.Second)
				continue
			}

			if i := findInstanceWithPublicAddr(out.Reservations); i != nil {
				cfg.Node.PublicIp = *i.PublicIpAddress
				cfg.Node.PrivateIp = *i.PrivateIpAddress
				log.Infof("[%s] - found public ip - %s for node %s", s.Name(), cfg.Node.PublicIp, nodeName)
				found = true
				break
			}
			// TODO(stgleb): Fix this wait loop
			time.Sleep(2 * time.Second)
		}
		if !found {
			log.Errorf("[%s] - failed to find public IP address after %d attempts", s.Name(), IPAttempts)
			cfg.Node.State = node.StateError
			cfg.NodeChan() <- cfg.Node
			return ErrNoPublicIP
		}
	}

	cfg.Node.Region = cfg.AWSConfig.Region
	cfg.Node.CreatedAt = instance.LaunchTime.Unix()
	cfg.Node.ID = *instance.InstanceId
	cfg.Node.State = node.StateProvisioning

	cfg.NodeChan() <- cfg.Node
	if cfg.IsMaster {
		cfg.AddMaster(&cfg.Node)
	} else {
		cfg.AddNode(&cfg.Node)
	}

	log.Infof("[%s] - success! Created node %s with instanceID %s ", s.Name(), nodeName, cfg.Node.ID)
	logrus.Debugf("%v", *instance)

	return nil
}

func (s *StepCreateInstance) FindAMI(ctx context.Context, w io.Writer, EC2 ec2iface.EC2API) (string, error) {
	out, err := EC2.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("architecture"),
				Values: []*string{
					aws.String("x86_64"),
				},
			},
			{
				Name: aws.String("virtualization-type"),
				Values: []*string{
					aws.String("hvm"),
				},
			},
			{
				Name: aws.String("root-device-type"),
				Values: []*string{
					aws.String("ebs"),
				},
			},
			//Owner should be Canonical
			{
				Name: aws.String("owner-id"),
				Values: []*string{
					aws.String("099720109477"),
				},
			},
			{
				Name: aws.String("description"),
				Values: []*string{
					aws.String("Canonical, Ubuntu, 16.04*"),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	amiID := ""

	log := util.GetLogger(w)

	for _, img := range out.Images {
		if img.Description == nil {
			continue
		}
		if strings.Contains(*img.Description, "UNSUPPORTED") {
			continue
		}
		amiID = *img.ImageId

		logMessage := fmt.Sprintf("[%s] - using AMI (ID: %s) %s", s.Name(), amiID, *img.Description)
		log.Info(logMessage)
		logrus.Info(logMessage)

		break
	}

	return amiID, nil
}

func (s *StepCreateInstance) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - rollback initiated", s.Name())

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	if cfg.Node.ID != "" {
		_, err := EC2.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
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
