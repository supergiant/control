package amazon

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteSecurityGroupsStepName = "aws_delete_security_groups"

var (
	deleteSecGroupTimeout      = time.Second * 10
	deleteSecGroupAttemptCount = 10
)

type deleteSecurityGroupService interface {
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	RevokeSecurityGroupIngressWithContext(aws.Context, *ec2.RevokeSecurityGroupIngressInput, ...request.Option) (*ec2.RevokeSecurityGroupIngressOutput, error)
	DeleteSecurityGroupWithContext(aws.Context, *ec2.DeleteSecurityGroupInput, ...request.Option) (*ec2.DeleteSecurityGroupOutput, error)
}

type DeleteSecurityGroup struct {
	getSvc func(steps.AWSConfig) (deleteSecurityGroupService, error)
}

func InitDeleteSecurityGroup(fn GetEC2Fn) {
	steps.RegisterStep(DeleteSecurityGroupsStepName,
		NewDeleteSecurityGroupService(fn))
}

func NewDeleteSecurityGroupService(fn GetEC2Fn) *DeleteSecurityGroup {
	return &DeleteSecurityGroup{
		getSvc: func(cfg steps.AWSConfig) (deleteSecurityGroupService, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteSecurityGroup) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.MastersSecurityGroupID == "" ||
		cfg.AWSConfig.NodesSecurityGroupID == "" {
		logrus.Debug("Skip deleting empty security groups")
		return nil
	}

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Getting service caused %v", err)
		return errors.Wrapf(err, "%s get service",
			DeleteSecurityGroupsStepName)
	}

	masterGroupName, err := s.getSecurityGroupNameByID(
		cfg.AWSConfig.MastersSecurityGroupID, svc)

	if err != nil {
		logrus.Debugf("get master security group ID %v", err)
		return errors.Wrapf(err, "get master security group ID")
	}

	logrus.Debugf("Master group name %s", masterGroupName)
	nodeGroupName, err := s.getSecurityGroupNameByID(
		cfg.AWSConfig.NodesSecurityGroupID, svc)

	if err != nil {
		logrus.Debugf("get node security group ID %v", err)
		return errors.Wrapf(err, "get node security group ID")
	}

	logrus.Debugf("Revoking dependent Node Security Group ingress rules %s", nodeGroupName)

	// Decouple security groups from each other
	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		IpPermissions: []*ec2.IpPermission{
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
		},
	}

	logrus.Debugf("Revoking dependent Master Security Group  %s ingress rules",
		cfg.AWSConfig.MastersSecurityGroupID)
	_, err = svc.RevokeSecurityGroupIngressWithContext(ctx, revokeInput)

	if err != nil {
		return errors.Wrapf(err, "revoke relation between master "+
			"and node security group %s caused %s",
			cfg.AWSConfig.NodesSecurityGroupID, err.Error())
	}

	revokeInput = &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		IpPermissions: []*ec2.IpPermission{
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
		},
	}

	_, err = svc.RevokeSecurityGroupIngressWithContext(ctx, revokeInput)

	if err != nil {
		logrus.Debugf("revoke relation between node and master "+
			"security group caused %s caused %s",
			cfg.AWSConfig.MastersSecurityGroupID, err.Error())
		return errors.Wrapf(err, "find security group %s",
			cfg.AWSConfig.NodesSecurityGroupID)
	}

	logrus.Debugf("Dependencies between security groups has been revoked")
	var deleteErr error
	var timeout = deleteSecGroupTimeout

	// Delete master security group first
	for i := 0; i < deleteSecGroupAttemptCount; i++ {
		logrus.Debugf("Delete master security group %s", cfg.AWSConfig.MastersSecurityGroupID)
		reqMaster := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		}
		_, deleteErr = svc.DeleteSecurityGroupWithContext(ctx, reqMaster)

		if deleteErr != nil {
			logrus.Debugf("delete master security group %s caused %s",
				cfg.AWSConfig.MastersSecurityGroupID, deleteErr.Error())
		} else {
			logrus.Debugf("master security group %s has been deleted",
				cfg.AWSConfig.MastersSecurityGroupID)
			break
		}

		logrus.Debugf("Sleep for %v", timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if deleteErr != nil {
		logrus.Errorf("Delete master security group %s", DeleteSecurityGroupsStepName)
		return errors.Wrapf(deleteErr, "%s delete master security group",
			DeleteSecurityGroupsStepName)
	}

	timeout = deleteSecGroupTimeout
	// Delete node security group
	for i := 0; i < deleteSecGroupAttemptCount; i++ {
		reqNode := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		}

		logrus.Debugf("Delete node security group %s",
			cfg.AWSConfig.NodesSecurityGroupID)
		_, deleteErr = svc.DeleteSecurityGroupWithContext(ctx, reqNode)

		if deleteErr != nil {
			logrus.Debugf("delete node security group %s %s",
				cfg.AWSConfig.NodesSecurityGroupID, deleteErr.Error())
		} else {
			logrus.Debugf("node security group %s has been deleted",
				cfg.AWSConfig.NodesSecurityGroupID)
			break
		}

		logrus.Debugf("Sleep for %v", timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}


	if deleteErr != nil {
		logrus.Errorf("Delete node security group %s", DeleteSecurityGroupsStepName)
		return errors.Wrapf(deleteErr, "%s delete node security group",
			DeleteSecurityGroupsStepName)
	}

	// Don't fail even if something not get deleted
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

func (*DeleteSecurityGroup) getSecurityGroupNameByID(name string, svc deleteSecurityGroupService) (string, error) {
	req := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{&name},
	}

	resp, err := svc.DescribeSecurityGroups(req)
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
