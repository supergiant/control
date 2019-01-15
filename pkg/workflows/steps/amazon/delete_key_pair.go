package amazon

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteKeyPairStepName = "aws_delete_key_pair"

type KeyService interface {
	DeleteKeyPair(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error)
}

var (
	deleteKeyPairTimeout      = time.Minute * 1
	deleteKeyPairAttemptCount = 3
)

type DeleteKeyPair struct {
	getSvc func(steps.AWSConfig) (KeyService, error)
	GetEC2 GetEC2Fn
}

func InitDeleteKeyPair(fn GetEC2Fn) {
	steps.RegisterStep(DeleteKeyPairStepName, NewDeleteKeyPairStep(fn))
}

func NewDeleteKeyPairStep(fn GetEC2Fn) *DeleteKeyPair {
	return &DeleteKeyPair{
		getSvc: func(config steps.AWSConfig) (KeyService, error) {
			EC2, err := fn(config)
			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
		GetEC2: fn,
	}
}

func (s *DeleteKeyPair) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.KeyPairName == "" || cfg.AWSConfig.KeyID == "" {
		logrus.Debugf("Skip deleting empty key pair")
		return nil
	}

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Error getting EC2 key service %v", err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	var (
		deleteErr error
		timeout   = deleteKeyPairTimeout
	)

	for i := 0; i < deleteKeyPairAttemptCount; i++ {
		logrus.Debugf("Delete Key pair %s %s in vpc %s",
			cfg.AWSConfig.KeyPairName, cfg.AWSConfig.KeyID, cfg.AWSConfig.VPCID)
		_, deleteErr = svc.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(cfg.AWSConfig.KeyPairName),
		})

		if deleteErr != nil {
			logrus.Debugf("Delete Key pair %s caused %s retry in %v ",
				cfg.AWSConfig.KeyPairName, deleteErr.Error(), timeout)
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
	return nil
}

func (*DeleteKeyPair) Description() string {
	return "Delete key pair"
}

func (*DeleteKeyPair) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
