package profile

type KubeProfile struct {
	ID             string        `json:"id" valid:"required"`
	MasterProfiles []NodeProfile `json:"masterProfiles" valid:"required"`
	NodesProfiles  []NodeProfile `json:"nodesProfiles" valid:"required"`

	// TODO(stgleb): In future releases arch will probably migrate to node profile
	// to allow user create heterogeneous cluster of machine with different arch
	Arch            string `json:"arch"`
	OperatingSystem string `json:"operatingSystem"`
	UbuntuVersion   string `json:"ubuntuVersion"`
	DockerVersion   string `json:"dockerVersion"`
	K8SVersion      string `json:"K8SVersion"`
	FlannelVersion  string `json:"flannelVersion"`
	NetworkType     string `json:"networkType"`
	HelmVersion     string `json:"helmVersion"`
	RBACEnabled     bool   `json:"rbacEnabled"`
}
