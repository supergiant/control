package steps

import (
	"github.com/supergiant/supergiant/pkg/runner"
)

type CertificatesConfig struct {
	KubernetesConfigDir   string
	CACert                string
	CACertName            string
	CAKeyCert             string
	CAKeyName             string
	APIServerCert         string
	APIServerCertName     string
	APIServerKey          string
	APIServerKeyName      string
	KubeletClientCert     string
	KubeletClientCertName string
	KubeletClientKey      string
	KubeletClientKeyName  string
}

type DOConfig struct {
	Name         string
	K8sVersion   string
	Region       string
	Size         string
	Role         string // master/node
	Image        string
	Fingerprints []string
	AccessToken  string
}

type FlannelConfig struct {
	Version     string
	Arch        string
	Network     string
	NetworkType string
}

type KubeletConfig struct {
	MasterPrivateIP   string
	ProxyPort         string
	EtcdClientPort    string
	KubernetesVersion string
}

type KubeletConfConfig struct {
	Host string
	Port string
}

type KubeProxyConfig struct {
	MasterPrivateIP   string
	ProxyPort         string
	EtcdClientPort    string
	KubernetesVersion string
}

type ManifestConfig struct {
	KubernetesVersion   string
	KubernetesConfigDir string
	RBACEnabled         bool
	EtcdHost            string
	EtcdPort            string
	PrivateIpv4         string
	ProviderString      string
	MasterHost          string
	MasterPort          string
}

type PostStartConfig struct {
	Host        string
	Port        string
	Username    string
	RBACEnabled bool
}

type KubeletSystemdServiceConfig struct {
	KubernetesVersion  string
	KubeletService     string
	KubernetesProvider string
}

type TillerConfig struct {
	HelmVersion     string
	OperatingSystem string
	Arch            string
}

type DockerConfig struct {
	DockerVersion  string
	ReleaseVersion string
	Arch           string
}
type Config struct {
	DockerConfig
	CertificatesConfig
	DOConfig
	FlannelConfig
	KubeletConfig
	KubeletConfConfig
	KubeProxyConfig
	ManifestConfig
	PostStartConfig
	KubeletSystemdServiceConfig
	TillerConfig

	runner.Runner
}
