package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteInternetGatewayStepName = "aws_delete_internet_gateway"

type DeleteInternetGateway struct {
	GetEC2 GetEC2Fn
}

func InitDeleteInternetGateWay(fn GetEC2Fn) {
	steps.RegisterStep(DeleteInternetGatewayStepName, &DeleteInternetGateway{
		GetEC2: fn,
	})
}

func (s *DeleteInternetGateway) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Detach internet gateway %s from vpc %s",
		cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)
	_, err = EC2.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
		VpcId:             aws.String(cfg.AWSConfig.VPCID),
	})

	if err != nil {
		logrus.Debugf("Detach internet gateway %s from vpc %s caused %v",
			cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID, err)
	}

	logrus.Debugf("Delete internet gateway %s from vpc %s",
		cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)

	_, err = EC2.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
	})

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("DisassociateRouteTable caused %s", err.Message())
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
