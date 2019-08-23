package model

import (
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/profile"
)

type KubeState string

const (
	StateProvisioning KubeState = "provisioning"
	StateFailed       KubeState = "failed"
	StateOperational  KubeState = "operational"
	StateDeleting     KubeState = "deleting"
	StateImporting    KubeState = "importing"
	StateUpgrading    KubeState = "upgrading"
)

// Kube represents a kubernetes cluster.
type Kube struct {
	ID           string      `json:"id" valid:"-"`
	State        KubeState   `json:"state"`
	Name         string      `json:"name" valid:"required"`
	Provider     clouds.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	RBACEnabled  bool        `json:"rbacEnabled"`
	AccountName  string      `json:"accountName"`
	Region       string      `json:"region"`
	Zone         string      `json:"zone" valid:"-"`
	ServicesCIDR string      `json:"servicesCIDR"`
	DNSIP        string      `json:"dnsIp"`
	// DEPRECATED: use APIServerPort instead.
	APIPort       string `json:"apiPort"`
	APIServerPort int    `json:"apibindPort"`
	Auth          Auth   `json:"auth"`

	User     string `json:"user" valid:"-"`
	Password string `json:"password" valid:"-"`

	Arch                   string            `json:"arch"`
	OperatingSystem        string            `json:"operatingSystem"`
	OperatingSystemVersion string            `json:"operatingSystemVersion"`
	DockerVersion          string            `json:"dockerVersion"`
	K8SVersion             string            `json:"K8SVersion"`
	HelmVersion            string            `json:"helmVersion"`
	Networking             Networking        `json:"networking"`
	Subnets                map[string]string `json:"subnets"`

	ExternalDNSName string `json:"externalDNSName"`
	InternalDNSName string `json:"internalDNSName"`
	BootstrapToken  string `json:"bootstrapToken"`

	CloudSpec profile.CloudSpecificSettings `json:"cloudSpec" valid:"-"`

	ProfileID string `json:"profileId"`

	Masters map[string]*Machine `json:"masters"`
	Nodes   map[string]*Machine `json:"nodes"`
	// Store taskIds of tasks that are made to provision this kube
	Tasks map[string][]string `json:"tasks"`

	SSHConfig SSHConfig `json:"sshConfig"`

	UserData         string              `json:"userData"`
	ExposedAddresses []profile.Addresses `json:"exposedAddresses"`
	Addons           []string            `json:"addons,omitempty"`
}

type SSHConfig struct {
	User                string `json:"user"`
	Port                string `json:"port"`
	BootstrapPrivateKey string `json:"bootstrapPrivateKey"`
	BootstrapPublicKey  string `json:"bootstrapPublicKey"`
	PublicKey           string `json:"publicKey"`
	Timeout             int    `json:"timeout"`
}

// Auth holds all possible auth parameters.
type Auth struct {
	// DEPRECATED: use static auth
	Username string `json:"username"`
	// DEPRECATED: use static auth
	Password       string             `json:"token"`
	ParentCert     string             `json:"parentCert"`
	CAKey          string             `json:"caKey"`
	CACert         string             `json:"caCert"`
	CACertHash     string             `json:"caCertHash"`
	AdminCert      string             `json:"adminCert"`
	AdminKey       string             `json:"adminKey"`
	CertificateKey string             `json:"certificateKey"`
	StaticAuth     profile.StaticAuth `json:"staticAuth"`
}

type Networking struct {
	// DEPRECATED: use Provider field instead
	Manager  string `json:"manager"`
	Provider string `json:"provider"`
	Version  string `json:"version"`
	Type     string `json:"type"`
	CIDR     string `json:"cidr"`
}
