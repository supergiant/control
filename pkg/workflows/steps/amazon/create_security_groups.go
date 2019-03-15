package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateSecurityGroups = "create_security_groups_step"

type secGroupService interface {
	CreateSecurityGroupWithContext(aws.Context, *ec2.CreateSecurityGroupInput, ...request.Option) (*ec2.CreateSecurityGroupOutput, error)
	AuthorizeSecurityGroupIngressWithContext(aws.Context, *ec2.AuthorizeSecurityGroupIngressInput, ...request.Option) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}

type CreateSecurityGroupsStep struct {
	getSvc         func(config steps.AWSConfig) (secGroupService, error)
	findOutboundIP func() (string, error)
}

func NewCreateSecurityGroupsStep(fn GetEC2Fn) *CreateSecurityGroupsStep {
	return &CreateSecurityGroupsStep{
		getSvc: func(config steps.AWSConfig) (secGroupService, error) {
			EC2, err := fn(config)
			if err != nil {
				return nil, ErrAuthorization
			}

			return EC2, nil
		},
		findOutboundIP: FindExternalIP,
	}
}

func InitCreateSecurityGroups(fn GetEC2Fn) {
	steps.RegisterStep(StepCreateSecurityGroups, NewCreateSecurityGroupsStep(fn))
}

func (s *CreateSecurityGroupsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("%s: Getting service caused %v",
			StepCreateSecurityGroups, err)
		return errors.Wrapf(err, "%s get service", StepCreateSecurityGroups)
	}

	logrus.Debugf("Create security groups for VPC %s",
		cfg.AWSConfig.VPCID)
	if cfg.AWSConfig.MastersSecurityGroupID == "" {
		groupName := fmt.Sprintf("%s-masters-secgroup", cfg.ClusterID)

		log.Infof("[%s] - masters security groups not specified, will create a new one...", s.Name())
		out, err := svc.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes masters for cluster " + cfg.ClusterID),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(groupName),
		})
		if err != nil {
			return errors.Wrapf(err, "create master security group")
		} else {
			cfg.AWSConfig.MastersSecurityGroupID = *out.GroupId
		}
	}
	//If there is no security group, create it
	if cfg.AWSConfig.NodesSecurityGroupID == "" {
		groupName := fmt.Sprintf("%s-nodes-secgroup", cfg.ClusterID)

		log.Infof("[%s] - node security groups not specified, will create a new one...", s.Name())
		out, err := svc.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes nodes for cluster " + cfg.ClusterID),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(groupName),
		})
		if err != nil {
			return errors.Wrapf(err, "create node security group")
		} else {
			cfg.AWSConfig.NodesSecurityGroupID = *out.GroupId
		}
	}

	logrus.Debugf("Security groups %s %s has been created",
		cfg.AWSConfig.MastersSecurityGroupID, cfg.AWSConfig.NodesSecurityGroupID)

	logrus.Debugf("Authorize SSH between groups")
	//In order to deploy the kubernetes cluster supergiant needs to open port 22
	if err := s.authorizeSSH(ctx, svc, cfg.AWSConfig.MastersSecurityGroupID); err != nil {
		logrus.Errorf("authorize ssh for masters caused %v", err)
		return errors.Wrapf(err, "%s authorize ssh for masters",
			StepCreateSecurityGroups)
	}

	if err := s.authorizeSSH(ctx, svc, cfg.AWSConfig.NodesSecurityGroupID); err != nil {
		logrus.Errorf("authorize ssh for nodes caused %v", err)
		return errors.Wrapf(err, "%s authorize ssh for nodes",
			StepCreateSecurityGroups)
	}

	logrus.Debugf("Allow traffic between groups")
	//Open ports between master <-> node security groups
	// nodes to nodes
	if err := s.allowAllTraffic(ctx, svc, cfg); err != nil {
		return err
	}

	logrus.Debugf("Whitelist SG IP address")
	if err := s.whiteListSupergiantIP(ctx, svc, cfg.AWSConfig.MastersSecurityGroupID); err != nil {
		logrus.Errorf("[%s] - failed to whitelist supergiant IP in master "+
			"security group: %v", s.Name(), err)
		return errors.Wrapf(err, "%s failed whitelisting supergiant IP", s.Name())
	}

	return nil
}

func (s *CreateSecurityGroupsStep) authorizeSSH(ctx context.Context, EC2 secGroupService, groupID string) error {
	_, err := EC2.AuthorizeSecurityGroupIngressWithContext(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(groupID),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		CidrIp:     aws.String("0.0.0.0/0"),
		IpProtocol: aws.String("tcp"),
	})

	return err
}

func (s *CreateSecurityGroupsStep) allowAllTraffic(ctx context.Context, EC2 secGroupService, cfg *steps.Config) error {
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

	if err != nil {
		return err
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

	return err
}

func (s *CreateSecurityGroupsStep) whiteListSupergiantIP(ctx context.Context, EC2 secGroupService, groupID string) error {
	supergiantIP, err := FindOutboundIP(ctx, s.findOutboundIP)
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
	return "Create security groups"
}

func (*CreateSecurityGroupsStep) Depends() []string {
	return nil
}

func (*CreateSecurityGroupsStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
