package amazon

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type FakeEC2KeyPair struct {
	ec2iface.EC2API

	importOutput   *ec2.ImportKeyPairOutput
	describeOutput *ec2.DescribeKeyPairsOutput
	importErr      error
	describeErr    error
}

func (f *FakeEC2KeyPair) ImportKeyPairWithContext(aws.Context, *ec2.ImportKeyPairInput, ...request.Option) (*ec2.ImportKeyPairOutput, error) {
	return f.importOutput, f.importErr
}

func (f *FakeEC2KeyPair) DescribeKeyPairs(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return f.describeOutput, f.describeErr
}

func TestKeyPairStep_Run(t *testing.T) {
	cfg := steps.NewConfig("TEST", "", "myacc", profile.Profile{})
	cfg.AWSConfig.KeyPairName = "mypair"
	fingerprint := "e9:b0:fe:0d:1a:4c:f9:00:dd:fd:c2:16:05:8b:3f:83"
	step := NewImportKeyPairStep(func(config steps.AWSConfig) (ec2iface.EC2API, error) {
		return &FakeEC2KeyPair{
			importOutput: &ec2.ImportKeyPairOutput{
				KeyName:        aws.String(cfg.AWSConfig.KeyPairName),
				KeyFingerprint: &fingerprint,
			},
			importErr:      nil,
			describeOutput: &ec2.DescribeKeyPairsOutput{},
			describeErr:    nil,
		}, nil
	})

	err := step.Run(context.Background(), os.Stdout, cfg)
	require.NoError(t, err)
}
