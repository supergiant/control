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

type fakeAuthorizer struct {
	err error
}

func (a fakeAuthorizer) Authorizer() (autorest.Authorizer, error) {
	return autorest.NullAuthorizer{}, a.err
}

func TestGetAuthorizerStep(t *testing.T) {
	s := NewGetAuthorizerStepStep()

	require.NotNil(t, s.clientCreadsFn, "client credential func shouldn't be nil")

	var nilStringSlice []string
	require.Equal(t, nil, s.Rollback(context.Background(), nil, nil), "rollback not implemented")
	require.Equal(t, nilStringSlice, s.Depends(), "depends not implemented")
	require.Equal(t, GetAuthorizerStepName, s.Name(), "check step name")
	require.Equal(t, "Azure: Get authentication token", s.Description(), "check description")
}

func TestGetAuthorizerStep_Run(t *testing.T) {
	for _, tc := range []struct {
		name        string
		inp         *steps.Config
		step        GetAuthorizerStep
		expectedErr error
	}{
		{
			name:        "nil steps config",
			step:        GetAuthorizerStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name:        "nil client credentials getter",
			inp:         &steps.Config{},
			step:        GetAuthorizerStep{},
			expectedErr: sgerrors.ErrNilEntity,
		},
		{
			name: "get token: error",
			inp:  &steps.Config{},
			step: GetAuthorizerStep{
				clientCreadsFn: func(clientID string, clientSecret string, tenantID string) Authorizerer {
					return fakeAuthorizer{err: errFake}
				},
			},
			expectedErr: errFake,
		},
		{
			name: "success",
			inp:  &steps.Config{},
			step: GetAuthorizerStep{
				clientCreadsFn: func(clientID string, clientSecret string, tenantID string) Authorizerer {
					return fakeAuthorizer{}
				},
			},
		},
	} {
		err := tc.step.Run(context.Background(), nil, tc.inp)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}
