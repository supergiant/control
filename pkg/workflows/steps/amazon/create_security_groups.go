package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const CreateSecurityGroupsStepName = "create_security_groups_step"

type CreateSecurityGroupsStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateSecurityGroupsStep(fn GetEC2Fn) *CreateSecurityGroupsStep {
	return &CreateSecurityGroupsStep{
		GetEC2: fn,
	}
}

func InitCreateSecurityGroups(fn GetEC2Fn) {
	steps.RegisterStep(CreateSecurityGroupsStepName, NewCreateSecurityGroupsStep(fn))
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

	if cfg.AWSConfig.MastersSecurityGroup == "" {
		log.Infof("[%s] - masters security groups not specified, will create a new one...", s.Name())
		out, err := EC2.CreateSecurityGroupWithContext(ctx, &ec2.CreateSecurityGroupInput{
			Description: aws.String("Security group for Kubernetes masters for cluster " + cfg.ClusterName),
			VpcId:       aws.String(cfg.AWSConfig.VPCID),
			GroupName:   aws.String(fmt.Sprintf("%s-masters-secgroup", cfg.ClusterName)),
		})
		if err != nil {
			log.Errorf("[%s] - create security groups for k8s masters: %v", s.Name(), err)
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
			log.Errorf("[%s] - create security groups for k8s nodes: %v", s.Name(), err)
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
