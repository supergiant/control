package azure

import (
	"context"
	"io"

	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/network/mgmt/network"
	"github.com/supergiant/control/pkg/clouds/azuresdk"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateVNetStepName = "create_vnet"

type CreateVnetStep struct {
}

func (*CreateVnetStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	sdk := azuresdk.New(cfg.AzureConfig)

	cl, err := sdk.GetNetworksClient()
	if err != nil {
		return err
	}

	if cfg.AzureConfig.VirtualNetworkName == "" {
		vnetName := fmt.Sprintf("sg-%s-%s", cfg.ClusterID, cfg.ClusterName)
		log.Infof("[%s] - trying to create virtual network %s", CreateGroupStepName, vnetName)

		_, err := cl.CreateOrUpdate(ctx, cfg.AzureConfig.ResourceGroupName, vnetName, network.VirtualNetwork{})
		if err != nil {
			return nil
		}
	} else {
		log.Infof("[%s] - using virtual network %s", cfg.AzureConfig.VirtualNetworkName)
	}

	return nil
}

func (*CreateVnetStep) Name() string {
	return CreateGroupStepName
}

func (*CreateVnetStep) Description() string {
	return ""
}

func (*CreateVnetStep) Depends() []string {
	return nil
}

func (*CreateVnetStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
