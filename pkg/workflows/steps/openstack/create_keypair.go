package openstack

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateKeyPairStepName = "create_keypair"

type CreateKeyPairStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateKeyPairStep() *CreateKeyPairStep {
	return &CreateKeyPairStep{
		getClient: func(config steps.OpenStackConfig) (client *gophercloud.ProviderClient, e error) {
			opts := gophercloud.AuthOptions{
				IdentityEndpoint: config.AuthURL,
				Username:         config.UserName,
				Password:         config.Password,
				TenantName:       config.TenantName,
				DomainID:         config.DomainID,
			}

			client, err := openstack.AuthenticatedClient(opts)
			if err != nil {
				return nil, err
			}

			return client, nil
		},
	}
}

func (s *CreateKeyPairStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateKeyPairStepName)
	}

	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", CreateKeyPairStepName)
	}

	keypair, err := keypairs.Create(computeClient, keypairs.CreateOpts{
		Name:         fmt.Sprintf("keypair-%s", config.Kube.ID),
		PublicKey: config.Kube.SSHConfig.BootstrapPublicKey,
	}).Extract()

	if err != nil {
		return errors.Wrapf(err, "create keypair error step %s", CreateKeyPairStepName)
	}

	config.OpenStackConfig.KeyPairName = keypair.Name

	return nil
}

func (s *CreateKeyPairStep) Name() string {
	return CreateKeyPairStepName
}

func (s *CreateKeyPairStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateKeyPairStep) Description() string {
	return "Create key pair"
}

func (s *CreateKeyPairStep) Depends() []string {
	return nil
}
