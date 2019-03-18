package azure

import (
	"context"
	"io"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const GetAuthorizerStepName = "GetAuthorizer"

type Authorizerer interface {
	Authorizer() (autorest.Authorizer, error)
}

type CreadentialsClientFn func(clientID string, clientSecret string, tenantID string) Authorizerer

type GetAuthorizerStep struct {
	clientCreadsFn CreadentialsClientFn
}

func NewGetAuthorizerStepStep() *GetAuthorizerStep {
	return &GetAuthorizerStep{
		clientCreadsFn: func(clientID string, clientSecret string, tenantID string) Authorizerer {
			a := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
			return &a
		},
	}
}

func (s *GetAuthorizerStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.clientCreadsFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "base client builder")
	}

	a, err := s.clientCreadsFn(
		config.AzureConfig.ClientID,
		config.AzureConfig.ClientSecret,
		config.AzureConfig.TenantID,
	).Authorizer()
	if err != nil {
		return err
	}

	config.SetAzureAuthorizer(a)
	return nil
}

func (s *GetAuthorizerStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *GetAuthorizerStep) Name() string {
	return GetAuthorizerStepName
}

func (s *GetAuthorizerStep) Depends() []string {
	return nil
}

func (s *GetAuthorizerStep) Description() string {
	return "Azure: Get authentication token"
}
