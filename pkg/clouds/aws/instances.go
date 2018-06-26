package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

// Tag keys for aws resources:
const (
	TagCluster = "Cluster"
	TagRole    = "Role"
)

// AWS client specific errors:
var (
	ErrInvalidKeys        = errors.New("aws: invalid keys")
	ErrNoRegionProvided   = errors.New("aws: region should shouldn't be emplty")
	ErrInstanceIDEmpty    = errors.New("aws: instance id shouldn't be emplty")
	ErrNoInstancesCreated = errors.New("aws: no instances were created")
)

type Client struct {
	session *session.Session

	ec2SvcFn func(s *session.Session, region string) ec2iface.EC2API
}

// New returns a configured AWS client.
func New(keyID, secret string) (*Client, error) {
	keyID, secret = strings.TrimSpace(keyID), strings.TrimSpace(secret)
	if keyID == "" || secret == "" {
		return nil, ErrInvalidKeys
	}

	session := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(keyID, secret, ""),
	})

	return &Client{
		session:  session,
		ec2SvcFn: ec2Svc,
	}, nil
}

// CreateInstance start a new instance due to the config.
func (c *Client) CreateInstance(ctx context.Context, cfg InstanceConfig) error {
	cfg.Region = strings.TrimSpace(cfg.Region)
	if cfg.Region == "" {
		return ErrNoRegionProvided
	}
	ec2S := c.ec2SvcFn(c.session, cfg.Region)

	instanceInp := &ec2.RunInstancesInput{
		ImageId:      aws.String(cfg.ImageID),
		InstanceType: aws.String(ec2.InstanceTypeT2Micro),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String(cfg.KeyName),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(cfg.IAMRole),
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeType:          aws.String(cfg.VolumeType),
					VolumeSize:          aws.Int64(cfg.VolumeSize),
				},
			},
		},
	}
	if cfg.HasPublicAddr {
		instanceInp.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				AssociatePublicIpAddress: aws.Bool(cfg.HasPublicAddr),
				DeleteOnTermination:      aws.Bool(true),
				Groups:                   cfg.SecurityGroups,
				SubnetId:                 aws.String(cfg.SubnetID),
			},
		}
	}

	res, err := ec2S.RunInstancesWithContext(ctx, instanceInp)
	if err != nil {
		return errors.Wrap(err, "aws: run instance")
	}
	if res.Instances == nil || len(res.Instances) < 1 {
		return ErrNoInstancesCreated
	}

	// add some metadata info to an instance
	if cfg.Tags == nil {
		cfg.Tags = make(map[string]string)
	}
	cfg.Tags[TagCluster] = cfg.ClusterName
	cfg.Tags[TagRole] = cfg.ClusterRole

	return tagAWSResource(ec2S, *(res.Instances[0].InstanceId), cfg.Tags)
}

// DeleteInstance terminates an instance with provided id and region.
func (c *Client) DeleteInstance(ctx context.Context, region, instanceID string) error {
	region, instanceID = strings.TrimSpace(region), strings.TrimSpace(instanceID)
	if region == "" {
		return ErrNoRegionProvided
	}
	if instanceID == "" {
		return ErrInstanceIDEmpty
	}
	ec2S := c.ec2SvcFn(c.session, region)

	_, err := ec2S.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})

	return errors.Wrap(err, "aws: terminate instance")
}

func ec2Svc(s *session.Session, region string) ec2iface.EC2API {
	return ec2.New(s, aws.NewConfig().WithRegion(region))
}

func tagAWSResource(ec2S ec2iface.EC2API, id string, tags map[string]string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInstanceIDEmpty
	}
	if len(tags) == 0 {
		return nil
	}

	awsTags := make([]*ec2.Tag, len(tags))
	for k, v := range tags {
		awsTags = append(awsTags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	_, err := ec2S.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{aws.String(id)},
		Tags:      awsTags,
	})
	return errors.Wrap(err, "aws: tag resource")
}
