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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"time"
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

	if err, ok := err.(awserr.Error); ok {
			logrus.Debugf("get master security group ID %v", err)
			return errors.Wrapf(err, "get master security group ID")
	}
	logrus.Debugf("Master group name %s", masterGroupName)

	nodeGroupName, err := s.getSecurityGroupNameByID(
		cfg.AWSConfig.NodesSecurityGroupID, EC2)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("get node security group ID %s", err.Message())
		return errors.Wrapf(err, "get node security group ID")
	}

	logrus.Debugf("Node group name %s", nodeGroupName)

	// Decouple security groups from each other
	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:   aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		SourceSecurityGroupName: aws.String(nodeGroupName),
	}

	logrus.Debugf("Revoke relation between master and node security groups")
	_, err = EC2.RevokeSecurityGroupIngressWithContext(ctx, revokeInput)

	if err, ok := err.(awserr.Error); ok {
		return errors.Wrapf(err, "revoke relation between master " +
			"and node security group %s caused %s",
			cfg.AWSConfig.NodesSecurityGroupID, err.Message())
	}

	revokeInput = &ec2.RevokeSecurityGroupIngressInput{
		GroupId:   aws.String(cfg.AWSConfig.NodesSecurityGroupID),
		SourceSecurityGroupName: aws.String(masterGroupName),
	}

	logrus.Debugf("Revoke relation between node and master security groups")
	_, err = EC2.RevokeSecurityGroupIngressWithContext(ctx, revokeInput)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("revoke relation between node and master " +
			"security group caused %s caused %s",
			cfg.AWSConfig.MastersSecurityGroupID, err.Message())
		return errors.Wrapf(err, "find security group %s",
			cfg.AWSConfig.NodesSecurityGroupID)
	}

	referenceInput := &ec2.DescribeSecurityGroupReferencesInput{
		GroupId: []*string{aws.String(cfg.AWSConfig.MastersSecurityGroupID),
		aws.String(cfg.AWSConfig.NodesSecurityGroupID)},
	}

	references, err := EC2.DescribeSecurityGroupReferences(referenceInput)
	for _, ref := range references.SecurityGroupReferenceSet {
		logrus.Debugf("GroupID: %s ReferenceGroup: %s", *ref.GroupId,
			*ref.ReferencingVpcId)
	}

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("DescribeSecurityGroupReferences caused %s", err.Error())
		return errors.Wrapf(err, "DescribeSecurityGroupReferences:",
			DeleteSecurityGroupsStepName)
	}

	// TODO(stgleb): Yes, this needs to be more interactive
	time.Sleep(time.Minute * 1)

	references, err = EC2.DescribeSecurityGroupReferences(referenceInput)
	for _, ref := range references.SecurityGroupReferenceSet {
		logrus.Debugf("GroupID: %s ReferenceGroup: %s", *ref.GroupId,
			*ref.ReferencingVpcId)
	}


	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("DescribeSecurityGroupReferences caused %s", err.Error())
		return errors.Wrapf(err, "DescribeSecurityGroupReferences:",
			DeleteSecurityGroupsStepName)
	}

	logrus.Debugf("Delete master security group %s", cfg.AWSConfig.MastersSecurityGroupID)
	reqMaster := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
	}

	_, err = EC2.DeleteSecurityGroupWithContext(ctx, reqMaster)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("delete master security group %s caused %s",
			aws.String(cfg.AWSConfig.MastersSecurityGroupID), err.Error())
		return errors.Wrapf(err, "%s: master security groups",
			DeleteSecurityGroupsStepName)
	}


	reqNode := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
	}

	logrus.Debugf("Delete node security group %s",
		cfg.AWSConfig.NodesSecurityGroupID)
	_, err = EC2.DeleteSecurityGroupWithContext(ctx, reqNode)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("delete node security group %s %s",
			aws.String(cfg.AWSConfig.NodesSecurityGroupID), err.Message())
		return errors.Wrapf(err, "%s: nodes security groups",
			DeleteSecurityGroupsStepName)
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
