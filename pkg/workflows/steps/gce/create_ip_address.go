package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateIPAddressStepName = "gce_create_ip_address"

type CreateAddressStep struct {
	Timeout       time.Duration
	AttemptCount  int
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateAddressStep() (*CreateAddressStep, error) {
	return &CreateAddressStep{
		Timeout:      time.Second * 10,
		AttemptCount: 6,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config.ClientEmail,
				config.PrivateKey, config.TokenURI)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertAddress: func(ctx context.Context, config steps.GCEConfig, address *compute.Address) (*compute.Operation, error) {
					return client.Addresses.Insert(config.ProjectID, config.Region, address).Do()
				},
				getAddress: func(ctx context.Context, config steps.GCEConfig, addressName string) (*compute.Address, error) {
					return client.Addresses.Get(config.ProjectID, config.Region, addressName).Do()
				},
			}, nil
		},
	}, nil
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
	config.GCEConfig.ExternalAddressName = externalAddressName

	externalAddress := &compute.Address{
		Name:        externalAddressName,
		Description: "External static IP address",
		AddressType: "EXTERNAL",
		IpVersion:   "IPV4",
	}

	_, err = svc.insertAddress(ctx, config.GCEConfig, externalAddress)

	if err != nil {
		logrus.Errorf("Error creating external ip address %v", err)
		return errors.Wrapf(err, "error creating external ip address types")
	}

	timeout := s.Timeout

	for i := 0; i < s.AttemptCount; i++ {
		externalAddress, err = svc.getAddress(ctx, config.GCEConfig, externalAddressName)

		if err == nil && externalAddress.Status == "RESERVED" {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error obtaining external ip address %v", err)
		return errors.Wrapf(err, "error obtaining external ip address types")
	}

	config.GCEConfig.ExternalIPAddress = externalAddress.Address
	config.ExternalDNSName = externalAddress.Address
	internalAddressName := fmt.Sprintf("in-ip-%s", config.ClusterID)
	config.GCEConfig.InternalAddressName = internalAddressName

	internalAddress := &compute.Address{
		Name:        internalAddressName,
		Description: "Internal static IP address",
		AddressType: "INTERNAL",
		IpVersion:   "IPV4",
	}
	_, err = svc.insertAddress(ctx, config.GCEConfig, internalAddress)

	if err != nil {
		logrus.Errorf("Error creating internal ip address %v", err)
		return errors.Wrapf(err, "error creating internal ip address types")
	}

	timeout = s.Timeout

	for i := 0; i < s.AttemptCount; i++ {
		internalAddress, err = svc.getAddress(ctx, config.GCEConfig, internalAddressName)

		if err == nil && internalAddress.Status == "RESERVED" {
			break
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error obtaining internal ip address %v", err)
		return errors.Wrapf(err, "error obtaining internal ip address types")
	}

	config.GCEConfig.InternalAddressName = internalAddress.Address
	config.InternalDNSName = internalAddress.Address

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
