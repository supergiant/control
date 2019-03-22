package azure

import (
	"context"
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestDeleteVMStep(t *testing.T) {
	s := NewDeleteVMStep(NewSDK())

	require.NotNil(t, s.sdk, "sdk shouldn't be nil")

	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, []string{CreateGroupStepName}, s.Depends())
	require.Equal(t, DeleteVMStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Delete virtual machine", s.Description(), "check description")
}

func TestDeleteVMStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		inp         *steps.Config
		a           autorest.Authorizer
		step        DeleteVMStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			step:        DeleteVMStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name:        "nil sdk",
			inp:         &steps.Config{},
			step:        DeleteVMStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "delete vm: error",
			inp:  &steps.Config{},
			a:    &autorest.APIKeyAuthorizer{},
			step: DeleteVMStep{
				sdk: fakeSDK{
					vm: fakeVMClient{
						deleteErr: fakeErr,
					},
				},
			},
			expectedErr: fakeErr,
		},
	} {
		if tc.a != nil && tc.inp != nil {
			// set authorizer
			tc.inp.SetAzureAuthorizer(tc.a)

		}
		err := tc.step.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
