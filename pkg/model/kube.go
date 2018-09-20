package model

import (
	"github.com/supergiant/supergiant/pkg/node"
)

type KubeState string

const (
	StateProvisioning KubeState = "provisioning"
	StateFailed       KubeState = "failed"
	StateOperational  KubeState = "operational"
)

// Kube represents a kubernetes cluster.
type Kube struct {
	State        KubeState `json:"state"`
	Name         string    `json:"name" valid:"required"`
	RBACEnabled  bool      `json:"rbacEnabled"`
	AccountName  string    `json:"accountName"`
	Region       string    `json:"region"`
	Auth         Auth      `json:"auth"`
	SshUser      string    `json:"sshUser"`
	SshPublicKey []byte    `json:"sshKey"`

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
	Manager string `json:"manager"`
	Version string `json:"version"`
	Type    string `json:"type"`
	CIDR    string `json:"cidr"`
}
