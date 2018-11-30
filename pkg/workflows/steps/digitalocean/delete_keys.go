package digitalocean

import (
	"context"
	"io"
	"time"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type DeleteKeysStep struct {
	keyService KeyService
	timeout          time.Duration
}

func NewDeleteKeysStep(keyService KeyService) *DeleteKeysStep {
	return &DeleteKeysStep{
		keyService: keyService,
	}
}

func (s *DeleteKeysStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	return nil
}

func (s *DeleteKeysStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteKeysStep) Name() string {
	return DeleteDeleteKeysStep
}

func (s *DeleteKeysStep) Depends() []string {
	return nil
}

func (s *DeleteKeysStep) Description() string {
	return ""
}
