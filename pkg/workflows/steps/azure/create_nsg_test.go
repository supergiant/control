package azure

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	fakeErr = errors.New("fake")
)

type fakeNSGClient struct {
	res network.SecurityGroupsCreateOrUpdateFuture
	err error
}

func (c fakeNSGClient) CreateOrUpdate(ctx context.Context, groupName string, nsgName string, params network.SecurityGroup) (network.SecurityGroupsCreateOrUpdateFuture, error) {
	return network.SecurityGroupsCreateOrUpdateFuture{}, c.err
}

type fakeSubnetGetter struct {
	res network.Subnet
	err error
}

func (f fakeSubnetGetter) Get(ctx context.Context, groupName, vnetName, subnetName, expand string) (network.Subnet, error) {
	return f.res, f.err
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
		inp         *steps.Config
		step        CreateSecurityGroupStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			step:        CreateSecurityGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name:        "nsg client builder is nil",
			inp:         &steps.Config{},
			step:        CreateSecurityGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "subnet getter is nil",
			inp:  &steps.Config{},
			step: CreateSecurityGroupStep{
				nsgClientFn: func(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
					return fakeNSGClient{}, autorest.Client{}
				},
			},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "get sg address: error",
			inp:  &steps.Config{},
			step: CreateSecurityGroupStep{
				nsgClientFn: func(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
					return fakeNSGClient{}, autorest.Client{}
				},
				subnetGetterFn: func(a autorest.Authorizer, subscriptionID string) SubnetGetter {
					return fakeSubnetGetter{}
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "", fakeErr
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "get subnet: error",
			inp:  &steps.Config{},
			step: CreateSecurityGroupStep{
				nsgClientFn: func(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
					return fakeNSGClient{}, autorest.Client{}
				},
				subnetGetterFn: func(a autorest.Authorizer, subscriptionID string) SubnetGetter {
					return fakeSubnetGetter{err: fakeErr}
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "", nil
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "create nsg: error",
			inp:  &steps.Config{},
			step: CreateSecurityGroupStep{
				nsgClientFn: func(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
					return fakeNSGClient{
						err: fakeErr,
					}, autorest.Client{}
				},
				subnetGetterFn: func(a autorest.Authorizer, subscriptionID string) SubnetGetter {
					return fakeSubnetGetter{}
				},
				findOutboundIP: func(ctx context.Context) (string, error) {
					return "1.2.3.4", nil
				},
			},
			expectedErr: fakeErr,
		},
		// TODO: mock azure.Future
		//{
		//	name: "create nsg: success",
		//	step: CreateSecurityGroupStep{
		//		nsgClientFn: func(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
		//			return fakeNSGClient{
		//				res: network.SecurityGroupsCreateOrUpdateFuture{
		//					Future: toAzureFuture(http.Response{
		//						Request: &http.Request{
		//							Method: http.MethodPost,
		//						},
		//					}),
		//				},
		//			}, autorest.Client{}
		//		},
		//		findOutboundIP: func(ctx context.Context) (string, error) {
		//			return "1.2.3.4", nil
		//		},
		//	},
		//},
	} {
		err := tc.step.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}

func toAzureFuture(r http.Response) azure.Future {
	f, err := azure.NewFutureFromResponse(&r)
	if err != nil {
		return azure.Future{}
	}
	return f
}
