package profile

type KubeProfile struct {
	Id                string        `json:"id" valid:"required"`
	MasterProfiles    []NodeProfile `json:"masterProfiles" valid:"required"`
	NodesProfiles     []NodeProfile `json:"nodesProfiles" valid:"required"`

	Arch              string        `json:"arch"`
	OperatingSystem   string        `json:"operatingSystem"`
	UbuntuVersion     string        `json:"ubuntuVersion"`
	Dockerversion     string        `json:"dockerVersion"`
	K8SVersion        string        `json:"K8SVersion"`
	FlannelVersion    string        `json:"flannelVersion"`
	NetworkType       string        `json:"networkType"`
	HelmVersion       string        `json:"helmVersion"`
	RBACEnabled       bool          `json:"rbacEnabled"`
}
