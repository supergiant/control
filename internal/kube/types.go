package kube

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
