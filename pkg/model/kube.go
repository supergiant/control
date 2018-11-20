package model

import (
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/clouds"
)

type KubeState string

const (
	StateProvisioning KubeState = "provisioning"
	StateFailed       KubeState = "failed"
	StateOperational  KubeState = "operational"
	StateDeleting     KubeState = "deleting"
)

// Kube represents a kubernetes cluster.
type Kube struct {
	ID           string    `json:"id" valid:"-"`
	State        KubeState `json:"state"`
	Name         string    `json:"name" valid:"required"`
	Provider     clouds.Name       `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	RBACEnabled  bool      `json:"rbacEnabled"`
	AccountName  string    `json:"accountName"`
	Region       string    `json:"region"`
	Zone         string    `json:"zone" valid:"-"`
	APIPort      string    `json:"apiPort"`
	Auth         Auth      `json:"auth"`
	SshUser      string    `json:"sshUser"`
	SshPublicKey []byte    `json:"sshKey"`
	User                   string                `json:"user" valid:"-"`
	Password               string                `json:"password" valid:"-"`

	Arch                   string     `json:"arch"`
	OperatingSystem        string     `json:"operatingSystem"`
	OperatingSystemVersion string     `json:"operatingSystemVersion"`
	DockerVersion          string     `json:"dockerVersion"`
	K8SVersion             string     `json:"K8SVersion"`
	HelmVersion            string     `json:"helmVersion"`
	Networking             Networking `json:"networking"`

	CloudSpec profile.CloudSpecificSettings `json:"cloudSpec" valid:"-"`

	Masters map[string]*node.Node `json:"masters"`
	Nodes   map[string]*node.Node `json:"nodes"`
}

// Auth holds all possible auth parameters.
type Auth struct {
	Username  string `json:"username"`
	Password  string `json:"token"`
	CAKey     string `json:"caKey"`
	CACert    string `json:"caCert"`
	AdminCert string `json:"adminCert"`
	AdminKey  string `json:"adminKey"`
}

type Networking struct {
	Manager string `json:"manager"`
	Version string `json:"version"`
	Type    string `json:"type"`
	CIDR    string `json:"cidr"`
}
