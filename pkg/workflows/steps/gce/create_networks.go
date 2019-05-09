package gce

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/model"
	"io"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateNetworksStepName = "gce_create_networks"

type CreateNetworksStep struct {
	Timeout       time.Duration
	AttemptCount  int
	accountGetter     accountGetter

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
	zoneGetterFactory func(context.Context, accountGetter, *steps.Config) (account.ZonesGetter, error)
}

type accountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

func NewCreateNetworksStep(getter accountGetter) *CreateNetworksStep {
	return &CreateNetworksStep{
		Timeout:      time.Second * 10,
		AttemptCount: 10,
		accountGetter: getter,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertAddress: func(ctx context.Context, config steps.GCEConfig, address *compute.Address) (*compute.Operation, error) {
					return client.Addresses.Insert(config.ServiceAccount.ProjectID, config.Region, address).Do()
				},
				getAddress: func(ctx context.Context, config steps.GCEConfig, addressName string) (*compute.Address, error) {
					return client.Addresses.Get(config.ServiceAccount.ProjectID, config.Region, addressName).Do()
				},
				getNetwork: func(ctx context.Context, config steps.GCEConfig,networkName string) (*compute.Network, error) {
					return client.Networks.Get(config.ProjectID, networkName).Do()
				},
				insertNetwork: func(ctx context.Context, config steps.GCEConfig, network *compute.Network) (*compute.Operation, error) {
					return client.Networks.Insert(config.ProjectID, network).Do()
				},
			}, nil
		},
		zoneGetterFactory: func(ctx context.Context, accountGetter accountGetter,
			cfg *steps.Config) (account.ZonesGetter, error) {
			acc, err := accountGetter.Get(ctx, cfg.CloudAccountName)

			if err != nil {
				logrus.Errorf("Get cloud account %s caused error %v",
					cfg.CloudAccountName, err)
				return nil, errors.Wrapf(err, "Get cloud account")
			}

			zoneGetter, err := account.NewZonesGetter(acc, cfg)

			return zoneGetter, err
		},
	}
}

func (s *CreateNetworksStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateNetworksStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateIPAddressStepName)
	}

	network := &compute.Network{
		AutoCreateSubnetworks: true,
		Name: fmt.Sprintf("network-%s", config.ClusterID),
	}

	_, err = svc.insertNetwork(ctx, config.GCEConfig, network)

	if err != nil {
		logrus.Errorf("Create network caused error %v", err)
		return errors.Wrap(err, "Create network caused error")
	}

	network, err = svc.getNetwork(ctx, config.GCEConfig, network.Name)

	if err != nil {
		logrus.Errorf("Get network caused error %v", err)
		return errors.Wrap(err, "Get network caused error")
	}

	// TODO(stgleb): store subnetworks from network
	config.GCEConfig.Subnets = network.Subnetworks
	config.GCEConfig.NetworkLink = network.SelfLink
	config.GCEConfig.NetworkName = network.Name

	logrus.Debugf("Created network name %s link %s",
		config.GCEConfig.NetworkName, config.GCEConfig.NetworkLink)
	zoneGetter, err := s.zoneGetterFactory(ctx, s.accountGetter, config)

	if err != nil {
		logrus.Errorf("Create zone getter caused error %v", err)
		return errors.Wrap(err, "Create zone getter caused")
	}

	azs, err := zoneGetter.GetZones(ctx, *config)

	if err != nil {
		logrus.Errorf("get availability zones %v", err)
		return errors.Wrap(err, "get availability zones")
	}

	config.GCEConfig.Subnets = make(map[string]string)

	// Use subnets as a storage of availability zones
	for _, az := range azs {
		config.GCEConfig.Subnets[az] = "dummy"
	}

	return nil
}

func (s *CreateNetworksStep) Name() string {
	return CreateNetworksStepName
}

func (s *CreateNetworksStep) Depends() []string {
	return nil
}

func (s *CreateNetworksStep) Description() string {
	return "Create network and subnetworks"
}

func (s *CreateNetworksStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
