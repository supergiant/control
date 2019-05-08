package gce

import (
	"context"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/model"
	"io"
	"math/rand"
	"net"
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
				getSubnetwork: func(ctx context.Context, config steps.GCEConfig, subnetworkName string) (*compute.Subnetwork, error) {
					return client.Subnetworks.Get(config.ProjectID, config.Region, subnetworkName).Do()
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
		IPv4Range: "10.0.0.0/16",
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

	config.GCEConfig.NetworkLink = network.SelfLink
	config.GCEConfig.NetworkName = network.Name

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

	for _, az := range azs {
		_, cidrIP, err := net.ParseCIDR(network.IPv4Range)

		if err != nil {
			logrus.Errorf("Error parsing VPC cidr %s",
				network.IPv4Range)
			return errors.Wrapf(err, "Error parsing VPC cidr %s",
				network.IPv4Range)
		}

		logrus.Info(cidrIP)
		subnetCidr, err := cidr.Subnet(cidrIP, 8, rand.Int()%256)
		logrus.Debugf("Subnet cidr %s", subnetCidr)

		if err != nil {
			logrus.Debugf("Calculating subnet cidr caused %s", err.Error())
			return errors.Wrapf(err, "%s Calculating subnet"+
				" cidr caused error", CreateNetworksStepName)
		}

		subnetwork := &compute.Subnetwork{
			Name: fmt.Sprint("subnetwork-%s-%s", az, config.ClusterID),
			IpCidrRange: subnetCidr.String(),
		}

		_, err = svc.insertSubnetwork(ctx, config.GCEConfig, subnetwork)

		if err != nil {
			logrus.Errorf("Create subnet caused error %v", err)
			return errors.Wrap(err, "Create subnet caused")
		}

		subnetwork, err = svc.getSubnetwork(ctx, config.GCEConfig, subnetwork.Name)

		if err != nil {
			logrus.Errorf("Get subnet caused error %v", err)
			return errors.Wrap(err, "Get subnet caused")
		}

		config.GCEConfig.Subnets[az] = subnetwork.SelfLink
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
