package openstack

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateSecurityGroupStepName = "create_security_group"

type CreateSecurityGroupStep struct {
	getClient func(steps.OpenStackConfig) (*gophercloud.ProviderClient, error)
}

func NewCreateSecurityGroup() *CreateSecurityGroupStep {
	return &CreateSecurityGroupStep{
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

func (s *CreateSecurityGroupStep) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	client, err := s.getClient(config.OpenStackConfig)

	if err != nil {
		return errors.Wrapf(err, "step %s", CreateSecurityGroupStepName)
	}

	identityClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: config.OpenStackConfig.Region,
	})

	if err != nil {
		return errors.Wrapf(err, "step %s get compute client", CreateSecurityGroupStepName)
	}

	masterSg, err := secgroups.Create(identityClient, secgroups.CreateOpts{
		Name:        fmt.Sprintf("master-sg-%s", config.Kube.ID),
		Description: "master nodes security group",
	}).Extract()

	if err != nil {
		return errors.Wrapf(err, "create master sg error step %s", CreateSecurityGroupStepName)
	}

	// TODO(stgleb): probably we need to save rule ID for deletion
	_, err = secgroups.CreateRule(identityClient, secgroups.CreateRuleOpts{
		ParentGroupID: masterSg.ID,
		FromPort:      22,
		ToPort:        22,
		IPProtocol:    "tcp",
		CIDR:          "0.0.0.0/0",
	}).Extract()

	if err != nil {
		return errors.Wrap(err, "create ssh rule for master security group")
	}

	// Allow traffic through API Server port to master nodes
	_, err = secgroups.CreateRule(identityClient, secgroups.CreateRuleOpts{
		ParentGroupID: masterSg.ID,
		FromPort:      config.Kube.APIServerPort,
		ToPort:        config.Kube.APIServerPort,
		IPProtocol:    "tcp",
		CIDR:          "0.0.0.0/0",
	}).Extract()

	if err != nil {
		return errors.Wrap(err, "create ssh rule for master security group")
	}

	config.OpenStackConfig.MasterSecurityGroupId = masterSg.ID

	workerSg, err := secgroups.Create(identityClient, secgroups.CreateOpts{
		Name:        fmt.Sprintf("worker-sg-%s", config.Kube.ID),
		Description: "worker security group",
	}).Extract()

	if err != nil {
		return errors.Wrapf(err, "create worker sg error step %s", CreateSecurityGroupStepName)
	}

	// TODO(stgleb): probably we need to save rule ID for deletion
	_, err = secgroups.CreateRule(identityClient, secgroups.CreateRuleOpts{
		ParentGroupID: workerSg.ID,
		FromPort:      22,
		ToPort:        22,
		IPProtocol:    "tcp",
		CIDR:          "0.0.0.0/0",
	}).Extract()

	if err != nil {
		return errors.Wrap(err, "create ssh rule for worker security group")
	}

	config.OpenStackConfig.WorkerSecurityGroupId = workerSg.ID

	return nil
}

func (s *CreateSecurityGroupStep) Name() string {
	return CreateSecurityGroupStepName
}

func (s *CreateSecurityGroupStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateSecurityGroupStep) Description() string {
	return "Create security group"
}

func (s *CreateSecurityGroupStep) Depends() []string {
	return nil
}
