package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
<<<<<<< HEAD
=======
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
>>>>>>> e6855504... WIP
)

type GetEC2Fn func(steps.AWSConfig) (ec2iface.EC2API, error)

func GetEC2(cfg steps.AWSConfig) (ec2iface.EC2API, error) {
	logrus.Debug("get EC2 client")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(cfg.KeyID, cfg.Secret, ""),
		},
	})

	if err != nil {
		return nil, err
	}
	return ec2.New(sess), nil
}

type GetIAMFn func(steps.AWSConfig) (iamiface.IAMAPI, error)

func GetIAM(cfg steps.AWSConfig) (iamiface.IAMAPI, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(cfg.KeyID, cfg.Secret, ""),
		},
	})

	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}

// NOT(stgleb): *elb.ELB doesn't implement elbiface.ELBAPI for some reasom
type GetELBFn func(steps.AWSConfig) (*elb.ELB, error)

func GetELB(cfg steps.AWSConfig) (elbiface.ELBAPI, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(cfg.KeyID, cfg.Secret, ""),
		},
	})

	if err != nil {
		return nil, err
	}
	return elb.New(sess), nil
}
