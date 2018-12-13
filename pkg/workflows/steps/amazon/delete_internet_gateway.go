package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteInternetGatewayStepName = "aws_delete_internet_gateway"

type IGWDeleter interface {
	DetachInternetGateway(*ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error)
	DeleteInternetGateway(*ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error)
}

type DeleteInternetGateway struct {
	getIGWService func(steps.AWSConfig) (IGWDeleter, error)
}

func InitDeleteInternetGateWay(fn GetEC2Fn) {
	steps.RegisterStep(DeleteInternetGatewayStepName,
		NewDeleteInernetGateway(fn))
}

func NewDeleteInernetGateway(fn GetEC2Fn) *DeleteInternetGateway {
	return &DeleteInternetGateway{
		getIGWService: func(config steps.AWSConfig) (IGWDeleter, error) {
			EC2, err := fn(config)
			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteInternetGateway) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.InternetGatewayID == "" {
		logrus.Debug("Skip deleting empty Internet GW")
		return nil
	}

	svc, err := s.getIGWService(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("Error while getting IGW deleter %v", err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Detach internet gateway %s from vpc %s",
		cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)
	_, err = svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
		VpcId:             aws.String(cfg.AWSConfig.VPCID),
	})

	if err != nil {
		logrus.Debugf("Detach internet gateway %s from vpc %s caused %v",
			cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID, err)
		return errors.Wrapf(err, "Detach internet gateway")
	}

	logrus.Debugf("Delete internet gateway %s from vpc %s",
		cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)

	_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
	})

	if err != nil {
		logrus.Debugf("DeleteInternetGateway caused %s", err.Error())
		return errors.Wrapf(err, "DeleteInternetGateway")
	}

	logrus.Debugf("Internet gateway %s was deleted from vpc %s",
		cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)
	return nil
}

func (*DeleteInternetGateway) Name() string {
	return DeleteInternetGatewayStepName
}

func (*DeleteInternetGateway) Depends() []string {
	return []string{DeleteSecurityGroupsStepName}
}

func (*DeleteInternetGateway) Description() string {
	return "Delete internet gateway from VPC"
}

func (*DeleteInternetGateway) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
