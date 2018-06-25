package profile

type KubeProfile struct {
	Id                string      `json:"id" valid:"required"`
	KubernetesVersion string      `json:"kubernetes_version" valid:"required"`
	MasterProfile     NodeProfile `json:"masterProfile" valid:"required"`
	NodesProfile      NodeProfile `json:"nodesProfile" valid:"required"`
	MasterNodeCount   int         `json:"master_node_count" valid:"required"`
	CustomFiles       string      `json:"customFiles" valid:"optional"`
	RBACEnabled       bool        `json:"rbacEnabled"`
}
