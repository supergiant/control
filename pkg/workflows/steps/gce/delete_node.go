package gce

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteNodeStepName = "gce_delete_node"

type DeleteNodeStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteNodeStep() *DeleteNodeStep {
	return &DeleteNodeStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteInstance: func(projectID string, region string, name string) (*compute.Operation, error) {
					return client.Instances.Delete(projectID, region, name).Do()
				},
			}, nil
		},
	}
}

func (s *DeleteNodeStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client
	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		return errors.Wrapf(err, "%s get service", DeleteClusterStepName)
	}

	logrus.Debugf("Delete node %s in %s",
		config.Node.Name, config.Node.Region)
	_, serr := svc.deleteInstance(config.GCEConfig.ServiceAccount.ProjectID,
		config.Node.Region,
		config.Node.Name)

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
	return "Google compute engine delete instance step"
}

func (s *DeleteNodeStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
