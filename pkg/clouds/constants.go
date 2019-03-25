package clouds

import "github.com/pkg/errors"

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

const (
	DigitalOceanFingerPrint    = "fingerprint"
	DigitalOceanAccessToken    = "accessToken"
	EnvDigitalOceanAccessToken = "DIGITALOCEAN_TOKEN"

	GCEProjectID   = "project_id"
	GCEPrivateKey  = "private_key"
	GCEClientEmail = "client_email"
	GCETokenURI    = "token_uri"

	GCETargetPoolName     = "gceTargetPoolName"
	GCEBackendServiceName = "gceBackendServiceName"

	GCEExternalIPAddress = "gceExternalIpAddress"
	GCEInternalIPAddress = "gceInternalIpAddress"

	GCEHealthCheckName = "gceHealthCheckName"

	GCEExternalIPAddressName = "gceExternalIpAddressName"
	GCEInternalIPAddressName = "gceInternalIpAddressName"

	TagClusterID         = "supergiant.io/cluster-id"
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

	// Use client credentials auth model for azure.
	// https://github.com/Azure/azure-sdk-for-go#more-authentication-details
	AzureTenantID       = "tenantId"
	AzureSubscriptionID = "subscriptionId"
	AzureClientID       = "clientId"
	AzureClientSecret   = "clientSecret"

	AzureVNetCIDR = "azureVNetCIDR"
)
