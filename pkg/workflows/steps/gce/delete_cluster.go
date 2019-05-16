package gce

import (
	"context"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/model"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteClusterStepName = "gce_delete_cluster"

type DeleteClusterStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteClusterStep() *DeleteClusterStep {
	return &DeleteClusterStep{
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

func (s *DeleteClusterStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.
	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		return errors.Wrapf(err, "%s get service", DeleteClusterStepName)
	}

	for _, master := range config.GetMasters() {
		if master.State == model.MachineStatePlanned || master.State == model.MachineStateBuilding {
			continue
		}

		logrus.Debugf("Delete master %s in %s", master.Name, master.Region)

		_, serr := svc.deleteInstance(config.GCEConfig.ServiceAccount.ProjectID,
			master.Region,
			master.Name)

		if serr != nil {
			return errors.Wrap(serr, "GCE delete instance")
		}
	}

	for _, node := range config.GetNodes() {
		if node.State == model.MachineStatePlanned || node.State == model.MachineStateBuilding {
			continue
		}

		logrus.Debugf("Delete node %s in %s", node.Name, node.Region)
		_, serr := svc.deleteInstance(config.GCEConfig.ServiceAccount.ProjectID,
			node.Region,
			node.Name)

		if serr != nil {
			return errors.Wrap(serr, "GCE delete instance")
		}
	}

	return nil
}

func (s *DeleteClusterStep) Name() string {
	return DeleteClusterStepName
}

func (s *DeleteClusterStep) Depends() []string {
	return nil
}

func (s *DeleteClusterStep) Description() string {
	return "Google compute engine delete cluster step"
}

func (s *DeleteClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
