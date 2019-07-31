package azure

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateLBStepName = "CreateLoadBalancer"
)

type CreateLBStep struct {
	sdk SDKInterface
}

func NewCreateLBStep(s SDK) *CreateLBStep {
	return &CreateLBStep{
		sdk: s,
	}
}

func (s *CreateLBStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	if config == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "config")
	}
	if s.sdk == nil {
		return errors.Wrap(sgerrors.ErrNilEntity, "azure sdk")
	}

	if err := ensureAuthorizer(s.sdk, config); err != nil {
		return errors.Wrap(err, "ensure authorization")
	}

	addr, lb, err := s.createLB(
		ctx,
		config.GetAzureAuthorizer(),
		config.AzureConfig.SubscriptionID,
		config.AzureConfig.Location,
		toResourceGroupName(config.Kube.ID, config.Kube.Name),
		toLBName(config.Kube.ID, config.Kube.Name),
		config.Kube.APIServerPort,
	)
	if err != nil {
		return errors.Wrap(err, "create load balancer")
	}
	if addr == "" {
		return errors.Wrapf(sgerrors.ErrRawError, "%s load balancer address is unknown", to.String(lb.Name))
	}

	config.Kube.ExternalDNSName = addr
	config.Kube.InternalDNSName = addr

	logrus.Debugf("azure: %s lb has been created", to.String(lb.Name))
	return nil
}

func (s *CreateLBStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateLBStep) Name() string {
	return CreateLBStepName
}

func (s *CreateLBStep) Depends() []string {
	return []string{CreateGroupStepName}
}

func (s *CreateLBStep) Description() string {
	return "Azure: Create Load Balancer"
}

func (s *CreateLBStep) createLB(ctx context.Context, a autorest.Authorizer, subsID, location, groupName, lbName string, lbPort int64) (string, network.LoadBalancer, error) {
	pipName := toIPName(lbName)
	probeName := "apiserver"
	frontEndIPConfigName := "apiserver"
	backEndAddressPoolName := "masters"
	idPrefix := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers", subsID, groupName)

	pip, err := s.createPublicIP(ctx, a, subsID, location, groupName, pipName)
	if err != nil {
		return "", network.LoadBalancer{}, errors.Wrap(err, "create public ip")
	}

	lbClient := s.sdk.LBClient(a, subsID)
	f, err := lbClient.CreateOrUpdate(
		ctx,
		groupName,
		lbName,
		network.LoadBalancer{
			Location: to.StringPtr(location),
			LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
					{
						Name: &frontEndIPConfigName,
						FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: network.Dynamic,
							PublicIPAddress:           &pip,
						},
					},
				},
				BackendAddressPools: &[]network.BackendAddressPool{
					{
						Name: &backEndAddressPoolName,
					},
				},
				Probes: &[]network.Probe{
					{
						Name: &probeName,
						ProbePropertiesFormat: &network.ProbePropertiesFormat{
							Protocol:          network.ProbeProtocolTCP,
							Port:              to.Int32Ptr(int32(lbPort)),
							IntervalInSeconds: to.Int32Ptr(5),
							NumberOfProbes:    to.Int32Ptr(4),
						},
					},
				},
				LoadBalancingRules: &[]network.LoadBalancingRule{
					{
						Name: to.StringPtr("apiserver"),
						LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
							Protocol:             network.TransportProtocolTCP,
							FrontendPort:         to.Int32Ptr(int32(lbPort)),
							BackendPort:          to.Int32Ptr(int32(lbPort)),
							IdleTimeoutInMinutes: to.Int32Ptr(4),
							EnableFloatingIP:     to.BoolPtr(false),
							LoadDistribution:     network.LoadDistributionDefault,
							FrontendIPConfiguration: &network.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, lbName, frontEndIPConfigName)),
							},
							BackendAddressPool: &network.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", idPrefix, lbName, backEndAddressPoolName)),
							},
							Probe: &network.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/probes/%s", idPrefix, lbName, probeName)),
							},
						},
					},
				},
			},
		},
	)

	if err = f.WaitForCompletionRef(ctx, lbClient.Client); err != nil {
		return "", network.LoadBalancer{}, err
	}

	lb, err := f.Result(lbClient)
	return to.String(pip.IPAddress), lb, err
}

func (s *CreateLBStep) createPublicIP(ctx context.Context, a autorest.Authorizer, subsID, location, groupName, ipName string) (network.PublicIPAddress, error) {
	f, err := s.sdk.PublicAddressesClient(a, subsID).CreateOrUpdate(
		ctx,
		groupName,
		ipName,
		network.PublicIPAddress{
			Name:     to.StringPtr(ipName),
			Location: to.StringPtr(location),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPv4,
				PublicIPAllocationMethod: network.Static,
			},
		},
	)
	if err != nil {
		return network.PublicIPAddress{}, err
	}

	err = f.WaitForCompletionRef(ctx, s.sdk.RestClient(a, subsID))
	if err != nil {
		return network.PublicIPAddress{}, errors.Wrap(err, "wait for public ip address is ready")
	}

	return s.sdk.PublicAddressesClient(a, subsID).Get(ctx, groupName, ipName, "")
}
