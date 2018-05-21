package provider

type Name string

const (
	DigitalOcean      Name = "digitalocean"
	Amazon            Name = "aws"
	GoogleCloudEngine Name = "gce"
	Packet            Name = "packet"
)
