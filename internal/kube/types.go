package kube

//
//type Kube struct {
//	Name              string
//	AccountName       string
//	KubeProfileName   string
//	KubernetesVersion string
//	ProfileName       string
//	Nodes             []Node
//	LoadBalancers     []LoadBalancer
//}
//
//type Node struct {
//	KubeName    string
//	ProfileName string
//}
//
//type LoadBalancer struct {
//	KubeName  string
//	Namespace string
//	Selector  map[string]string
//	Ports     map[int]int
//	IP        string
//}

type Kube struct {
	ID      string `json:"id"`
	APIHost string `json:"apiHost"`
	APIPort string `json:"apiPort"`
	Auth    Auth   `json:"auth"`
}

type Auth struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	CA       string `json:"ca"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
