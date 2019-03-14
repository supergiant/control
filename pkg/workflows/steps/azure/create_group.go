package azure

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateGroupStepName = "CreateResourceGroup"

type GroupsClientFn func(authorizer Autorizerer, subscriptionID string) (GroupsInterface, error)

type CreateGroupStep struct {
	groupsClientFn GroupsClientFn
}

func NewCreateGroupStep() *CreateGroupStep {
	return &CreateGroupStep{
		groupsClientFn: GroupsClientFor,
	}
}

func (s *CreateGroupStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if s.groupsClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "base client builder")
	}

	c, err := s.groupsClientFn(auth.ClientCredentialsConfig{
		ClientID:     config.AzureConfig.ClientID,
		ClientSecret: config.AzureConfig.ClientSecret,
		TenantID:     config.AzureConfig.TenantID,
	},
		config.AzureConfig.SubscriptionID,
	)
	if err != nil {
		return err
	}

	_, err = c.CreateOrUpdate(ctx, "", resources.Group{
		Name:     toStrPtr(toResourceGroupName(config.ClusterID, config.ClusterName)),
		Location: toStrPtr(config.AzureConfig.Location),
		Tags:     map[string]*string{},
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
