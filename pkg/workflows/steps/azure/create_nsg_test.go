package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	fakeErr = errors.New("fake")
)

type fakeNSGClient struct {
	err error
}

func (c fakeNSGClient) CreateOrUpdate(ctx context.Context, groupName string, nsgName string, params network.SecurityGroup) (network.SecurityGroupsCreateOrUpdateFuture, error) {
	return network.SecurityGroupsCreateOrUpdateFuture{}, c.err
}

func TestCreateSecurityGroupStep(t *testing.T) {
	s := NewCreateSecurityGroupStep()
	require.NotNil(t, s.nsgClientFn, "nsg client shouldn't be nil")
	require.NotNil(t, s.findOutboundIP, "findOutboundIP func shouldn't be nil")

	var nilStringSlice []string
	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, nilStringSlice, s.Depends(), "depends not implemented")
	require.Equal(t, CreateSecurityGroupStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Create master/node network security groups", s.Description(), "check description")
}

func TestCreateSecurityGroupStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		step        CreateSecurityGroupStep
		expectedErr error
	}{
		{
			name:        "nsg client builder is empty",
			step:        CreateSecurityGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "build nsg client: error",
			step: CreateSecurityGroupStep{
				nsgClientFn: func(authorizer Autorizerer, subscriptionID string) (SecurityGroupCreator, error) {
					return nil, fakeErr
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "get sg address: error",
			step: CreateSecurityGroupStep{
				nsgClientFn: func(authorizer Autorizerer, subscriptionID string) (SecurityGroupCreator, error) {
					return fakeNSGClient{}, nil
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "", fakeErr
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "create nsg: error",
			step: CreateSecurityGroupStep{
				nsgClientFn: func(authorizer Autorizerer, subscriptionID string) (SecurityGroupCreator, error) {
					return fakeNSGClient{
						err: fakeErr,
					}, nil
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "1.2.3.4", nil
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "create nsg: success",
			step: CreateSecurityGroupStep{
				nsgClientFn: func(authorizer Autorizerer, subscriptionID string) (SecurityGroupCreator, error) {
					return fakeNSGClient{}, nil
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "1.2.3.4", nil
				},
			},
		},
	} {
		err := tc.step.Run(context.Background(), nil, &steps.Config{})

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
