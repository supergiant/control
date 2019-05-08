package gce

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateIPAddressStepName = "gce_create_ip_address"

type CreateAddressStep struct {
	Timeout       time.Duration
	AttemptCount  int
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateAddressStep() *CreateAddressStep {
	return &CreateAddressStep{
		Timeout:      time.Second * 10,
		AttemptCount: 10,
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
			}, nil
		},
	}
}

func (s *CreateAddressStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateIPAddressStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateIPAddressStepName)
	}

	externalAddressName := fmt.Sprintf("ex-ip-%s", config.ClusterID)
	logrus.Debugf("create external ip address name %s", externalAddressName)

	externalAddress := &compute.Address{
		Name:        externalAddressName,
		Description: "External static IP address",
		AddressType: "EXTERNAL",
	}

	_, err = svc.insertAddress(ctx, config.GCEConfig, externalAddress)

	if err != nil {
		logrus.Errorf("Error creating external ip address %v", err)
		return errors.Wrapf(err, "error creating external ip address types")
	}

	timeout := s.Timeout

	for i := 0; i < s.AttemptCount; i++ {
		externalAddress, err = svc.getAddress(ctx, config.GCEConfig, externalAddressName)

		if err == nil && externalAddress.Address != "" {
			config.GCEConfig.ExternalIPAddressLink = externalAddress.SelfLink
			config.ExternalDNSName = externalAddress.Address
			config.InternalDNSName = externalAddress.Address
			logrus.Debugf("External IP %s", externalAddress.Address)
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error obtaining external ip address %v", err)
		return errors.Wrapf(err, "error obtaining external ip address types")
	}

	logrus.Debugf("Save external IP address SelfLink %s", externalAddress.SelfLink)
	internalAddressName := fmt.Sprintf("in-ip-%s", config.ClusterID)

	internalAddress := &compute.Address{
		Name:        internalAddressName,
		Description: "Internal static IP address",
		AddressType: "INTERNAL",
	}

	logrus.Debugf("create internal ip address %s", internalAddressName)
	_, err = svc.insertAddress(ctx, config.GCEConfig, internalAddress)

	if err != nil {
		logrus.Errorf("Error creating internal ip address %v", err)
		return errors.Wrapf(err, "error creating internal ip address types")
	}

	timeout = s.Timeout

	for i := 0; i < s.AttemptCount; i++ {
		internalAddress, err = svc.getAddress(ctx, config.GCEConfig, internalAddressName)

		if err == nil && internalAddress.Address != "" {
			config.GCEConfig.InternalIPAddressLink = internalAddress.SelfLink
			config.InternalDNSName = internalAddress.Address
			logrus.Debugf("Internal IP %s", internalAddress.Address)
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error obtaining internal ip address %v", err)
		return errors.Wrapf(err, "error obtaining internal ip address types")
	}

	logrus.Debugf("Save internal IP address SelfLink %s", internalAddress.SelfLink)


	return nil
}

func (s *CreateAddressStep) Name() string {
	return CreateIPAddressStepName
}

func (s *CreateAddressStep) Depends() []string {
	return nil
}

func (s *CreateAddressStep) Description() string {
	return "Create static ip addresses"
}

func (s *CreateAddressStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
