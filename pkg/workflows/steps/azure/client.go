package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

var (
	_ GroupsInterface       = resources.GroupsClient{}
	_ VirtualNetworkCreator = network.VirtualNetworksClient{}
	_ SecurityGroupCreator  = network.SecurityGroupsClient{}
)

type GroupsInterface interface {
	CreateOrUpdate(ctx context.Context, name string, grp resources.Group) (resources.Group, error)
	Delete(ctx context.Context, name string) (resources.GroupsDeleteFuture, error)
}

type SecurityGroupCreator interface {
	CreateOrUpdate(ctx context.Context, groupName string, nsgName string, params network.SecurityGroup) (network.SecurityGroupsCreateOrUpdateFuture, error)
}

type SubnetGetter interface {
	Get(ctx context.Context, groupName, vnetName, subnetName, expand string) (network.Subnet, error)
}

type VirtualNetworkCreator interface {
	CreateOrUpdate(ctx context.Context, groupName string, vnetName string, params network.VirtualNetwork) (network.VirtualNetworksCreateOrUpdateFuture, error)
}

func NSGClientFor(a autorest.Authorizer, subscriptionID string) (SecurityGroupCreator, autorest.Client) {
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
