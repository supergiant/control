package kube

type Kube struct {
	Name              string
	AccountName       string
	KubeProfileName   string
	KubernetesVersion string
	ProfileName       string
	Nodes             []Node
	LoadBalancers     []LoadBalancer
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
