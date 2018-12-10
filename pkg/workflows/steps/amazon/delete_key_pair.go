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
	"time"
)

const DeleteKeyPairStepName = "aws_delete_key_pair"

var (
	deleteKeyPairTimeout      = time.Minute * 1
	deleteKeyPairAttemptCount = 3
)

type DeleteKeyPair struct {
	GetEC2 GetEC2Fn
}

func InitDeleteKeyPair(fn GetEC2Fn) {
	steps.RegisterStep(DeleteKeyPairStepName, &DeleteKeyPair{
		GetEC2: fn,
	})
}

func (s *DeleteKeyPair) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.KeyPairName == "" || cfg.AWSConfig.KeyID == "" {
		logrus.Debugf("Skip deleting empty key pair")
		return nil
	}

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	var (
		deleteErr error
		timeout   = deleteKeyPairTimeout
	)

	for i := 0; i < deleteKeyPairAttemptCount; i++ {
		logrus.Debugf("Delete Key pair %s %s in vpc %s",
			cfg.AWSConfig.KeyPairName, cfg.AWSConfig.KeyID, cfg.AWSConfig.VPCID)
		_, deleteErr = EC2.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(cfg.AWSConfig.KeyPairName),
		})

		if err, ok := deleteErr.(awserr.Error); ok {
			logrus.Debugf("Delete Key pair %s caused %s retry in %v ",
				cfg.AWSConfig.KeyPairName, err.Message(), timeout)
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			break
		}
	}

	return nil
}

func (*DeleteKeyPair) Name() string {
	return DeleteKeyPairStepName
}

func (*DeleteKeyPair) Depends() []string {
	return []string{DeleteKeyPairStepName}
}

func (*DeleteKeyPair) Description() string {
	return "Delete route table from vpc"
}

func (*DeleteKeyPair) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
