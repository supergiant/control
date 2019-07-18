package azure

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateVMStepName = "CreateVirtualMachine"

	UbuntuPublisher = "Canonical"
	UbuntuOffer     = "UbuntuServer"
	UbuntuSKU       = "18.04-LTS"

	ifaceName = "ip0"
)

type CreateVMStep struct {
	sdk SDKInterface
}

func NewCreateVMStep(s SDK) *CreateVMStep {
	return &CreateVMStep{
		sdk: s,
	}
}

func (s *CreateVMStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.sdk == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "azure sdk")
	}

	if err := ensureAuthorizer(s.sdk, config); err != nil {
		return errors.Wrap(err, "ensure authorization")
	}

	// TODO: set user with config
	config.Kube.SSHConfig.User = clouds.OSUser

	vmName := util.MakeNodeName(config.ClusterName, config.TaskID, config.IsMaster)

	config.Node = model.Machine{
		Name:     vmName,
		TaskID:   config.TaskID,
		Region:   config.AzureConfig.Location,
		Role:     model.ToRole(config.IsMaster),
		Size:     config.AzureConfig.VMSize,
		Provider: clouds.Azure,
		State:    model.MachineStatePlanned,
	}

	// Update node state in cluster
	config.NodeChan() <- config.Node

	if err := s.setupVM(ctx, config, vmName); err != nil {
		config.Node.State = model.MachineStateError
		config.NodeChan() <- config.Node
		return errors.Wrapf(err, "setup %s vm", vmName)
	}

	config.Node.CreatedAt = time.Now().Unix()
	config.Node.State = model.MachineStateProvisioning

	config.NodeChan() <- config.Node
	if config.IsMaster {
		config.AddMaster(&config.Node)
	} else {
		config.AddNode(&config.Node)
	}

	logrus.Debugf("Machine created %s/%s", config.ClusterName, config.Node.Name)
	return nil
}

func (s *CreateVMStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateVMStep) Name() string {
	return CreateVMStepName
}

func (s *CreateVMStep) Depends() []string {
	return []string{CreateGroupStepName}
}

func (s *CreateVMStep) Description() string {
	return "Azure: Create virtual machine"
}

func (s *CreateVMStep) setupVM(ctx context.Context, config *steps.Config, vmName string) error {
	var lbName string
	if config.IsMaster {
		lbName = toLBName(config.ClusterID, config.ClusterName)
	}

	asName := toASName(config.ClusterID, config.ClusterName, model.ToRole(config.IsMaster).String())
	as, err := s.ensureAvailabilitySet(
		ctx,
		config.GetAzureAuthorizer(),
		config.AzureConfig.SubscriptionID,
		config.AzureConfig.Location,
		toResourceGroupName(config.ClusterID, config.ClusterName),
		asName,
	)
	if err != nil {
		return errors.Wrapf(err, "ensure %s availability set: %s", asName, err)
	}

	nic, err := s.setupNIC(
		ctx,
		config.GetAzureAuthorizer(),
		config.AzureConfig.SubscriptionID,
		config.AzureConfig.Location,
		toResourceGroupName(config.ClusterID, config.ClusterName),
		toVNetName(config.ClusterID, config.ClusterName),
		toSubnetName(config.ClusterID, config.ClusterName, model.ToRole(config.IsMaster).String()),
		toNSGName(config.ClusterID, config.ClusterName, model.ToRole(config.IsMaster).String()),
		toIPName(vmName),
		toNICName(vmName),
		lbName,
	)
	if err != nil {
		config.Node.State = model.MachineStateError
		config.NodeChan() <- config.Node
		return errors.Wrap(err, "setup nic")
	}

	volumeSize, err := strconv.Atoi(config.AzureConfig.VolumeSize)

	if err != nil {
		return errors.Wrapf(err, "convert volume size to int")
	}

	volumeSize32 := int32(volumeSize)
	vmClient := s.sdk.VMClient(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)
	f, err := vmClient.CreateOrUpdate(
		ctx,
		toResourceGroupName(config.ClusterID, config.ClusterName),
		vmName,
		compute.VirtualMachine{
			Location: to.StringPtr(config.AzureConfig.Location),
			VirtualMachineProperties: &compute.VirtualMachineProperties{
				AvailabilitySet: &compute.SubResource{
					ID: as.ID,
				},
				HardwareProfile: &compute.HardwareProfile{
					VMSize: compute.VirtualMachineSizeTypes(config.AzureConfig.VMSize),
				},
				StorageProfile: &compute.StorageProfile{
					ImageReference: &compute.ImageReference{
						Publisher: to.StringPtr(UbuntuPublisher),
						Offer:     to.StringPtr(UbuntuOffer),
						Sku:       to.StringPtr(UbuntuSKU),
						Version:   to.StringPtr("latest"),
					},
					OsDisk: &compute.OSDisk{
						CreateOption: compute.DiskCreateOptionTypesFromImage,
						Caching:      compute.CachingTypesReadWrite,
						OsType:       compute.Linux,
						DiskSizeGB:   &volumeSize32,
						ManagedDisk: &compute.ManagedDiskParameters{
							StorageAccountType: compute.StorageAccountTypesStandardLRS,
						},
					},
				},
				OsProfile: &compute.OSProfile{
					ComputerName:  to.StringPtr(vmName),
					AdminUsername: to.StringPtr(clouds.OSUser),
					LinuxConfiguration: &compute.LinuxConfiguration{
						DisablePasswordAuthentication: to.BoolPtr(true),
						SSH: &compute.SSHConfiguration{
							PublicKeys: toPublicKeys(config.Kube.SSHConfig.BootstrapPublicKey, config.Kube.SSHConfig.PublicKey),
						},
					},
				},
				NetworkProfile: &compute.NetworkProfile{
					NetworkInterfaces: &[]compute.NetworkInterfaceReference{
						{
							ID: nic.ID,
							NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return errors.Wrapf(err, "run %s vm", vmName)
	}

	restclient := s.sdk.RestClient(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)
	if err = f.WaitForCompletionRef(ctx, restclient); err != nil {
		return errors.Wrapf(err, "wait for %s vm is ready", vmName)
	}

	vm, err := vmClient.Get(ctx, toResourceGroupName(config.ClusterID, config.ClusterName), vmName, compute.InstanceView)
	if err != nil {
		return errors.Wrapf(err, "get %s vm", vmName)
	}

	config.Node.PrivateIp = getPrivateIP(nic)
	config.Node.PublicIp, err = s.getPublicIP(
		ctx,
		config.GetAzureAuthorizer(),
		config.AzureConfig.SubscriptionID,
		toResourceGroupName(config.ClusterID, config.ClusterName),
		toIPName(vmName),
	)
	if err != nil {
		return errors.Wrapf(err, "get %s public ip addresses", toIPName(vmName))
	}

	if config.Node.PublicIp == "" || config.Node.PrivateIp == "" {
		return errors.Wrapf(sgerrors.ErrRawError, "failed to get ip addresses for %s vm", vmName)
	}

	config.Node.ID = to.String(vm.ID)
	return nil
}

func (s *CreateVMStep) setupNIC(ctx context.Context, a autorest.Authorizer, subsID, location, groupName,
	vnetName, subnetName, nsgName, ipName, nicName, lbName string) (network.Interface, error) {

	subnet, err := s.sdk.SubnetClient(a, subsID).Get(ctx, groupName, vnetName, subnetName, "")
	if err != nil {
		return network.Interface{}, errors.Wrap(err, "get subnet")
	}

	nsg, err := s.sdk.NSGClient(a, subsID).Get(ctx, groupName, nsgName, "")
	if err != nil {
		return network.Interface{}, errors.Wrap(err, "get network security group")
	}

	ip, err := s.createPublicIP(ctx, a, subsID, location, groupName, ipName)
	if err != nil {
		return network.Interface{}, errors.Wrap(err, "create public ip address")
	}

	var lbPool *[]network.BackendAddressPool
	if len(lbName) > 0 {
		// add to the lb address pool
		lb, err := s.sdk.LBClient(a, subsID).Get(ctx, groupName, lbName, "")
		if err != nil {
			return network.Interface{}, err
		}
		if lb.BackendAddressPools == nil || len(*lb.BackendAddressPools) == 0 {
			return network.Interface{}, errors.Wrapf(sgerrors.ErrRawError, "%s load balancer has no address pools", lbName)
		}
		lbPool = &[]network.BackendAddressPool{
			{
				ID: (*lb.BackendAddressPools)[0].ID,
			},
		}
	}

	nicParams := network.Interface{
		Name:     to.StringPtr(nicName),
		Location: to.StringPtr(location),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			NetworkSecurityGroup: &nsg,
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					Name: to.StringPtr(ifaceName),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						Subnet:                          &subnet,
						PrivateIPAllocationMethod:       network.Dynamic,
						PublicIPAddress:                 &ip,
						LoadBalancerBackendAddressPools: lbPool,
					},
				},
			},
			EnableIPForwarding: to.BoolPtr(true),
		},
	}

	future, err := s.sdk.NICClient(a, subsID).CreateOrUpdate(ctx, groupName, nicName, nicParams)
	if err != nil {
		return network.Interface{}, errors.Wrap(err, "create nic")
	}

	err = future.WaitForCompletionRef(ctx, s.sdk.RestClient(a, subsID))
	if err != nil {
		return network.Interface{}, errors.Wrap(err, "wait for a nic is ready")
	}

	return s.sdk.NICClient(a, subsID).Get(ctx, groupName, nicName, "")
}

func (s *CreateVMStep) createPublicIP(ctx context.Context, a autorest.Authorizer, subsID, location, groupName, ipName string) (network.PublicIPAddress, error) {
	f, err := s.sdk.PublicAddressesClient(a, subsID).CreateOrUpdate(
		ctx,
		groupName,
		ipName,
		network.PublicIPAddress{
			Name:     to.StringPtr(ipName),
			Location: to.StringPtr(location),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion: network.IPv4,
			},
		},
	)
	if err != nil {
		return network.PublicIPAddress{}, err
	}

	err = f.WaitForCompletionRef(ctx, s.sdk.RestClient(a, subsID))
	if err != nil {
		return network.PublicIPAddress{}, errors.Wrap(err, "wait for public ip address is ready")
	}

	return s.sdk.PublicAddressesClient(a, subsID).Get(ctx, groupName, ipName, "")
}

func toPublicKeys(bootstrapPKey, pkey string) *[]compute.SSHPublicKey {
	keys := make([]compute.SSHPublicKey, 0)
	if len(bootstrapPKey) > 0 {
		keys = append(keys, compute.SSHPublicKey{
			Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", clouds.OSUser)),
			KeyData: to.StringPtr(bootstrapPKey),
		})
	}
	if len(pkey) > 0 {
		keys = append(keys, compute.SSHPublicKey{
			Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", clouds.OSUser)),
			KeyData: to.StringPtr(pkey),
		})
	}
	return &keys
}

func (s *CreateVMStep) getPublicIP(ctx context.Context, a autorest.Authorizer, subsID, groupName, ipName string) (string, error) {
	ip, err := s.sdk.PublicAddressesClient(a, subsID).Get(ctx, groupName, ipName, "")
	if err != nil {
		return "", err
	}
	return to.String(ip.IPAddress), nil
}

func (s *CreateVMStep) ensureAvailabilitySet(ctx context.Context, a autorest.Authorizer, subsID, location, groupName, asName string) (compute.AvailabilitySet, error) {
	as, err := s.sdk.AvailabilitySetClient(a, subsID).CreateOrUpdate(ctx, groupName, asName, compute.AvailabilitySet{
		Location: to.StringPtr(location),
		Sku: &compute.Sku{
			Name: to.StringPtr(string(compute.Aligned)),
		},
		AvailabilitySetProperties: &compute.AvailabilitySetProperties{
			// TODO: set it dynamically? https://docs.microsoft.com/en-us/azure/virtual-machines/windows/manage-availability#number-of-fault-domains-per-region
			PlatformFaultDomainCount: to.Int32Ptr(2),
		},
	})
	if err != nil {
		return compute.AvailabilitySet{}, err
	}
	return as, err
}

func getPrivateIP(nic network.Interface) string {
	for _, iface := range *nic.IPConfigurations {
		if to.String(iface.Name) != ifaceName {
			continue
		}
		return to.String(iface.PrivateIPAddress)
	}
	return ""
}
