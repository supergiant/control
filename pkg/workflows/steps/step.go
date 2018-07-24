package steps

import (
	"context"
)

type Step interface {
	Run(ctx context.Context, config Config) error
}
