package steps

import (
	"time"

	"github.com/supergiant/supergiant/pkg/runner"
)

type CertificatesConfig struct {
	KubernetesConfigDir   string `json:"kubernetes_config_dir"`
	CACert                string `json:"ca_cert"`
	CACertName            string `json:"ca_cert_name"`
	CAKeyCert             string `json:"ca_key_cert"`
	CAKeyName             string `json:"ca_key_name"`
	APIServerCert         string `json:"api_server_cert"`
	APIServerCertName     string `json:"api_server_cert_name"`
	APIServerKey          string `json:"api_server_key"`
	APIServerKeyName      string `json:"api_server_key_name"`
	KubeletClientCert     string `json:"kubelet_client_cert"`
	KubeletClientCertName string `json:"kubelet_client_cert_name"`
	KubeletClientKey      string `json:"kubelet_client_key"`
	KubeletClientKeyName  string `json:"kubelet_client_key_name"`
}

type DOConfig struct {
	Name         string   `json:"name"`
	K8sVersion   string   `json:"k_8_s_version"`
	Region       string   `json:"region"`
	Size         string   `json:"size"`
	Role         string   `json:"role"` // master/node
	Image        string   `json:"image"`
	Fingerprints []string `json:"fingerprints"`
	AccessToken  string   `json:"access_token"`
}

type FlannelConfig struct {
	Version     string `json:"version"`
	Arch        string `json:"arch"`
	Network     string `json:"network"`
	NetworkType string `json:"network_type"`
}

type KubeletConfig struct {
	MasterPrivateIP   string `json:"master_private_ip"`
	ProxyPort         string `json:"proxy_port"`
	EtcdClientPort    string `json:"etcd_client_port"`
	KubernetesVersion string `json:"kubernetes_version"`
}

type KubeletConfConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type KubeProxyConfig struct {
	MasterPrivateIP   string `json:"master_private_ip"`
	ProxyPort         string `json:"proxy_port"`
	EtcdClientPort    string `json:"etcd_client_port"`
	KubernetesVersion string `json:"kubernetes_version"`
}

type ManifestConfig struct {
	KubernetesVersion   string `json:"kubernetes_version"`
	KubernetesConfigDir string `json:"kubernetes_config_dir"`
	RBACEnabled         bool   `json:"rbac_enabled"`
	EtcdHost            string `json:"etcd_host"`
	EtcdPort            string `json:"etcd_port"`
	PrivateIpv4         string `json:"private_ipv_4"`
	ProviderString      string `json:"provider_string"`
	MasterHost          string `json:"master_host"`
	MasterPort          string `json:"master_port"`
}

type PostStartConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	RBACEnabled bool   `json:"rbac_enabled"`
}

type KubeletSystemdServiceConfig struct {
	KubernetesVersion  string `json:"kubernetes_version"`
	KubeletService     string `json:"kubelet_service"`
	KubernetesProvider string `json:"kubernetes_provider"`
}

type TillerConfig struct {
	HelmVersion     string `json:"helm_version"`
	OperatingSystem string `json:"operating_system"`
	Arch            string `json:"arch"`
}

type DockerConfig struct {
	DockerVersion  string `json:"docker_version"`
	ReleaseVersion string `json:"release_version"`
	Arch           string `json:"arch"`
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

	Timeout time.Duration `json:"timeout"`
	runner.Runner
}
