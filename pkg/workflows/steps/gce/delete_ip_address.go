package gce

import (
	"context"
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
			client, err := GetClient(ctx, config)

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
		logrus.Errorf("Error deleting external address %v", err)
	}

	_, err = svc.deleteIpAddress(ctx, config.GCEConfig, config.GCEConfig.ExternalAddressName)

	if err != nil {
		logrus.Errorf("Error deleting internal address %v", err)
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
