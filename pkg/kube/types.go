package kube

// Kube represents a kubernetes cluster.
type Kube struct {
	Name              string `json:"name" valid:"required"`
	Version           string
	RBACEnabled       bool
	AccountName       string
	KubeProfileName   string
	KubernetesVersion string
	ProfileName       string
	Nodes             []Node
	LoadBalancers     []LoadBalancer
	APIAddr           string
	Auth              Auth
}

// Node represents a kubernetes worker node.
type Node struct {
	KubeName    string
	ProfileName string
}

// LoadBalancer represents a cloud provider one.
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
