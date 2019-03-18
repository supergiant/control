package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type fakeVNetClient struct {
	res network.VirtualNetworksCreateOrUpdateFuture
	err error
}

func (c fakeVNetClient) CreateOrUpdate(ctx context.Context, groupName string, vnetName string, params network.VirtualNetwork) (network.VirtualNetworksCreateOrUpdateFuture, error) {
	return network.VirtualNetworksCreateOrUpdateFuture{}, c.err
}

func TestCreateVirtualNetworkStep(t *testing.T) {
	s := NewCreateVirtualNetworkStep()
	require.NotNil(t, s.vnetClientFn, "nsg client shouldn't be nil")

	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, []string{CreateGroupStepName}, s.Depends(), "depends not implemented")
	require.Equal(t, CreateVirtualNetworkStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Create virtual network", s.Description(), "check description")
}

func TestCreateVirtualNetworkStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		inp         *steps.Config
		step        CreateVirtualNetworkStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			step:        CreateVirtualNetworkStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name:        "vnet client builder is empty",
			inp:         &steps.Config{},
			step:        CreateVirtualNetworkStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "create vnet: error",
			inp:  &steps.Config{},
			step: CreateVirtualNetworkStep{
				vnetClientFn: func(a autorest.Authorizer, subscriptionID string) (VirtualNetworkCreator, autorest.Client) {
					return fakeVNetClient{
						err: fakeErr,
					}, autorest.Client{}
				},
			},
			expectedErr: fakeErr,
		},
	} {
		err := tc.step.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
