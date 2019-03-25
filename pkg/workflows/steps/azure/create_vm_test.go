package azure

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestCreateVMStep(t *testing.T) {
	s := NewCreateVMStep(NewSDK())

	require.NotNil(t, s.sdk, "sdk shouldn't be nil")

	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, []string{CreateGroupStepName}, s.Depends())
	require.Equal(t, CreateVMStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Create virtual machine", s.Description(), "check description")
}

func TestCreateVMStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		inp         *steps.Config
		step        CreateVMStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			step:        CreateVMStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
	} {
		err := tc.step.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
