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

// AWS provider specific errors:
var (
	ErrInvalidKeys        = errors.New("aws: invalid keys")
	ErrNoRegionProvided   = errors.New("aws: region should shouldn't be emplty")
	ErrInstanceIDEmpty    = errors.New("aws: instance id shouldn't be emplty")
	ErrNoInstancesCreated = errors.New("aws: no instances were created")
)

type Provider struct {
	session *session.Session

	ec2SvcFn func(s *session.Session, region string) ec2iface.EC2API
}

// New returns a configured AWS Provider.
func New(keyID, secret string) (*Provider, error) {
	if strings.TrimSpace(keyID) == "" || strings.TrimSpace(secret) == "" {
		return nil, ErrInvalidKeys
	}

	session := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(keyID, secret, ""),
	})

	return &Provider{
		session:  session,
		ec2SvcFn: ec2Svc,
	}, nil
}

// CreateInstance start a new instance due to the provided config.
func (p *Provider) CreateInstance(ctx context.Context, c InstanceConfig) error {
	if strings.TrimSpace(c.Region) == "" {
		return ErrNoRegionProvided
	}
	ec2S := p.ec2SvcFn(p.session, c.Region)

	instanceInp := &ec2.RunInstancesInput{
		ImageId:      aws.String(c.ImageID),
		InstanceType: aws.String(ec2.InstanceTypeT2Micro),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String(c.KeyName),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(c.IAMRole),
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeType:          aws.String(c.VolumeType),
					VolumeSize:          aws.Int64(c.VolumeSize),
				},
			},
		},
	}
	if c.HasPublicAddr {
		instanceInp.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				AssociatePublicIpAddress: aws.Bool(c.HasPublicAddr),
				DeleteOnTermination:      aws.Bool(true),
				Groups:                   c.SecurityGroups,
				SubnetId:                 aws.String(c.SubnetID),
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
	if c.Tags == nil {
		c.Tags = make(map[string]string)
	}
	c.Tags[TagCluster] = c.ClusterName
	c.Tags[TagRole] = c.ClusterRole

	return tagAWSResource(ec2S, *(res.Instances[0].InstanceId), c.Tags)
}

// DeleteInstance terminates an instance with provided id and region.
func (p *Provider) DeleteInstance(ctx context.Context, region, instanceID string) error {
	if strings.TrimSpace(region) == "" {
		return ErrNoRegionProvided
	}
	if strings.TrimSpace(instanceID) == "" {
		return ErrInstanceIDEmpty
	}
	ec2S := p.ec2SvcFn(p.session, region)

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
	if strings.TrimSpace(id) == "" {
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
