package azuresdk

import (
	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
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

func (s *SDK) Authorizer() (autorest.Authorizer, error) {
	creds := auth.NewClientCredentialsConfig(s.ClientID, s.ClientSecret, s.TenantID)

	authorither, err := creds.Authorizer()
	if err != nil {
		return nil, errors.Wrap(ErrInvalidAuth, err.Error())
	}

	return authorither, nil
}

func (s *SDK) NetworksClient() (network.VirtualNetworksClient, error) {
	a, err := s.Authorizer()
	if err != nil {
		return network.VirtualNetworksClient{}, err
	}

	vnetClient := network.NewVirtualNetworksClient(s.SubscriptionID)
	vnetClient.Authorizer = a

	return vnetClient, nil
}

func (s *SDK) VMClient() (compute.VirtualMachinesClient, error) {
	a, err := s.Authorizer()
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	bc := compute.BaseClient{
		SubscriptionID: s.SubscriptionID,
		BaseURI:        compute.DefaultBaseURI,
		Client: autorest.Client{
			Authorizer: a,
		},
	}

	computeClient := compute.VirtualMachinesClient{
		BaseClient: bc,
	}

	return computeClient, nil
}

func (s *SDK) NetworkInterfaceClient() (network.InterfacesClient, error) {
	a, err := s.Authorizer()
	if err != nil {
		return network.InterfacesClient{}, err
	}

	bc := network.BaseClient{
		SubscriptionID: s.SubscriptionID,
		Client: autorest.Client{
			Authorizer: a,
		},
		BaseURI: network.DefaultBaseURI,
	}

	networkClient := network.InterfacesClient{
		BaseClient: bc,
	}

	return networkClient, nil
}

func (s *SDK) GroupsClient() (resources.GroupsClient, error) {
	a, err := s.Authorizer()
	if err != nil {
		return resources.GroupsClient{}, err
	}

	bc := resources.BaseClient{
		SubscriptionID: s.SubscriptionID,
		Client: autorest.Client{
			Authorizer: a,
		},
		BaseURI: resources.DefaultBaseURI,
	}

	return resources.GroupsClient{
		BaseClient: bc,
	}, nil
}
