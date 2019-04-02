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

func TestCreateGroupStep(t *testing.T) {
	s := NewCreateGroupStep()

	require.NotNil(t, s.groupsClientFn, "base client shouldn't be nil")

	var nilStringSlice []string
	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, nilStringSlice, s.Depends(), "depends not implemented")
	require.Equal(t, CreateGroupStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Create ResourceGroup", s.Description(), "check description")
}

func TestCreateGroupStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		inp         *steps.Config
		createGroup CreateGroupStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			createGroup: CreateGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name:        "nil groups client builder",
			inp:         &steps.Config{},
			createGroup: CreateGroupStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "delete cluster error",
			inp:  &steps.Config{},
			createGroup: CreateGroupStep{
				groupsClientFn: func(a autorest.Authorizer, subscriptionID string) GroupsInterface {
					return fakeGroupsClient{createErr: errFake}
				},
			},
			expectedErr: errFake,
		},
		{
			name: "success",
			inp:  &steps.Config{},
			createGroup: CreateGroupStep{
				groupsClientFn: func(a autorest.Authorizer, subscriptionID string) GroupsInterface {
					return fakeGroupsClient{}
				},
			},
		},
	} {
		err := tc.createGroup.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
