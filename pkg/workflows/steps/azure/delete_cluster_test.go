package azure

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	errFake = errors.New("fake error")
)

func TestDeleteClusterStep(t *testing.T) {
	s := NewDeleteClusterStep(NewSDK())
	require.NotNil(t, s.groupsClientFn, "groups client shouldn't be nil")

	var nilStringSlice []string
	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, nilStringSlice, s.Depends(), "depends not implemented")
	require.Equal(t, DeleteClusterStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Delete Cluster", s.Description(), "check description")
}

func TestDeleteClusterStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name          string
		inp           *steps.Config
		deleteCluster DeleteClusterStep
		expectedErr   error
	}{
		{
			name:          "nil steps config",
			deleteCluster: DeleteClusterStep{},
			expectedErr:   sgerrors.ErrNilEntity,
		},
		{
			name:          "nil groups client builder",
			inp:           &steps.Config{},
			deleteCluster: DeleteClusterStep{},
			expectedErr:   sgerrors.ErrNilEntity,
		},
	} {
		err := tc.deleteCluster.Run(context.Background(), nil, tc.inp)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
