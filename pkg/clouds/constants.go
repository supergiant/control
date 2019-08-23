package clouds

import (
	"github.com/pkg/errors"
)

type Name string

const (
	AWS          Name = "aws"
	DigitalOcean Name = "digitalocean"
	Packet       Name = "packet"
	GCE          Name = "gce"
	Azure        Name = "azure"
	OpenStack    Name = "openstack"

	Unknown Name = "unknown"
)

var versions = []string{"1.11.5", "1.12.7", "1.13.7", "1.14.3", "1.15.1"}

func ToProvider(name string) (Name, error) {
	switch name {
	case string(AWS):
		return AWS, nil
	case string(DigitalOcean):
		return DigitalOcean, nil
	case string(Packet):
		return Packet, nil
	case string(GCE):
		return GCE, nil
	case string(OpenStack):
		return OpenStack, nil
	}
	return Unknown, errors.New("invalid provider")
}

func GetVersions() []string {
	if versions == nil {
		return nil
	}
	c := make([]string, len(versions))
	copy(c, versions)
	return c
}

const (
	OSUser = "supergiant"

	DigitalOceanFingerPrint = "fingerprint"
	DigitalOceanAccessToken = "accessToken"

	DigitalOceanExternalLoadBalancerID = "externalLoadBalancerID"
	DigitalOceanInternalLoadBalancerID = "internalLoadBalancerID"

	EnvDigitalOceanAccessToken = "DIGITALOCEAN_TOKEN"

	GCEProjectID   = "project_id"
	GCEPrivateKey  = "private_key"
	GCEClientEmail = "client_email"
	GCETokenURI    = "token_uri"

	GCETargetPoolName = "gceTargetPoolName"

	GCEExternalIPAddress = "gceExternalIpAddress"
	GCEInternalIPAddress = "gceInternalIpAddress"

	GCEExternalIPAddressName = "gceExternalIpAddressName"
	GCEInternalIPAddressName = "gceInternalIpAddressName"

	GCEExternalForwardingRuleName = "gceExternalForwardingRuleName"
	GCEInternalForwardingRuleName = "gceInternalForwardingRuleName"

	GCEBackendServiceName = "gceBackendServiceName"
	GCEBackendServiceLink = "gceBackendServiceLink"

	GCEHealthCheckName = "gceHealthCheckName"

	GCENetworkLink = "gceNetworkLink"
	GCENetworkName = "gceNetworkName"

	GCEImageFamily = "gceImageFamily"

	TagClusterID         = "supergiant.io/cluster-id"
	TagNodeName          = "Name"
	TagKubernetesCluster = "KubernetesCluster"

	AWSAccessKeyID              = "access_key"
	AWSSecretKey                = "secret_key"
	AwsAZ                       = "aws_az"
	AwsVpcCIDR                  = "aws_vpc_cidr"
	AwsVpcID                    = "aws_vpc_id"
	AwsKeyPairName              = "aws_keypair_name"
	AwsSubnets                  = "aws_subnets"
	AwsMastersSecGroupID        = "aws_masters_secgroup_id"
	AwsNodesSecgroupID          = "aws_nodes_secgroup_id"
	AwsSshBootstrapPrivateKey   = "aws_ssh_bootstrap_private_key"
	AwsUserProvidedSshPublicKey = "aws_user_provided_public_key"
	AwsRouteTableID             = "aws_route_table_id"
	AwsInternetGateWayID        = "aws_internet_gateway_id"
	AwsMasterInstanceProfile    = "aws_master_instance_profile"
	AwsNodeInstanceProfile      = "aws_node_instance_profile"
	AwsImageID                  = "aws_image_id"
	AwsExternalLoadBalancerName = "AwsExternalLoadBalancerName"
	AwsInternalLoadBalancerName = "AwsInternalLoadBalancerName"
	AwsVolumeSize               = "AwsVolumeSize"

	// Use client credentials auth model for azure.
	// https://github.com/Azure/azure-sdk-for-go#more-authentication-details
	AzureTenantID       = "tenantId"
	AzureSubscriptionID = "subscriptionId"
	AzureClientID       = "clientId"
	AzureClientSecret   = "clientSecret"
	AzureVolumeSize     = "azureVolumeSize"
	AzureVNetCIDR       = "azureVNetCIDR"

	OpenStackAuthUrl               = "authUrl"
	OpenStackDomainId              = "domainId"
	OpenStackTenantName            = "tenantName"
	OpenStackUserName              = "userName"
	OpenStackPassword              = "password"
	OpenStackDomainName            = "domainName"
	OpenStackTenantId              = "tenantId"
	OpenStackFlavorName            = "flavorName"
	OpenStackNetworkId             = "networkId"
	OpenStackSubnetId              = "subnetId"
	OpenStackSubnetCIDR            = "subnetCIDR"
	OpenStackRouterId              = "routerId"
	OpenStackListenerId            = "listenerId"
	OpenStackHealthCheckId         = "healthCheckId"
	OpenStackPoolId                = "poolId"
	OpenStackImageId               = "imageId"
	OpenStackKeypairName           = "keyPairName"
	OpenStackMasterSecurityGroupId = "masterSecurityGroupId"
	OpenStackWorkerSecurityGroupId = "workerSecurityGroupId"
)
