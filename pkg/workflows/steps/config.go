package steps

import (
	"sync"
	"time"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/storage"
)

type CertificatesConfig struct {
	KubernetesConfigDir string `json:"kubernetesConfigDir"`
	MasterPrivateIP     string `json:"masterPrivateIP"`
	Username            string `json:"username"`
	Password            string `json:"password"`
}

type DOConfig struct {
	Name string `json:"name" valid:"required"`
	// These come from UI select
	Region string `json:"region" valid:"required"`
	Size   string `json:"size" valid:"required"`
	Image  string `json:"image" valid:"required"`

	// These come from clouda ccount
	Fingerprint string `json:"fingerprint" valid:"required"`
	AccessToken string `json:"accessToken" valid:"required"`
}

// TODO(stgleb): Fill struct with fields when provisioning on other providers is done
type AWSConfig struct{}

type GCEConfig struct{}

type PacketConfig struct{}

type OSConfig struct{}

type FlannelConfig struct {
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	EtcdHost string `json:"etcdHost"`
}

type NetworkConfig struct {
	EtcdRepositoryUrl string `json:"etcdRepositoryUrl"`
	EtcdVersion       string `json:"etcdVersion"`
	EtcdHost          string `json:"etcdHost"`

	Arch            string `json:"arch"`
	OperatingSystem string `json:"operatingSystem"`

	Network     string `json:"network"`
	NetworkType string `json:"networkType"`
}

type KubeletConfig struct {
	MasterPrivateIP string `json:"masterPrivateIP"`
	ProxyPort       string `json:"proxyPort"`
	EtcdClientPort  string `json:"etcdClientPort"`
	K8SVersion      string `json:"k8sVersion"`
}

type ManifestConfig struct {
	IsMaster            bool   `json:"isMaster"`
	K8SVersion          string `json:"k8sVersion"`
	KubernetesConfigDir string `json:"kubernetesConfigDir"`
	RBACEnabled         bool   `json:"rbacEnabled"`
	ProviderString      string `json:"providerString"`
	MasterHost          string `json:"masterHost"`
	MasterPort          string `json:"masterPort"`
}

type PostStartConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	RBACEnabled bool   `json:"rbacEnabled"`
	Timeout     int    `json:"timeout"`
}

type TillerConfig struct {
	HelmVersion     string `json:"helmVersion"`
	OperatingSystem string `json:"operatingSystem"`
	Arch            string `json:"arch"`
}

type DockerConfig struct {
	Version        string `json:"version"`
	ReleaseVersion string `json:"releaseVersion"`
	Arch           string `json:"arch"`
}

type DownloadK8sBinary struct {
	K8SVersion      string `json:"k8sVersion"`
	Arch            string `json:"arch"`
	OperatingSystem string `json:"operatingSystem"`
}

type EtcdConfig struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	DiscoveryUrl   string `json:"discoveryUrl"`
	Host           string `json:"host"`
	DataDir        string `json:"dataDir"`
	ServicePort    string `json:"servicePort"`
	ManagementPort string `json:"managementPort"`
	StartTimeout   string `json:"startTimeout"`
	RestartTimeout string `json:"restartTimeout"`
}

type SshConfig struct {
	User       string `json:"user"`
	Port       string `json:"port"`
	PrivateKey []byte `json:"privateKey"`
	Timeout    int    `json:"timeout"`
}

// TODO(stgleb): rename to context and embed context.Context here
type Config struct {
	Provider    clouds.Name `json:"provider"`
	IsMaster    bool        `json:"isMaster"`
	ClusterName string      `json:"clusterName"`

	DigitalOceanConfig DOConfig     `json:"digitalOceanConfig"`
	AWSConfig          AWSConfig    `json:"awsConfig"`
	GCEConfig          GCEConfig    `json:"gceConfig"`
	OSConfig           OSConfig     `json:"osConfig"`
	PacketConfig       PacketConfig `json:"packetConfig"`

	DockerConfig       DockerConfig       `json:"dockerConfig"`
	DownloadK8sBinary  DownloadK8sBinary  `json:"downloadK8sBinary"`
	CertificatesConfig CertificatesConfig `json:"certificatesConfig"`
	FlannelConfig      FlannelConfig      `json:"flannelConfig"`
	NetworkConfig      NetworkConfig      `json:"networkConfig"`
	KubeletConfig      KubeletConfig      `json:"kubeletConfig"`
	ManifestConfig     ManifestConfig     `json:"manifestConfig"`
	PostStartConfig    PostStartConfig    `json:"postStartConfig"`
	TillerConfig       TillerConfig       `json:"tillerConfig"`
	EtcdConfig         EtcdConfig         `json:"etcdConfig"`
	SshConfig          SshConfig          `json:"sshConfig"`

	Node             node.Node     `json:"node"`
	CloudAccountName string        `json:"cloudAccountName" valid:"required, length(1|32)"`
	Timeout          time.Duration `json:"timeout"`
	Runner           runner.Runner `json:"-"`

	repository storage.Interface `json:"-"`

	// TODO(stgleb): make MasterNodes private, add public method for initialization.
	m           sync.RWMutex
	MasterNodes map[string]*node.Node `json:"masterNodes"`
}

func (c *Config) AddMaster(n *node.Node) {
	c.m.Lock()
	defer c.m.Unlock()
	c.MasterNodes[n.Id] = n
}

func (c *Config) GetMaster() *node.Node {
	c.m.RLock()
	defer c.m.RUnlock()

	if len(c.MasterNodes) == 0 {
		return nil
	}

	for _, value := range c.MasterNodes {
		return value
	}

	return nil
}
