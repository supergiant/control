package clouds

type Name string

const (
	AWS          Name = "aws"
	DigitalOcean Name = "digitalocean"
	Packet       Name = "packet"
	GCE          Name = "gce"
	OpenStack    Name = "openstack"
)

const (
	CredsPrivateKey = "privateKey"
)
