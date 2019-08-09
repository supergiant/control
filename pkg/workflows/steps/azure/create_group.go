package azure

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateGroupStepName = "CreateResourceGroup"

type GroupsClientFn func(a autorest.Authorizer, subscriptionID string) GroupsInterface

type CreateGroupStep struct {
	groupsClientFn GroupsClientFn
}

func NewCreateGroupStep() *CreateGroupStep {
	return &CreateGroupStep{
		groupsClientFn: GroupsClientFor,
	}
}

func (s *CreateGroupStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.groupsClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "base client builder")
	}

	groupsClient := s.groupsClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)
	_, err := groupsClient.CreateOrUpdate(ctx, toResourceGroupName(config.Kube.ID, config.Kube.Name), resources.Group{
		Name:     to.StringPtr(toResourceGroupName(config.Kube.ID, config.Kube.Name)),
		Location: to.StringPtr(config.AzureConfig.Location),
	})

	return errors.Wrap(err, "create resource group")
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
