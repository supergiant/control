package amazon

import (
	"context"
	"io"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	ImportClusterStepName = "import_cluster_aws"
)

type ImportClusterStep struct {
}

func NewImportClusterStep(fn GetEC2Fn) *ImportClusterStep {
	return &ImportClusterStep{
	}
}

func InitImportClusterStep(fn GetEC2Fn) {
	steps.RegisterStep(StepFindAMI, NewImportClusterStep(fn))
}

func (s ImportClusterStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	return nil
}

func (s ImportClusterStep) Name() string {
	return ImportClusterStepName
}

func (s ImportClusterStep) Description() string {
	return ImportClusterStepName
}

func (s ImportClusterStep) Depends() []string {
	return nil
}

func (s ImportClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
