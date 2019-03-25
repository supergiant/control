package azure

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
)

const CreateSecurityGroupStepName = "CreateNetworkSecurityGroup"

type NSGClientFn func(a autorest.Authorizer, subscriptionID string) (SecurityGroupInterface, autorest.Client)

type SubnetClientFn func(a autorest.Authorizer, subscriptionID string) SubnetGetter

type CreateSecurityGroupStep struct {
	nsgClientFn    NSGClientFn
	subnetGetterFn SubnetClientFn
	findOutboundIP func(ctx context.Context) (string, error)
}

func NewCreateSecurityGroupStep() *CreateSecurityGroupStep {
	return &CreateSecurityGroupStep{
		nsgClientFn:    NSGClientFor,
		subnetGetterFn: SubnetClientFor,
		findOutboundIP: func(ctx context.Context) (string, error) {
			return amazon.FindOutboundIP(ctx, amazon.FindExternalIP)
		},
	}
}

func (s *CreateSecurityGroupStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.nsgClientFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "security group client builder")
	}
	if s.subnetGetterFn == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "subnets client builder")
	}

	nsgClient, restclient := s.nsgClientFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)
	subnetClient := s.subnetGetterFn(config.GetAzureAuthorizer(), config.AzureConfig.SubscriptionID)

	sgAddr, err := s.findOutboundIP(ctx)
	if err != nil {
		return errors.Wrap(err, "get supergiant address")
	}

	// default security rules:
	// https://docs.microsoft.com/en-us/azure/virtual-network/security-overview#default-security-rules
	for _, r := range []struct {
		role  string
		rules []network.SecurityRule
	}{
		{
			role:  model.RoleMaster.String(),
			rules: masterSecurityRules(sgAddr),
		},
		{
			role:  model.RoleNode.String(),
			rules: nodeSecurityRules(sgAddr),
		},
	} {
		subnet, err := subnetClient.Get(
			ctx,
			toResourceGroupName(config.ClusterID, config.ClusterName),
			toVNetName(config.ClusterID, config.ClusterName),
			toSubnetName(config.ClusterID, config.ClusterName, r.role),
			"",
		)
		if err != nil {
			return errors.Wrapf(err, "get %s subnet", toSubnetName(config.ClusterID, config.ClusterName, r.role))
		}

		name := toNSGName(config.ClusterID, config.ClusterName, r.role)
		f, err := nsgClient.CreateOrUpdate(
			ctx,
			toResourceGroupName(config.ClusterID, config.ClusterName),
			name,
			network.SecurityGroup{
				Location: to.StringPtr(config.AzureConfig.Location),
				SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
					SecurityRules: &r.rules,
					Subnets:       &[]network.Subnet{subnet},
				},
			},
		)
		if err != nil {
			return errors.Wrapf(err, "create %s network security group", name)
		}

		if err = f.WaitForCompletionRef(ctx, restclient); err != nil {
			return errors.Wrapf(err, "wait for %s security group is ready", name)
		}

	}

	return nil
}

func (s *CreateSecurityGroupStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateSecurityGroupStep) Name() string {
	return CreateSecurityGroupStepName
}

func (s *CreateSecurityGroupStep) Depends() []string {
	return nil
}

func (s *CreateSecurityGroupStep) Description() string {
	return "Azure: Create master/node network security groups"
}

func masterSecurityRules(sgAddr string) []network.SecurityRule {
	return []network.SecurityRule{
		{
			Name: to.StringPtr("allow_ssh_for_sg"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr(fmt.Sprintf("%s/32", sgAddr)),
				SourcePortRange:          to.StringPtr("1-65535"),
				DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
				DestinationPortRange:     to.StringPtr("22"),
				Access:                   network.SecurityRuleAccessAllow,
				Direction:                network.SecurityRuleDirectionInbound,
				Priority:                 to.Int32Ptr(100),
			},
		},
		{
			Name: to.StringPtr("allow_https_for_sg"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr(fmt.Sprintf("%s/32", sgAddr)),
				SourcePortRange:          to.StringPtr("1-65535"),
				DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
				DestinationPortRange:     to.StringPtr("443"),
				Access:                   network.SecurityRuleAccessAllow,
				Direction:                network.SecurityRuleDirectionInbound,
				Priority:                 to.Int32Ptr(200),
			},
		},
	}
}

func nodeSecurityRules(sgAddr string) []network.SecurityRule {
	return []network.SecurityRule{
		{
			Name: to.StringPtr("allow_ssh_for_sg"),
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr(fmt.Sprintf("%s/32", sgAddr)),
				SourcePortRange:          to.StringPtr("1-65535"),
				DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
				DestinationPortRange:     to.StringPtr("22"),
				Access:                   network.SecurityRuleAccessAllow,
				Direction:                network.SecurityRuleDirectionInbound,
				Priority:                 to.Int32Ptr(100),
			},
		},
	}
}
