package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"os"
	"testing"
	"time"
)

func TestCreateVPCStep_Run(t *testing.T) {
	tt := []struct {
		awsFN GetEC2Fn
		err   error
	}{
		{
			func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return nil, nil
			},
			nil,
		},
	}

	for i, tc := range tt {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		cfg := steps.NewConfig("TEST", "", "TEST", profile.Profile{
			Region:   "us-east-1",
			Provider: clouds.AWS,
		})

		step := NewCreateVPCStep(tc.awsFN)
		err := step.Run(ctx, os.Stdout, cfg)
		if tc.err != nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.EqualError(t, tc.err, err.Error(), "TC%d, %v", i, err)
		}
	}

}
