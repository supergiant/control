package util

import "time"

const (
	DefaultTimeout = 3 * time.Second

	UbuntuSlug1604 = "ubuntu-16-04-x64"
	K8SVersion1102 = "1.10.2"
	K8SVersion187  = "1.8.7"
	K8SVersion177  = "1.7.7"
	K8SVersion167  = "1.6.7"
	K8SVersion157  = "1.5.7"
)

type ProviderName string

const (
	DigitalOcean      ProviderName = "digitalocean"
	Amazon            ProviderName = "aws"
	GoogleCloudEngine ProviderName = "gce"
	Packet            ProviderName = "packet"
)
