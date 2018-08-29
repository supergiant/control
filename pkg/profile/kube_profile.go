package profile

import "github.com/supergiant/supergiant/pkg/clouds"

type KubeProfile struct {
	ID string `json:"id" valid:"required"`

	MasterProfiles []map[string]string `json:"masterProfiles" valid:"required"`
	NodesProfiles  []map[string]string `json:"nodesProfiles" valid:"required"`

	// TODO(stgleb): In future releases arch will probably migrate to node profile
	// to allow user create heterogeneous cluster of machine with different arch
	Provider        clouds.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	Region          string      `json:"region"`
	Arch            string      `json:"arch"`
	OperatingSystem string      `json:"operatingSystem"`
	UbuntuVersion   string      `json:"ubuntuVersion"`
	DockerVersion   string      `json:"dockerVersion"`
	K8SVersion      string      `json:"K8SVersion"`
	FlannelVersion  string      `json:"flannelVersion"`
	NetworkType     string      `json:"networkType"`
	CIDR            string      `json:"cidr"`
	HelmVersion     string      `json:"helmVersion"`
	RBACEnabled     bool        `json:"rbacEnabled"`
}
