package azure

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestCreateGroupStep(t *testing.T) {
	s := NewCreateGroupStep()

	require.NotNil(t, s.groupsClientFn, "base client shouldn't be nil")

	var nilStringSlice []string
	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, nilStringSlice, s.Depends(), "depends not implemented")
	require.Equal(t, CreateGroupStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Create ResourceGroup", s.Description(), "check description")
}

func TestCreateClusterStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		createGroup CreateGroupStep
		expectedErr error
	}{
		{
			name:        "nil groups client builder",
			createGroup: CreateGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "build groups client error",
			createGroup: CreateGroupStep{
				groupsClientFn: func(authorizer Autorizerer, subscriptionID string) (GroupsInterface, error) {
					return nil, errFake
				},
			},
			expectedErr: errFake,
		},
		{
			name: "delete cluster error",
			createGroup: CreateGroupStep{
				groupsClientFn: func(authorizer Autorizerer, subscriptionID string) (GroupsInterface, error) {
					return fakeGroupsClient{createErr: errFake}, nil
				},
			},
			expectedErr: errFake,
		},
		{
			name: "success",
			createGroup: CreateGroupStep{
				groupsClientFn: func(authorizer Autorizerer, subscriptionID string) (GroupsInterface, error) {
					return fakeGroupsClient{}, nil
				},
			},
		},
	} {
		err := tc.createGroup.Run(context.Background(), nil, &steps.Config{})

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
