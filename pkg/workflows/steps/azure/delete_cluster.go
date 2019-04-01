package azure

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteClusterStepName = "DeleteCluster"

type DeleteClusterStep struct {
	sdk            SDK
	groupsClientFn GroupsClientFn
}

func NewDeleteClusterStep(s SDK) *DeleteClusterStep {
	return &DeleteClusterStep{
		sdk:            s,
		groupsClientFn: GroupsClientFor,
	}
}

func (s DeleteClusterStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.groupsClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "groups client builder")
	}

	if err := ensureAuthorizer(s.sdk, config); err != nil {
		return errors.Wrap(err, "ensure authorization")
	}

	groupsClient := s.groupsClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)

	// All cluster resources have been added to the this resources group.
	name := toResourceGroupName(config.Kube.ID, config.Kube.Name)
	logrus.Debugf("deleting %s azure resource group", name)
	f, err := groupsClient.Delete(ctx, name)
	if err != nil {
		return errors.Wrap(err, "delete cluster: delete resource group")
	}

	// TODO: deletion stacks here
	if err = f.WaitForCompletionRef(ctx, s.sdk.RestClient(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)); err != nil {
		return errors.Wrapf(err, "delete %s resource group", name)
	}

	logrus.Debugf("%s azure resource group has been deleted", name)
	return nil
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
