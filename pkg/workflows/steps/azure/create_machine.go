package azure

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-10-01/network"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/clouds/azuresdk"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateMachineStepName = "create_machine_azure"

type CreateMachineStep struct {
}

func (*CreateMachineStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	sdk := azuresdk.New(cfg.AzureConfig)

	vms, err := sdk.VMClient()
	if err != nil {
		return err
	}

	nics, err := sdk.NetworkInterfaceClient()
	if err != nil {
		return err
	}

	role := model.RoleMaster
	if !cfg.IsMaster {
		role = model.RoleNode
	}

	vmName := util.MakeNodeName(cfg.ClusterName, cfg.TaskID, cfg.IsMaster)
	nicName := "ipConfig1"

	cfg.Node = model.Machine{
		Name:     vmName,
		TaskID:   cfg.TaskID,
		Region:   cfg.AzureConfig.Location,
		Role:     role,
		Size:     cfg.AzureConfig.Size,
		Provider: clouds.Azure,
		State:    model.MachineStatePlanned,
	}

	nicFuture, err := nics.CreateOrUpdate(ctx, cfg.AzureConfig.ResourceGroupName, vmName, network.Interface{
		Name:     toStrPtr(nicName),
		Location: toStrPtr(cfg.AzureConfig.Location),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{},
			Primary:          toBoolPtr(true),
		},
	})

	if err != nil {
		return err
	}

	interfce, err := nicFuture.Result(nics)
	if err != nil {
		return err
	}

	future, err := vms.CreateOrUpdate(ctx, cfg.AzureConfig.ResourceGroupName, vmName, compute.VirtualMachine{
		Location: toStrPtr(cfg.AzureConfig.Location),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(cfg.AzureConfig.Size),
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  toStrPtr(vmName),
				AdminUsername: toStrPtr(cfg.AzureConfig.User),
				AdminPassword: toStrPtr(cfg.AzureConfig.Password),
				LinuxConfiguration: &compute.LinuxConfiguration{
					DisablePasswordAuthentication: toBoolPtr(true),
					SSH: &compute.SSHConfiguration{
						PublicKeys: &[]compute.SSHPublicKey{
							{
								Path:    toStrPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", cfg.AzureConfig.User)),
								KeyData: toStrPtr(cfg.Kube.SSHConfig.PublicKey),
							},
						},
					},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: interfce.ID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: toBoolPtr(true),
						},
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	vm, err := future.Result(vms)
	if err != nil {
		return err
	}

	cfg.Node.ID = *vm.ID
	cfg.Node.CreatedAt = time.Now().Unix()
	cfg.Node.State = model.MachineStateProvisioning

	for _, ip := range *interfce.IPConfigurations {
		cfg.Node.PublicIp = *ip.PublicIPAddress.IPAddress
		cfg.Node.PrivateIp = *ip.PrivateIPAddress
	}

	cfg.NodeChan() <- cfg.Node

	//TODO add node
	return nil
}

func (*CreateMachineStep) Name() string {
	return CreateMachineStepName
}

func (*CreateMachineStep) Description() string {
	return ""
}

func (*CreateMachineStep) Depends() []string {
	return nil
}

func (*CreateMachineStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
