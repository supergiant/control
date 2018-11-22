package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type GetEC2Fn func(steps.AWSConfig) (ec2iface.EC2API, error)

func GetEC2(cfg steps.AWSConfig) (ec2iface.EC2API, error) {
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
