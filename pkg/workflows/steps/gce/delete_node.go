package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"

	"fmt"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteNodeStepName = "gce_delete_node"

type DeleteNodeStep struct {
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func NewDeleteNodeStep() (steps.Step, error) {
	return &DeleteNodeStep{
		getClient: GetClient,
	}, nil
}

func (s *DeleteNodeStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.)
	client, err := s.getClient(ctx, config.GCEConfig.ClientEmail,
		config.GCEConfig.PrivateKey, config.GCEConfig.TokenURI)
	if err != nil {
		return err
	}

	_, serr := client.Instances.Delete(config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone,
		config.Node.Name).Do()

	if serr != nil {
		return errors.Wrap(serr, fmt.Sprintf("GCE delete instance %s", config.Node.Name))
	}

	return nil
}

func (s *DeleteNodeStep) Name() string {
	return DeleteNodeStepName
}

func (s *DeleteNodeStep) Depends() []string {
	return nil
}

func (s *DeleteNodeStep) Description() string {
	return "Google compute engine step for creating instance"
}

func (s *DeleteNodeStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
