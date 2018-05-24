package profile

type KubeProfile struct {
	KubernetesVersion string        `json:"kubernetes_version"`
	RBACEnabled       bool          `json:"rbac_enabled"`
	Nodes             []NodeProfile `json:"nodes"`
	MasterNodeCount   int           `json:"master_node_count"`
	CustomFiles       string        `json:"custom_files"`
}
