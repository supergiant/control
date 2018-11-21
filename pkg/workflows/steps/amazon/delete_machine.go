package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteNodeStepName = "aws_delete_node"

type DeleteNodeStep struct {
	GetEC2 GetEC2Fn
}

func InitDeleteNode(fn GetEC2Fn) {
	steps.RegisterStep(DeleteNodeStepName, &DeleteNodeStep{
		GetEC2: fn,
	})
}

func (s *DeleteNodeStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	logrus.Infof("[%s] - deleting node %s", s.Name(), cfg.Node.Name)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Get instance by name filter %s", cfg.Node.Name)
	describeInstanceOutput, err := EC2.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: aws.StringSlice([]string{cfg.Node.Name}),
			},
		},
	})
	if err != nil {
		return errors.Wrap(ErrDeleteNode, err.Error())
	}

	logrus.Debugf("Got %d instance outputs",
		len(describeInstanceOutput.Reservations))
	instanceIDS := make([]string, 0)
	for _, res := range describeInstanceOutput.Reservations {
		for _, instance := range res.Instances {
			instanceIDS = append(instanceIDS, *instance.InstanceId)
		}
	}
	if len(instanceIDS) == 0 {
		logrus.Infof("[%s] - node %s not found in cluster %s",
			s.Name(), cfg.Node.Name, cfg.ClusterName)
		return nil
	}

	logrus.Debugf("Node to be deleted Name: %s AWS id: %v",
		cfg.Node.Name, instanceIDS)
	output, err := EC2.TerminateInstancesWithContext(ctx,
		&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice(instanceIDS),
		})

	if err != nil {
		return errors.Wrap(err, "terminate instance")
	}

	terminatedInstances := make([]string, 0)
	for _, instance := range output.TerminatingInstances {
		terminatedInstances = append(terminatedInstances, *instance.InstanceId)
	}

	msg := fmt.Sprintf("terminated instances %s", strings.Join(terminatedInstances, " , "))
	log.Infof("[%s] - %s", s.Name(), msg)
	logrus.Infof("[%s] Deleted AWS node %s ", cfg.Node.Name, msg)

	return nil
}

func (*DeleteNodeStep) Name() string {
	return DeleteNodeStepName
}

func (*DeleteNodeStep) Depends() []string {
	return nil
}

func (*DeleteNodeStep) Description() string {
	return "Deletes node in aws cluster"
}

func (*DeleteNodeStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
