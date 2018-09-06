package kube

import (
	"github.com/supergiant/supergiant/pkg/node"
)

// Kube represents a kubernetes cluster.
type Kube struct {
	Name        string `json:"name" valid:"required"`
	RBACEnabled bool   `json:"rbacEnabled"`
	AccountName string `json:"accountName"`
	Auth        Auth   `json:"auth"`
	SSHUser     string `json:"sshUser"`
	SSHKey      []byte `json:"sshKey"`

	Arch                   string     `json:"arch"`
	OperatingSystem        string     `json:"operatingSystem"`
	OperatingSystemVersion string     `json:"operatingSystemVersion"`
	DockerVersion          string     `json:"dockerVersion"`
	K8SVersion             string     `json:"K8SVersion"`
	HelmVersion            string     `json:"helmVersion"`
	Networking             Networking `json:"networking"`

	Masters []*node.Node `json:"masters"`
	Nodes   []*node.Node `json:"nodes"`
}

// Auth holds all possible auth parameters.
type Auth struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	CA       string `json:"ca"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

type Networking struct {
	Manager string `json:"networkManager"`
	Version string `json:"flannelVersion"`
	Type    string `json:"networkType"`
	CIDR    string `json:"cidr"`
}
