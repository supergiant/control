package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

type GroupsInterface interface {
	CreateOrUpdate(ctx context.Context, name string, grp resources.Group) (resources.Group, error)
	Delete(ctx context.Context, name string) (resources.GroupsDeleteFuture, error)
}

type SecurityGroupInterface interface {
	CreateOrUpdate(ctx context.Context, groupName string, nsgName string, params network.SecurityGroup) (network.SecurityGroupsCreateOrUpdateFuture, error)
	Get(ctx context.Context, groupName string, nsgName string, expand string) (network.SecurityGroup, error)
}

type VirtualNetworkCreator interface {
	CreateOrUpdate(ctx context.Context, groupName string, vnetName string, params network.VirtualNetwork) (network.VirtualNetworksCreateOrUpdateFuture, error)
}

type NICInterface interface {
	CreateOrUpdate(ctx context.Context, groupName string, nicName string, params network.Interface) (network.InterfacesCreateOrUpdateFuture, error)
	Get(ctx context.Context, groupName string, nicName string, expand string) (network.Interface, error)
}

type PublicAddressInterface interface {
	CreateOrUpdate(ctx context.Context, groupName string, ipName string, params network.PublicIPAddress) (network.PublicIPAddressesCreateOrUpdateFuture, error)
	Get(ctx context.Context, groupName string, ipName string, expand string) (network.PublicIPAddress, error)
}

type SubnetGetter interface {
	Get(ctx context.Context, groupName, vnetName, subnetName, expand string) (network.Subnet, error)
}

type VMInterface interface {
	CreateOrUpdate(ctx context.Context, groupName string, vmName string, params compute.VirtualMachine) (compute.VirtualMachinesCreateOrUpdateFuture, error)
	Get(ctx context.Context, groupName string, vmName string, expand compute.InstanceViewTypes) (compute.VirtualMachine, error)
	Delete(ctx context.Context, groupName string, vmName string) (compute.VirtualMachinesDeleteFuture, error)
}

type SDKInterface interface {
	RestClient(a autorest.Authorizer, subscriptionID string) autorest.Client
	PublicAddressesClient(a autorest.Authorizer, subscriptionID string) PublicAddressInterface
	NICClient(a autorest.Authorizer, subscriptionID string) NICInterface
	SubnetClient(a autorest.Authorizer, subscriptionID string) SubnetGetter
	NSGClient(a autorest.Authorizer, subscriptionID string) SecurityGroupInterface
	VMClient(a autorest.Authorizer, subscriptionID string) VMInterface
	LBClient(a autorest.Authorizer, subscriptionID string) network.LoadBalancersClient
}

type SDK struct {
}

func NewSDK() SDK {
	return SDK{}
}

func (s SDK) RestClient(a autorest.Authorizer, subscriptionID string) autorest.Client {
	gclient := resources.NewGroupsClient(subscriptionID)
	gclient.Authorizer = a
	return gclient.Client
}

func (s SDK) NICClient(a autorest.Authorizer, subscriptionID string) NICInterface {
	nicClient := network.NewInterfacesClient(subscriptionID)
	nicClient.Authorizer = a
	return nicClient
}

func (s SDK) PublicAddressesClient(a autorest.Authorizer, subscriptionID string) PublicAddressInterface {
	ipClient := network.NewPublicIPAddressesClient(subscriptionID)
	ipClient.Authorizer = a
	return ipClient
}

func (s SDK) SubnetClient(a autorest.Authorizer, subscriptionID string) SubnetGetter {
	subnetClient := network.NewSubnetsClient(subscriptionID)
	subnetClient.Authorizer = a
	return subnetClient
}

func (s SDK) NSGClient(a autorest.Authorizer, subscriptionID string) SecurityGroupInterface {
	nsgClient := network.NewSecurityGroupsClient(subscriptionID)
	nsgClient.Authorizer = a
	return nsgClient
}

func (s SDK) VMClient(a autorest.Authorizer, subscriptionID string) VMInterface {
	vmclient := compute.NewVirtualMachinesClient(subscriptionID)
	vmclient.Authorizer = a
	return vmclient
}

func (s SDK) LBClient(a autorest.Authorizer, subscriptionID string) network.LoadBalancersClient {
	lbclient := network.NewLoadBalancersClient(subscriptionID)
	lbclient.Authorizer = a
	return lbclient
}
