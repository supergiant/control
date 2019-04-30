package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func NSGClientFor(a autorest.Authorizer, subscriptionID string) (SecurityGroupInterface, autorest.Client) {
	nsgClient := network.NewSecurityGroupsClient(subscriptionID)
	nsgClient.Authorizer = a
	return nsgClient, nsgClient.Client
}

func SubnetClientFor(a autorest.Authorizer, subscriptionID string) SubnetGetter {
	subnetClient := network.NewSubnetsClient(subscriptionID)
	subnetClient.Authorizer = a
	return subnetClient
}

func VNetClientFor(a autorest.Authorizer, subscriptionID string) (VirtualNetworkCreator, autorest.Client) {
	vnetClient := network.NewVirtualNetworksClient(subscriptionID)
	vnetClient.Authorizer = a
	return vnetClient, vnetClient.Client
}

func GroupsClientFor(a autorest.Authorizer, subscriptionID string) GroupsInterface {
	gclient := resources.NewGroupsClient(subscriptionID)
	gclient.Authorizer = a
	return gclient
}

// TODO: use sdk here
func ensureAuthorizer(sdk SDKInterface, config *steps.Config) error {
	if sdk == nil || config == nil {
		return errors.Wrap(sgerrors.ErrRawError, "sdk or config is nil")
	}
	if config.GetAzureAuthorizer() != nil {
		return nil
	}

	a, err := auth.NewClientCredentialsConfig(
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

func toResourceGroupName(clusterID, clusterName string) string {
	return fmt.Sprintf("sg-%s-%s", clusterName, clusterID)
}

func toVNetName(clusterID, clusterName string) string {
	return fmt.Sprintf("sg-vnet-%s-%s", clusterName, clusterID)
}

func toNSGName(clusterID, clusterName, role string) string {
	return fmt.Sprintf("sg-nsg-%s-%s-%s", clusterName, clusterID, role)
}

func toSubnetName(clusterID, clusterName, role string) string {
	return fmt.Sprintf("sg-subnet-%s-%s-%s", clusterName, clusterID, role)
}

func toASName(clusterID, clusterName, role string) string {
	return fmt.Sprintf("sg-as-%s-%s-%s", clusterName, clusterID, role)
}

func toLBName(clusterID, clusterName string) string {
	return fmt.Sprintf("sg-lb-%s-%s", clusterName, clusterID)
}

func toIPName(meta string) string {
	return fmt.Sprintf("ip0-%s", meta)
}

func toNICName(vmName string) string {
	return fmt.Sprintf("nic0-%s", vmName)
}
