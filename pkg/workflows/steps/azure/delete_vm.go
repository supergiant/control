package azure

import (
	"context"
	"io"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteVMStepName = "DeleteVirtualMachine"

type DeleteVMStep struct {
	sdk SDKInterface
}

func NewDeleteVMStep(s SDK) *DeleteVMStep {
	return &DeleteVMStep{
		sdk: s,
	}
}

func (s *DeleteVMStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.sdk == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "azure sdk")
	}

	if err := ensureAuthorizer(s.sdk, config); err != nil {
		return errors.Wrap(err, "ensure authorization")
	}

	f, err := s.sdk.VMClient(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID).Delete(
		ctx,
		toResourceGroupName(config.Kube.ID, config.Kube.Name),
		config.Node.Name,
	)
	if err != nil {
		return errors.Wrapf(err, "delete %s vm", config.Node.Name)
	}

	restclient := s.sdk.RestClient(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)
	if err = f.WaitForCompletionRef(ctx, restclient); err != nil {
		return errors.Wrapf(err, "wait for %s vm is ready", config.Node.Name)
	}

	log.Debugf("cluster %s: %s machine has been deleted", config.Kube.Name, config.Node.Name)
	return nil
}

func (s *DeleteVMStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteVMStep) Name() string {
	return DeleteVMStepName
}

func (s *DeleteVMStep) Depends() []string {
	return []string{CreateGroupStepName}
}

func (s *DeleteVMStep) Description() string {
	return "Azure: Delete virtual machine"
}
