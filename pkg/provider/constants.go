package provider

type (
	// Name of provider
	Name string
	// K8SVersion version in format major.minor.patch
	K8SVersion string
	// NodeRole either master or minion
	NodeRole string
)

const (
	AWS          Name = "aws"
	DigitalOcean Name = "digitalocean"
	Packet       Name = "packet"
	GCE          Name = "gce"
	OpenStack    Name = "openstack"

	K8S15 K8SVersion = "1.5"
	K8S16 K8SVersion = "1.6"
	K8S17 K8SVersion = "1.7"
	K8S18 K8SVersion = "1.8"

	Master NodeRole = "master"
	Minion NodeRole = "minion"
)
