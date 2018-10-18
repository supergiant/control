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
)
