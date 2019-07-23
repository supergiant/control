package azure

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateVNetAndSubnetsStepName = "CreateVirtualNetworkNetworkAndSubnets"

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

	_, vnet, err := net.ParseCIDR(config.AzureConfig.VNetCIDR)
	if err != nil {
		return errors.Wrap(err, "parse vnet cidr")
	}
	mastersSubnet, err := buildSubnet(vnet, 1, toSubnetName(config.Kube.ID, config.Kube.Name, model.RoleMaster.String()))
	if err != nil {
		return errors.Wrap(err, "calculate subnet for masters")
	}
	nodesSubnet, err := buildSubnet(vnet, 5, toSubnetName(config.Kube.ID, config.Kube.Name, model.RoleNode.String()))
	if err != nil {
		return errors.Wrap(err, "calculate subnet for nodes")
	}

	vnetClient, restclient := s.vnetClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)

	f, err := vnetClient.CreateOrUpdate(
		ctx,
		toResourceGroupName(config.Kube.ID, config.Kube.Name),
		toVNetName(config.Kube.ID, config.Kube.Name),
		network.VirtualNetwork{
			Location: to.StringPtr(config.AzureConfig.Location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{config.AzureConfig.VNetCIDR},
				},
				Subnets: &[]network.Subnet{
					mastersSubnet,
					nodesSubnet,
				},
			},
		},
	)
	if err != nil {
		return errors.Wrapf(err, "create %s vnet", toVNetName(config.Kube.ID, config.Kube.Name))
	}

	ctx, _ = context.WithTimeout(ctx, 60*VNetCreationTimeout)
	err = f.WaitForCompletionRef(ctx, restclient)
	return errors.Wrap(err, "wait for vnet is ready")
}

func (s *CreateVirtualNetworkStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateVirtualNetworkStep) Name() string {
	return CreateVNetAndSubnetsStepName
}

func (s *CreateVirtualNetworkStep) Depends() []string {
	return []string{CreateGroupStepName}
}

func (s *CreateVirtualNetworkStep) Description() string {
	return "Azure: Create virtual network"
}

// TODO: this doesn't calculate a subcidr just updates the net number and mask (aws uses the same)
func buildSubnet(baseCIDR *net.IPNet, netNum int, name string) (network.Subnet, error) {
	subnetCidr, err := cidr.Subnet(baseCIDR, 8, netNum)
	if err != nil {
		return network.Subnet{}, errors.Wrapf(sgerrors.ErrRawError, "%s", err)
	}
	return network.Subnet{
		Name: to.StringPtr(name),
		SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr(subnetCidr.String()),
		},
	}, nil
}
