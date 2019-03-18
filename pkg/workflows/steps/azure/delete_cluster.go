package azure

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteClusterStepName = "DeleteCluster"

type DeleteClusterStep struct {
	groupsClientFn GroupsClientFn
}

func NewDeleteClusterStep() *DeleteClusterStep {
	return &DeleteClusterStep{
		groupsClientFn: GroupsClientFor,
	}
}

func (s DeleteClusterStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.groupsClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "base client builder")
	}

	groupsClient := s.groupsClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)

	// All cluster resources have been added to the this resources group.
	_, err := groupsClient.Delete(ctx, toResourceGroupName(config.ClusterID, config.ClusterName))

	return errors.Wrap(err, "delete cluster: delete resource group")
}

func (s DeleteClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s DeleteClusterStep) Name() string {
	return DeleteClusterStepName
}

func (s DeleteClusterStep) Depends() []string {
	return nil
}

func (s DeleteClusterStep) Description() string {
	return "Azure: Delete Cluster"
}
