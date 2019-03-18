package azure

import (
	"context"
	"io"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const GetAuthorizerName = "GetAuthorizer"

type CreadentialsClientFn func(clientID string, clientSecret string, tenantID string) auth.ClientCredentialsConfig

type GetAuthorizer struct {
	clientCreadsFn CreadentialsClientFn
}

func NewGetAuthorizer() *GetAuthorizer {
	return &GetAuthorizer{
		clientCreadsFn: auth.NewClientCredentialsConfig,
	}
}

func (s *GetAuthorizer) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
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

func (s *GetAuthorizer) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *GetAuthorizer) Name() string {
	return GetAuthorizerName
}

func (s *GetAuthorizer) Depends() []string {
	return nil
}

func (s *GetAuthorizer) Description() string {
	return "Azure: Get authentication token"
}
