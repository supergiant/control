package aws

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/pkg/errors"
)

// Tag keys for aws resources:
const (
	TagName    = "Name"
	TagCluster = "KubernetesCluster"
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
	session  *session.Session
	ec2SvcFn func(s *session.Session, region string) ec2iface.EC2API
	tags     map[string]string
}

// New returns a configured AWS client.
func New(keyID, secret string, tags map[string]string) (*Client, error) {
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
		tags:     tags,
	}, nil
}

// AvailableInstanceTypes returns a list of valid instance types for the region.
func (c *Client) AvailableInstanceTypes(ctx context.Context) ([]*EC2TypeInfo, error) {
	ec2Types, err := c.getEC2Types()
	if err != nil {
		return nil, errors.Wrap(err, "aws: get ec2 types")
	}

	ec2Infos := make([]*EC2TypeInfo, 0, len(ec2Types))
	for _, t := range ec2Types {
		info, err := c.getEC2TypeInfo(t)
		if err != nil {
			return nil, errors.Wrapf(err, "aws: get ec2 %s info", t)
		}
		ec2Infos = append(ec2Infos, info)
	}

	return ec2Infos, nil
}

// CreateInstance startd a new instance due to the config.
func (c *Client) CreateInstance(ctx context.Context, cfg InstanceConfig) error {
	cfg.Region = strings.TrimSpace(cfg.Region)
	if cfg.Region == "" {
		return ErrNoRegionProvided
	}
	ec2S := c.ec2SvcFn(c.session, cfg.Region)

	instanceInp := &ec2.RunInstancesInput{
		ImageId:      aws.String(cfg.ImageID),
		InstanceType: aws.String(cfg.Type),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String(cfg.KeyName),
		EbsOptimized: aws.Bool(true),
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
		SecurityGroupIds: cfg.SecurityGroups,
		SubnetId:         aws.String(cfg.SubnetID),
	}
	if cfg.UsedData != "" {
		instanceInp.UserData = aws.String(cfg.UsedData)
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
	return tagAWSResource(ec2S, *(res.Instances[0].InstanceId), getTags(c.tags, &cfg))
}

// ListInstances returns a list of instances available to the client.
func (c *Client) ListRegionInstances(ctx context.Context, region string, tags map[string]string) ([]*ec2.Instance, error) {
	region = strings.TrimSpace(region)
	if region == "" {
		return nil, ErrNoRegionProvided
	}

	var token *string
	instList := make([]*ec2.Instance, 0)
	for {
		out, err := c.ec2SvcFn(c.session, region).DescribeInstancesWithContext(ctx,
			&ec2.DescribeInstancesInput{NextToken: token, Filters: c.buildFilter(tags)})
		if err != nil {
			return nil, err
		}
		for _, reservation := range out.Reservations {
			instList = append(instList, reservation.Instances...)
		}
		if out.NextToken == nil {
			break
		}
		token = out.NextToken
	}

	return instList, nil
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

func (c *Client) getEC2Types() ([]string, error) {
	svc := pricing.New(c.session, &aws.Config{Region: aws.String("us-east-1")})

	var nextToken *string
	instanceTypes := make([]string, 0)
	for {
		input := &pricing.GetAttributeValuesInput{
			AttributeName: aws.String("instanceType"),
			MaxResults:    aws.Int64(100),
			ServiceCode:   aws.String("AmazonEC2"),
			NextToken:     nextToken,
		}

		result, err := svc.GetAttributeValues(input)
		if err != nil {
			return nil, err
		}

		for _, v := range result.AttributeValues {
			if *v.Value != "" {
				instanceTypes = append(instanceTypes, *v.Value)
			}
		}

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return instanceTypes, nil
}

func (c *Client) getEC2TypeInfo(instanceType string) (*EC2TypeInfo, error) {
	svc := pricing.New(c.session, &aws.Config{Region: aws.String("us-east-1")})

	productsInput := &pricing.GetProductsInput{
		Filters: []*pricing.Filter{
			{
				Field: aws.String("ServiceCode"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("AmazonEC2"),
			},
			{
				Field: aws.String("instanceType"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(instanceType),
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(3),
		ServiceCode:   aws.String("AmazonEC2"),
	}

	productsRes, err := svc.GetProducts(productsInput)
	if err != nil {
		return nil, err
	}

	info := &EC2TypeInfo{}
	for _, product := range productsRes.PriceList {
		if p, ok := product["product"]; ok {
			// TODO: replace with reflect/type assetion
			d, err := json.Marshal(p)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal(d, info); err != nil {
				return nil, err
			}
			break
		}
	}

	return info, nil
}

func getTags(clientTags map[string]string, cfg *InstanceConfig) map[string]string {
	tags := make(map[string]string)

	// default set of tags
	tags[TagName] = cfg.Name
	tags[TagCluster] = cfg.ClusterName
	tags[TagRole] = cfg.ClusterRole

	// apply client tags
	for k, v := range clientTags {
		tags[k] = v
	}

	// apply instance config tags
	for k, v := range cfg.Tags {
		tags[k] = v
	}

	return tags
}

func (c *Client) buildFilter(tags map[string]string) []*ec2.Filter {
	if len(c.tags) == 0 && len(tags) == 0 {
		return nil
	}
	filters := make([]*ec2.Filter, 0)
	for k, v := range c.tags {
		filters = append(filters, &ec2.Filter{Name: aws.String("tag:" + k), Values: []*string{aws.String(v)}})
	}
	for k, v := range tags {
		filters = append(filters, &ec2.Filter{Name: aws.String("tag:" + k), Values: []*string{aws.String(v)}})
	}
	return filters
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
