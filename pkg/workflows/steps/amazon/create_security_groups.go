package amazon

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const CreateSecurityGroupsStepName = "create_security_groups_step"

type CreateSecurityGroupsStep struct {
	GetEC2 GetEC2Fn
}

func (s *CreateSecurityGroupsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	if cfg.AWSConfig.VPCID == "" {
		err := errors.New("aws: no vpc id provided for security groups")
		log.Errorf("[%s] %v", s.Name(), err)
		return err
	}

	if cfg.AWSConfig.MastersSecurityGroup == "" {
		log.Infof("[%s] - masters security groups not specified, will create a new one...")
		out, err := EC2.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes masters for cluster " + cfg.ClusterName),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(fmt.Sprintf("%s-masters-secgroup", cfg.ClusterName)),
		})
		if err != nil {
			log.Error(err)
			return err
		}
		cfg.AWSConfig.MastersSecurityGroup = *out.GroupId
	}
	//If there is no security group, create it
	if cfg.AWSConfig.NodesSecurityGroup == "" {
		log.Infof("[%s] - node security groups not specified, will create a new one...")
		out, err := EC2.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes nodes for cluster " + cfg.ClusterName),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(fmt.Sprintf("%s-nodes-secgroup", cfg.ClusterName)),
		})
		if err != nil {
			log.Error(err)
			return err
		}
		cfg.AWSConfig.NodesSecurityGroup = *out.GroupId
	}
	return nil
}

func (*CreateSecurityGroupsStep) Name() string {
	return CreateSecurityGroupsStepName
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
