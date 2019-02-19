package azure

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/clouds/azuresdk"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateGroupStepName = "CreateResourceGroup"

type CreateGroupStep struct {
}

func NewCreateInstanceStep() *CreateGroupStep {
	return &CreateGroupStep{}
}

func (s *CreateGroupStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	sdk := azuresdk.New(cfg.AzureConfig)

	groupName := toResourceGroupName(cfg.ClusterID, cfg.ClusterName)

	groupsClient, err := sdk.GroupsClient()
	if err != nil {
		return err
	}
	result, err := groupsClient.CreateOrUpdate(ctx, groupName, resources.Group{
		Name:     toStrPtr(groupName),
		Location: toStrPtr(cfg.AzureConfig.Location),
		Tags:     map[string]*string{},
	})

	if err != nil {
		return errors.Wrap(err, "create resources group")
	}
	fmt.Printf("%+v", result)

	cfg.AzureConfig.ResourceGroupName = groupName

	return nil
}

func (s *CreateGroupStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateGroupStep) Name() string {
	return CreateGroupStepName
}

func (s *CreateGroupStep) Depends() []string {
	return nil
}

func (s *CreateGroupStep) Description() string {
	return "Azure: Create ResourceGroup"
}
