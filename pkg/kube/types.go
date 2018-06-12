package kube

//type Kube struct {
//	Name string `json:"name" valid:"required"`
//	// Kubernetes
//	KubernetesVersion string `json:"kubernetes_version" validate:"nonzero" sg:"default=1.8.7"`
//	SSHPubKey         string `json:"ssh_pub_key"`
//	ETCDDiscoveryURL  string `json:"etcd_discovery_url" sg:"readonly"`
//
//	Username string `json:"username" validate:"nonzero" sg:"immutable"`
//	Password string `json:"password" validate:"nonzero" sg:"immutable"`
//
//	RBACEnabled bool `json:"rbac_enabled"`
//
//	MasterPublicIP string `json:"master_public_ip" sg:"readonly"`
//
//	Ready bool `json:"ready" sg:"readonly"`
//}

type Kube struct {
	Name              string
	AccountName       string
	KubeProfileName   string
	KubernetesVersion string
	ProfileName       string
	Nodes             []Node
	LoadBalancers     []LoadBalancer
	APIAddr           string
	Auth              Auth
}

type Node struct {
	KubeName    string
	ProfileName string
}

type LoadBalancer struct {
	KubeName  string
	Namespace string
	Selector  map[string]string
	Ports     map[int]int
	IP        string
}

// Auth holds all possible auth parameters.
type Auth struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	CA       string `json:"ca"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
