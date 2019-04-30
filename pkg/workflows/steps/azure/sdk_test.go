package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/stretchr/testify/require"
)

var (
	_ SDKInterface = fakeSDK{}
)

type fakeSDK struct {
	rest   autorest.Client
	ip     fakeIPClient
	nic    fakeNICClient
	subnet fakeSubnetGetter
	vm     fakeVMClient
}

func (s fakeSDK) AvailabilitySetClient(a autorest.Authorizer, subscriptionID string) compute.AvailabilitySetsClient {
	panic("implement me")
}

func (s fakeSDK) LBClient(a autorest.Authorizer, subscriptionID string) network.LoadBalancersClient {
	panic("implement me")
}

func (s fakeSDK) NSGClient(a autorest.Authorizer, subscriptionID string) SecurityGroupInterface {
	panic("implement me")
}

func (s fakeSDK) RestClient(a autorest.Authorizer, subscriptionID string) autorest.Client {
	return s.rest
}
func (s fakeSDK) PublicAddressesClient(a autorest.Authorizer, subscriptionID string) PublicAddressInterface {
	return s.ip
}
func (s fakeSDK) NICClient(a autorest.Authorizer, subscriptionID string) NICInterface {
	return s.nic
}
func (s fakeSDK) SubnetClient(a autorest.Authorizer, subscriptionID string) SubnetGetter {
	return s.subnet
}
func (s fakeSDK) VMClient(a autorest.Authorizer, subscriptionID string) VMInterface {
	return s.vm
}

type fakeIPClient struct {
}

func (fakeIPClient) CreateOrUpdate(ctx context.Context, groupName string, ipName string, params network.PublicIPAddress) (network.PublicIPAddressesCreateOrUpdateFuture, error) {
	panic("implement me")
}
func (fakeIPClient) Get(ctx context.Context, groupName string, ipName string, expand string) (network.PublicIPAddress, error) {
	panic("implement me")
}

type fakeNICClient struct {
}

func (fakeNICClient) CreateOrUpdate(ctx context.Context, groupName string, nicName string, params network.Interface) (network.InterfacesCreateOrUpdateFuture, error) {
	panic("implement me")
}
func (fakeNICClient) Get(ctx context.Context, groupName string, nicName string, expand string) (network.Interface, error) {
	panic("implement me")
}

type fakeVMClient struct {
	deleteRes compute.VirtualMachinesDeleteFuture
	deleteErr error
}

func (f fakeVMClient) Get(ctx context.Context, groupName string, vmName string, expand compute.InstanceViewTypes) (compute.VirtualMachine, error) {
	panic("implement me")
}

func (f fakeVMClient) Delete(ctx context.Context, groupName string, vmName string) (compute.VirtualMachinesDeleteFuture, error) {
	return f.deleteRes, f.deleteErr
}

func (f fakeVMClient) CreateOrUpdate(ctx context.Context, groupName string, vmName string, params compute.VirtualMachine) (compute.VirtualMachinesCreateOrUpdateFuture, error) {
	panic("implement me")
}

func TestSDK(t *testing.T) {
	sdk := NewSDK()

	restclient := sdk.RestClient(autorest.NullAuthorizer{}, "id")
	require.NotNil(t, restclient.Authorizer)

	nicclient := sdk.NICClient(autorest.NullAuthorizer{}, "id")
	require.NotNil(t, nicclient)

	ipclient := sdk.PublicAddressesClient(autorest.NullAuthorizer{}, "id")
	require.NotNil(t, ipclient)

	subnetclient := sdk.SubnetClient(autorest.NullAuthorizer{}, "id")
	require.NotNil(t, subnetclient)

	vmclient := sdk.VMClient(autorest.NullAuthorizer{}, "id")
	require.NotNil(t, vmclient)
}
