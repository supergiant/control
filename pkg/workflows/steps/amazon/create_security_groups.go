package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateSecurityGroups = "create_security_groups_step"

type CreateSecurityGroupsStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateSecurityGroupsStep(fn GetEC2Fn) *CreateSecurityGroupsStep {
	return &CreateSecurityGroupsStep{
		GetEC2: fn,
	}
}

func InitCreateSecurityGroups(fn GetEC2Fn) {
	steps.RegisterStep(StepCreateSecurityGroups, NewCreateSecurityGroupsStep(fn))
}

func (s *CreateSecurityGroupsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return ErrAuthorization
	}

	if cfg.AWSConfig.VPCID == "" {
		err := errors.New("no vpc id provided for security groups creation")
		log.Errorf("[%s] - %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Create security groups for VPC %s",
		cfg.AWSConfig.VPCID)
	if cfg.AWSConfig.MastersSecurityGroupID == "" {
		groupName := fmt.Sprintf("%s-masters-secgroup", cfg.ClusterID)

		log.Infof("[%s] - masters security groups not specified, will create a new one...", s.Name())
		out, err := EC2.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes masters for cluster " + cfg.ClusterID),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(groupName),
		})
		if err != nil {
			duplicateErr := false
			if err, ok := err.(awserr.Error); ok {
				logrus.Debugf("Create security group for masters caused %s",
					err.Message())
				if err.Code() == "InvalidGroup.Duplicate" {
					duplicateErr = true
				}
			}
			if !duplicateErr {
				log.Errorf("[%s] - create security groups for k8s masters: %v", s.Name(), err)
				return err
			}

			groupID, err := s.getSecurityGroupByName(ctx, EC2, cfg.AWSConfig.VPCID, groupName)
			if err != nil {
				return err
			}

			cfg.AWSConfig.MastersSecurityGroupID = groupID
		} else {
			cfg.AWSConfig.MastersSecurityGroupID = *out.GroupId
		}
	}
	//If there is no security group, create it
	if cfg.AWSConfig.NodesSecurityGroupID == "" {
		groupName := fmt.Sprintf("%s-nodes-secgroup", cfg.ClusterID)

		log.Infof("[%s] - node security groups not specified, will create a new one...", s.Name())
		out, err := EC2.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes nodes for cluster " + cfg.ClusterID),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(groupName),
		})
		if err != nil {
			duplicateErr := false
			if err, ok := err.(awserr.Error); ok {
				logrus.Debugf("Create security group for nodes caused %s",
					err.Message())
				if err.Code() == "InvalidGroup.Duplicate" {
					duplicateErr = true
				}
			}
			if !duplicateErr {
				log.Errorf("[%s] - create security groups for k8s nodes: %v", s.Name(), err)
				return err
			}

			groupID, err := s.getSecurityGroupByName(ctx, EC2, cfg.AWSConfig.VPCID, groupName)
			if err != nil {
				return err
			}

			cfg.AWSConfig.NodesSecurityGroupID = groupID
		} else {
			cfg.AWSConfig.NodesSecurityGroupID = *out.GroupId
		}
	}

	logrus.Debugf("Security groups %s %s has been created",
		cfg.AWSConfig.MastersSecurityGroupID, cfg.AWSConfig.NodesSecurityGroupID)

	logrus.Debugf("Authorize SSH between groups")
	//In order to deploy the kubernetes cluster supergiant needs to open port 22
	if err := s.authorizeSSH(ctx, EC2, cfg.AWSConfig.MastersSecurityGroupID); err != nil {
		return err
	}

	if err := s.authorizeSSH(ctx, EC2, cfg.AWSConfig.NodesSecurityGroupID); err != nil {
		return err
	}

	logrus.Debugf("Allow traffic between groups")
	//Open ports between master <-> node security groups
	// nodes to nodes
	if err := s.allowAllTraffic(ctx, EC2, cfg); err != nil {
		return err
	}

	logrus.Debugf("Whitelist SG IP address")
	if err := s.whiteListSupergiantIP(ctx, EC2, cfg.AWSConfig.MastersSecurityGroupID); err != nil {
		logrus.Errorf("[%s] - failed to whitelist supergiant IP in master "+
			"security group: %v", s.Name(), err)
	}

	return nil
}

func (s *CreateSecurityGroupsStep) getSecurityGroupByName(ctx context.Context, EC2 ec2iface.EC2API, vpcID, groupName string) (string, error) {
	logrus.Debugf("Get security group %s by name in vpc %s",
		groupName, vpcID)
	out, err := EC2.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: aws.StringSlice([]string{groupName}),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{vpcID}),
			},
		},
	})

	if err != nil {
		return "", err
	}
	if len(out.SecurityGroups) == 0 {
		return "", errors.Errorf("no security group %s found in aws", groupName)
	}
	logrus.Debugf("Found %s %s", groupName, *out.SecurityGroups[0].VpcId)
	return *out.SecurityGroups[0].GroupId, nil
}

func (s *CreateSecurityGroupsStep) getSecurityGroupById(ctx context.Context, EC2 ec2iface.EC2API, vpcID, groupID string) (*ec2.SecurityGroup, error) {
	logrus.Debugf("Get security group by Id %s in VPC %s",
		groupID, vpcID)
	out, err := EC2.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: aws.StringSlice([]string{groupID}),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{vpcID}),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.SecurityGroups) == 0 {
		return nil, errors.Errorf("no security group with ID %s found in aws", groupID)
	}
	return out.SecurityGroups[0], nil
}

func (s *CreateSecurityGroupsStep) authorizeSSH(ctx context.Context, EC2 ec2iface.EC2API, groupID string) error {
	_, err := EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(groupID),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		CidrIp:     aws.String("0.0.0.0/0"),
		IpProtocol: aws.String("tcp"),
	})
	if err, ok := err.(awserr.Error); ok {
		if err.Code() == "InvalidPermission.Duplicate" {
			return nil
		}
	}
	return err
}

func (s *CreateSecurityGroupsStep) allowAllTraffic(ctx context.Context, EC2 ec2iface.EC2API, cfg *steps.Config) error {
	_, err := EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(0),
				IpProtocol: aws.String("-1"),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
					},
				},
			},
			{
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(0),
				IpProtocol: aws.String("-1"),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
					},
				},
			},
		},
	})

	if err, ok := err.(awserr.Error); ok {
		if err.Code() == "InvalidPermission.Duplicate" {
			return nil
		}
	}

	_, err = EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(0),
				IpProtocol: aws.String("-1"),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
					},
				},
			},
			{
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(0),
				IpProtocol: aws.String("-1"),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
					},
				},
			},
		},
	})

	if err, ok := err.(awserr.Error); ok {
		if err.Code() == "InvalidPermission.Duplicate" {
			return nil
		}
	}

	return err
}

func (s *CreateSecurityGroupsStep) whiteListSupergiantIP(ctx context.Context, EC2 ec2iface.EC2API, groupID string) error {
	supergiantIP, err := FindOutboundIP(ctx, findOutBoundIP)
	if err != nil {
		return err
	}

	_, err = EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(groupID),
		FromPort:   aws.Int64(8080),
		ToPort:     aws.Int64(8080),
		CidrIp:     aws.String(fmt.Sprintf("%s/32", supergiantIP)),
		IpProtocol: aws.String("tcp"),
	})
	if err != nil {
		return err
	}

	_, err = EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(groupID),
		FromPort:   aws.Int64(443),
		ToPort:     aws.Int64(443),
		CidrIp:     aws.String(fmt.Sprintf("%s/32", supergiantIP)),
		IpProtocol: aws.String("tcp"),
	})
	if err != nil {
		return err
	}

	return err
}

func (*CreateSecurityGroupsStep) Name() string {
	return StepCreateSecurityGroups
}

func (*CreateSecurityGroupsStep) Description() string {
	return ""
}

func (*CreateSecurityGroupsStep) Depends() []string {
	return nil
}

func (*CreateSecurityGroupsStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
