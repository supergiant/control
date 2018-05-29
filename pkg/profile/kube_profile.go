package profile

type KubeProfile struct {
	Id                string        `json:"id" valid:"required"`
	KubernetesVersion string        `json:"kubernetes_version" `
	Provider          string        `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	Nodes             []NodeProfile `json:"nodes" valid:"required"`
	MasterNodeCount   int           `json:"master_node_count" valid:"required"`
	CustomFiles       string        `json:"custom_files" valid:"optional"`
	RBACEnabled       bool          `json:"rbac_enabled"`
}
