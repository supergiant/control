package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
)

const DeleteSecurityGroupsStepName = "aws_delete_security_groups"

type DeleteSecurityGroup struct {
	GetEC2 GetEC2Fn
}

func InitDeleteSecurityGroup(fn GetEC2Fn) {
	steps.RegisterStep(DeleteSecurityGroupsStepName, &DeleteSecurityGroup{
		GetEC2: fn,
	})
}

func (s *DeleteSecurityGroup) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	masterGroupName, err := s.getSecurityGroupNameByID(
		cfg.AWSConfig.MastersSecurityGroupID, EC2)

	if err != nil {
		return errors.Wrapf(err, "get master security group ID")
	}
	logrus.Debugf("Master group name %s", masterGroupName)

	nodeGroupName, err := s.getSecurityGroupNameByID(
		cfg.AWSConfig.NodesSecurityGroupID, EC2)

	if err != nil {
		return errors.Wrapf(err, "get node security group ID")
	}

	logrus.Debugf("Node group name %s", nodeGroupName)

	// Decouple security groups from each other
	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:   aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		GroupName: aws.String(nodeGroupName),
	}

	output, err := EC2.RevokeSecurityGroupIngress(revokeInput)

	logrus.Debugf(output.String())
	if err != nil {
		return errors.Wrapf(err, "find security group %s",
			cfg.AWSConfig.NodesSecurityGroupID)
	}

	revokeInput = &ec2.RevokeSecurityGroupIngressInput{
		GroupId:   aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		GroupName: aws.String(masterGroupName),
	}

	_, err = EC2.RevokeSecurityGroupIngress(revokeInput)

	if err != nil {
		return errors.Wrapf(err, "find security group %s",
			cfg.AWSConfig.NodesSecurityGroupID)
	}

	logrus.Debugf("Delete master security group %s", cfg.AWSConfig.MastersSecurityGroupID)
	reqMaster := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
	}

	_, err = EC2.DeleteSecurityGroup(reqMaster)

	if err != nil {
		return errors.Wrapf(err, "%s: master security groups", DeleteSecurityGroupsStepName)
	}

	reqNode := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
	}

	logrus.Debugf("Delete node security group %s", cfg.AWSConfig.NodesSecurityGroupID)
	_, err = EC2.DeleteSecurityGroup(reqNode)

	if err != nil {
		return errors.Wrapf(err, "%s: nodes security groups", DeleteSecurityGroupsStepName)
	}

	logrus.Debugf("Deleting security group finished")
	return nil
}

func (*DeleteSecurityGroup) Name() string {
	return DeleteSecurityGroupsStepName
}

func (*DeleteSecurityGroup) Depends() []string {
	return []string{DeleteClusterMachinesStepName}
}

func (*DeleteSecurityGroup) Description() string {
	return "Deletes security groups"
}

func (*DeleteSecurityGroup) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (*DeleteSecurityGroup) getSecurityGroupNameByID(name string, EC2 ec2iface.EC2API) (string, error) {
	req := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{&name},
	}

	resp, err := EC2.DescribeSecurityGroups(req)
	if err != nil {
		return "", errors.Wrapf(err,
			"find security group %s", name)
	}

	if len(resp.SecurityGroups) == 0 {
		return "", errors.Wrapf(sgerrors.ErrNotFound,
			"find security group %s", name)
	}

	return *resp.SecurityGroups[0].GroupName, nil
}
