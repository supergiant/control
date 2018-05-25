package kube

type Kube struct {
	KubeName          string
	AccountName       string
	KubeProfileName   string
	KubernetesVersion string
	ETCDDiscoveryURL  string

	Nodes         []Node
	LoadBalancers []LoadBalancer
}

type Node struct {
}

type LoadBalancer struct {
}
