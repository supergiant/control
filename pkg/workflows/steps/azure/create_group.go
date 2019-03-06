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

type BaseClientFn func(authorizer Autorizerer, subscriptionID string) (resources.BaseClient, error)

type CreateGroupStep struct {
	baseClientFn BaseClientFn
}

func NewCreateInstanceStep() *CreateGroupStep {
	return &CreateGroupStep{
		baseClientFn: BaseClientFor,
	}
}

func (s *CreateGroupStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if s.baseClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "base client builder")
	}

	baseClient, err := s.baseClientFn(auth.ClientCredentialsConfig{
		ClientID:     config.AzureConfig.ClientID,
		ClientSecret: config.AzureConfig.ClientSecret,
		TenantID:     config.AzureConfig.TenantID,
	},
		config.AzureConfig.SubscriptionID,
	)
	if err != nil {
		return err
	}

	groupsClient := resources.GroupsClient{BaseClient: baseClient}
	_, err = groupsClient.CreateOrUpdate(ctx, "", resources.Group{
		Name:     toStrPtr(toResourceGroupName(config.ClusterID, config.ClusterName)),
		Location: toStrPtr(config.AzureConfig.Location),
		Tags:     map[string]*string{},
	})

	return errors.Wrap(err, "create resources group")
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
