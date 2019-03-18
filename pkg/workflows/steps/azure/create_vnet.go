package azure

import (
	"context"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateVirtualNetworkStepName = "CreateVirtualNetworkNetwork"

var VNetCreationTimeout = 60 * time.Second

type VNetClientFn func(a autorest.Authorizer, subscriptionID string) (VirtualNetworkCreator, autorest.Client)

type CreateVirtualNetworkStep struct {
	vnetClientFn VNetClientFn
}

func NewCreateVirtualNetworkStep() *CreateVirtualNetworkStep {
	return &CreateVirtualNetworkStep{
		vnetClientFn: VNetClientFor,
	}
}

func (s *CreateVirtualNetworkStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.vnetClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "virtual network client builder")
	}

	vnetClient, restclient := s.vnetClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)

	f, err := vnetClient.CreateOrUpdate(
		ctx,
		toResourceGroupName(config.ClusterID, config.ClusterName),
		toVNetName(config.ClusterID, config.ClusterName),
		network.VirtualNetwork{
			Location: to.StringPtr(config.AzureConfig.Location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{config.AzureConfig.VNetCIDR},
				},
			},
		},
	)
	if err != nil {
		return errors.Wrapf(err, "create %s vnet", toVNetName(config.ClusterID, config.ClusterName))
	}

	ctx, _ = context.WithTimeout(ctx, 60*VNetCreationTimeout)
	err = f.WaitForCompletionRef(ctx, restclient)
	return errors.Wrap(err, "wait for vnet is ready")
}

func (s *CreateVirtualNetworkStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateVirtualNetworkStep) Name() string {
	return CreateVirtualNetworkStepName
}

func (s *CreateVirtualNetworkStep) Depends() []string {
	return []string{CreateGroupStepName}
}

func (s *CreateVirtualNetworkStep) Description() string {
	return "Azure: Create virtual network"
}
