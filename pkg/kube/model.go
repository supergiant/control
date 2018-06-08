package kube

type Kube struct {
	Name string `json:"name" valid:"required"`
	// Kubernetes
	KubernetesVersion string `json:"kubernetes_version" validate:"nonzero" sg:"default=1.8.7"`
	SSHPubKey         string `json:"ssh_pub_key"`
	ETCDDiscoveryURL  string `json:"etcd_discovery_url" sg:"readonly"`

	Username string `json:"username" validate:"nonzero" sg:"immutable"`
	Password string `json:"password" validate:"nonzero" sg:"immutable"`

	RBACEnabled bool `json:"rbac_enabled"`

	MasterPublicIP string `json:"master_public_ip" sg:"readonly"`

	Ready bool `json:"ready" sg:"readonly"`
}
