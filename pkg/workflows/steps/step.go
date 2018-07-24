package steps

import (
	"context"
	"github.com/supergiant/supergiant/pkg/workflows"
)

type Step interface {
	Run(ctx context.Context, config workflows.Config) error
}
