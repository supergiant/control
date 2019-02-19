package azure

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-10-01/network"
	"github.com/supergiant/control/pkg/clouds/azuresdk"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateVNetStepName = "create_vnet"

type CreateVnetStep struct {
}

func (s *CreateVnetStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	sdk := azuresdk.New(cfg.AzureConfig)

	cl, err := sdk.NetworksClient()
	if err != nil {
		return err
	}

	if cfg.AzureConfig.VirtualNetworkName == "" {
		vnetName := fmt.Sprintf("sg-%s-%s", cfg.ClusterID, cfg.ClusterName)
		log.Infof("[%s] - trying to create virtual network %s", CreateGroupStepName, vnetName)

		res, err := cl.CreateOrUpdate(ctx, cfg.AzureConfig.ResourceGroupName, vnetName, network.VirtualNetwork{
			Location: toStrPtr(cfg.AzureConfig.Location),
			Name:     toStrPtr(vnetName),
		})
		if err != nil {
			return err
		}
		done, err := res.DoneWithContext(ctx, cl.Sender)
		if err != nil {
			return err
		}
		if !done {
			log.Errorf("[%s] - failed to create virtual network", s.Name())
			return errors.New("failed to create vnet")
		}

		cfg.AzureConfig.VirtualNetworkName = vnetName
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
