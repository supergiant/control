package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteClusterMachinesStepName = "aws_delete_cluster_machines"

type DeleteClusterMachines struct {
	GetEC2 GetEC2Fn
	getSvc func(steps.AWSConfig) (instanceDeleter, error)
}

func InitDeleteClusterMachines(fn GetEC2Fn) {
	steps.RegisterStep(DeleteClusterMachinesStepName, NewDeleteClusterInstances(fn))
}

func NewDeleteClusterInstances(fn GetEC2Fn) *DeleteClusterMachines {
	return &DeleteClusterMachines{
		getSvc: func(config steps.AWSConfig) (instanceDeleter, error) {
			EC2, err := fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteClusterMachines) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	logrus.Infof("[%s] - deleting cluster %s machines",
		s.Name(), cfg.ClusterName)

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("%s Error getting service %v",
			DeleteClusterMachinesStepName, err)
		return errors.Wrapf(err, "%s error getting service",
			DeleteClusterMachinesStepName)
	}

	describeInstanceOutput, err := svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", clouds.TagClusterID)),
				Values: aws.StringSlice([]string{cfg.ClusterID}),
			},
		},
	})

	if err != nil {
		return errors.Wrap(ErrDeleteCluster, err.Error())
	}

	instanceIDS := make([]string, 0)
	spotRequestIDS := make([]string, 0)

	for _, res := range describeInstanceOutput.Reservations {
		for _, instance := range res.Instances {
			instanceIDS = append(instanceIDS, *instance.InstanceId)

			if instance.SpotInstanceRequestId != nil {
				spotRequestIDS = append(spotRequestIDS, *instance.SpotInstanceRequestId)
			}
		}
	}
	if len(instanceIDS) == 0 {
		logrus.Infof("[%s] - no nodes in k8s cluster %s", s.Name(), cfg.ClusterName)
		return nil
	}

	_, err = svc.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDS),
	})

	if err != nil {
		logrus.Error(ErrDeleteCluster, err.Error())
	}

	_, err = svc.CancelSpotInstanceRequestsWithContext(ctx, &ec2.CancelSpotInstanceRequestsInput{
		SpotInstanceRequestIds: aws.StringSlice(spotRequestIDS),
	})

	if err != nil {
		logrus.Error(ErrDeleteCluster, err.Error())
	}

	log.Infof("[%s] - completed", s.Name())
	logrus.Infof("[%s] Deleted AWS cluster %s",
		s.Name(), cfg.ClusterName)

	return nil
}

func (*DeleteClusterMachines) Name() string {
	return DeleteClusterMachinesStepName
}

func (*DeleteClusterMachines) Depends() []string {
	return nil
}

func (*DeleteClusterMachines) Description() string {
	return "Deletes all nodes in aws cluster"
}

func (*DeleteClusterMachines) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
