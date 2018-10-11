package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"os"
	"testing"
)

type FakeEC2KeyPair struct {
	ec2iface.EC2API

	output *ec2.ImportKeyPairOutput
	err    error
}

func (f *FakeEC2KeyPair) ImportKeyPairWithContext(aws.Context, *ec2.ImportKeyPairInput, ...request.Option) (*ec2.ImportKeyPairOutput, error) {
	return f.output, f.err
}

func TestKeyPairStep_Run(t *testing.T) {
	cfg := steps.NewConfig("TEST", "", "myacc", profile.Profile{})
	cfg.AWSConfig.KeyPairName = "mypair"

	step := NewImportKeyPairStep(func(config steps.AWSConfig) (ec2iface.EC2API, error) {
		return &FakeEC2KeyPair{
			output: &ec2.ImportKeyPairOutput{
				KeyName: aws.String(cfg.AWSConfig.KeyPairName),
			},
			err: nil,
		}, nil
	})

	err := step.Run(context.Background(), os.Stdout, cfg)
	require.NoError(t, err)
}
