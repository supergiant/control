package azuresdk

import (
	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	ErrInvalidAuth = errors.New("azure: auth")
)

type SDK struct {
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	TenantID       string
}

func New(ac steps.AzureConfig) *SDK {
	return &SDK{
		TenantID:       ac.TenantID,
		ClientID:       ac.ClientID,
		ClientSecret:   ac.ClientSecret,
		SubscriptionID: ac.SubscriptionID,
	}
}

func (s *SDK) GetAuthorizer() (autorest.Authorizer, error) {
	creds := auth.NewClientCredentialsConfig(s.ClientID, s.ClientSecret, s.TenantID)

	authorither, err := creds.Authorizer()
	if err != nil {
		return nil, errors.Wrap(ErrInvalidAuth, err.Error())
	}

	return authorither, nil
}

func (s *SDK) GetNetworksClient() (network.VirtualNetworksClient, error) {
	a, err := s.GetAuthorizer()
	if err != nil {
		return network.VirtualNetworksClient{}, err
	}

	bc := network.BaseClient{
		Client: autorest.Client{
			Authorizer: a,
		},
	}

	vnetClient := network.VirtualNetworksClient{
		BaseClient: bc,
	}

	return vnetClient, nil
}
