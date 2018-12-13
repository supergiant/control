package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteNodeStepName = "aws_delete_node"

type instanceDeleter interface {
	DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error)
	TerminateInstancesWithContext(aws.Context, *ec2.TerminateInstancesInput, ...request.Option) (*ec2.TerminateInstancesOutput, error)
}

type DeleteNodeStep struct {
	getSvc func(steps.AWSConfig) (instanceDeleter, error)
}

func InitDeleteNode(fn GetEC2Fn) {
	steps.RegisterStep(DeleteNodeStepName, NewDeleteVPC(fn))
}

func NewDeleteNode(fn GetEC2Fn) *DeleteNodeStep {
	return &DeleteNodeStep{
		getSvc: func(cfg steps.AWSConfig) (instanceDeleter, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteNodeStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	logrus.Infof("[%s] - deleting node %s", s.Name(), cfg.Node.Name)

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Get instance by name filter %s", cfg.Node.Name)
	describeInstanceOutput, err := svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
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
	_, err = svc.TerminateInstancesWithContext(ctx,
		&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice(instanceIDS),
		})

	if err != nil {
		return errors.Wrapf(err, "%s terminate instance", DeleteNodeStepName)
	}

	log.Infof("[%s] - finished successfully", s.Name())

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
