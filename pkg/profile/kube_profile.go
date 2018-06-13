package profile

type KubeProfile struct {
	Id                string      `json:"id" valid:"required"`
	KubernetesVersion string      `json:"kubernetes_version" valid:"required"`
	MasterProfile     NodeProfile `json:"master_profile" valid:"required"`
	NodesProfile      NodeProfile `json:"nodes_profile" valid:"required"`
	MasterNodeCount   int         `json:"master_node_count" valid:"required"`
	CustomFiles       string      `json:"custom_files" valid:"optional"`
	RBACEnabled       bool        `json:"rbac_enabled"`
}
