package clouds

import "github.com/pkg/errors"

type Name string

const (
	AWS          Name = "aws"
	DigitalOcean Name = "digitalocean"
	Packet       Name = "packet"
	GCE          Name = "gce"
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
	CredsPrivateKey            = "privateKey"
	CredsPublicKey             = "publicKey"
	AWSAccessKeyID             = "access_key"
	AWSSecretKey               = "secret_key"

	ClusterIDTag = "supergiant.io/cluster-id"
	AwsAZ = "aws_az"
	AwsVpcCIDR = "aws_vpc_cidr"
	AwsVpcID = "aws_vpc_id"
	AwsKeyPairName = "aws_keypair_name"
	AwsSubnetID = "aws_subnet_id"
	AwsMastersSecGroupID = "aws_masters_secgroup_id"
	AwsNodesSecgroupID = "aws_nodes_secgroup_id"
	AwsSshBootstrapPrivateKey = "aws_ssh_bootstrap_private_key"
	AwsUserProvidedSshPublicKey = "aws_user_provided_public_key"
)

