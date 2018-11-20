package profile

import "github.com/supergiant/supergiant/pkg/clouds"

type Profile struct {
	ID string `json:"id" valid:"required"`

	MasterProfiles []NodeProfile `json:"masterProfiles" valid:"-"`
	NodesProfiles  []NodeProfile `json:"nodesProfiles" valid:"-"`

	// TODO(stgleb): In future releases arch will probably migrate to node profile
	// to allow user create heterogeneous cluster of machine with different arch
	Provider               clouds.Name           `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)" valid:"-"`
	Region                 string                `json:"region" valid:"-"`
	Zone                   string                `json:"zone" valid:"-"`
	Arch                   string                `json:"arch" valid:"-"`
	OperatingSystem        string                `json:"operatingSystem" valid:"-"`
	UbuntuVersion          string                `json:"ubuntuVersion" valid:"-"`
	DockerVersion          string                `json:"dockerVersion" valid:"-"`
	K8SVersion             string                `json:"K8SVersion" valid:"-"`
	FlannelVersion         string                `json:"flannelVersion" valid:"-"`
	NetworkType            string                `json:"networkType" valid:"-"`
	CIDR                   string                `json:"cidr" valid:"-"`
	HelmVersion            string                `json:"helmVersion" valid:"-"`
	User                   string                `json:"user" valid:"-"`
	Password               string                `json:"password" valid:"-"`
	RBACEnabled            bool                  `json:"rbacEnabled" valid:"-"`
	CloudSpecificSettings  CloudSpecificSettings `json:"cloudSpecificSettings" valid:"-"`
	PublicKey              string                `json:"publicKey" valid:"-"`
	LogBootstrapPrivateKey bool                  `json:"logBootstrapPrivateKey" valid:"-"`
}

type NodeProfile map[string]string
type CloudSpecificSettings map[string]string
