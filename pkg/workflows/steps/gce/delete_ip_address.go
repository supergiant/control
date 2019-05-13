package gce

import (
	"context"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteIpAddressStepName = "gce_delete_ip_address"

type DeleteIpAddressStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteIpAddressStep() *DeleteIpAddressStep {
	return &DeleteIpAddressStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)
			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteIpAddress: func(ctx context.Context, config steps.GCEConfig, addressName string) (*compute.Operation, error) {
					return client.Addresses.Delete(config.ServiceAccount.ProjectID, config.Region, addressName).Do()
				},
			}, nil
		},
	}
}

func (s *DeleteIpAddressStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteIpAddressStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteIpAddressStepName)
	}

	_, err = svc.deleteIpAddress(ctx, config.GCEConfig, config.GCEConfig.ExternalAddressName)

	if err != nil {
		logrus.Errorf("Error deleting external address %s %v", config.GCEConfig.ExternalAddressName, err)
	}

	_, err = svc.deleteIpAddress(ctx, config.GCEConfig, config.GCEConfig.ExternalAddressName)

	if err != nil {
		logrus.Errorf("Error deleting internal address %s %v", config.GCEConfig.InternalAddressName, err)
	}

	return nil
}

func (s *DeleteIpAddressStep) Name() string {
	return DeleteIpAddressStepName
}

func (s *DeleteIpAddressStep) Depends() []string {
	return nil
}

func (s *DeleteIpAddressStep) Description() string {
	return "Delete ip addresses"
}

func (s *DeleteIpAddressStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
