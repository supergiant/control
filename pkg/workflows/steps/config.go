package steps

import (
	"sync"
	"time"

	"encoding/json"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
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

type GCEConfig struct{}

type PacketConfig struct{}

type OSConfig struct{}

type AWSConfig struct {
	KeyID  string `json:"keyID"`
	Secret string `json:"secret"`

	EC2Config        EC2Config `json:"ec2config"`
	Region           string    `json:"region"`
	AvailabilityZone string    `json:"availabilityZone"`

	KeyPairName string `json:"keyPairName"`
}

type EC2Config struct {
	VolumeSize    int    `json:"volumeSize"`
	EbsOptimized  bool   `json:"ebsOptimized"`
	GPU           bool   `json:"gpu"`
	ImageID       string `json:"imageId"`
	InstanceType  string `json:"instanceType"`
	SubnetID      string `json:"subnetID"`
	HasPublicAddr bool   `json:"hasPublicAddr"`
}

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
	PrivateKey string `json:"privateKey"`
	Timeout    int    `json:"timeout"`
}

type MasterMap struct {
	internal map[string]*node.Node
}

func (m *MasterMap) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.internal)
}

func (m *MasterMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.internal)
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

	//TODO @stgleb @yegor Add possiblity to not preserve ssh keys after provisioning
	DeleteSSHKeys    bool          `json:"deleteSSHKeys"`
	Node             node.Node     `json:"node"`
	CloudAccountName string        `json:"cloudAccountName" valid:"required, length(1|32)"`
	Timeout          time.Duration `json:"timeout"`
	Runner           runner.Runner `json:"-"`

	repository storage.Interface `json:"-"`

	m           sync.RWMutex
	MasterNodes MasterMap `json:"MasterNodes"`
}

func NewConfig(clusterName, discoveryUrl, cloudAccountName string, profile profile.Profile) *Config {
	return &Config{
		ClusterName: clusterName,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		AWSConfig:    AWSConfig{},
		GCEConfig:    GCEConfig{},
		OSConfig:     OSConfig{},
		PacketConfig: PacketConfig{},

		DockerConfig: DockerConfig{
			Version:        profile.DockerVersion,
			ReleaseVersion: profile.UbuntuVersion,
			Arch:           profile.Arch,
		},
		DownloadK8sBinary: DownloadK8sBinary{
			K8SVersion:      profile.K8SVersion,
			Arch:            profile.Arch,
			OperatingSystem: profile.OperatingSystem,
		},
		CertificatesConfig: CertificatesConfig{
			KubernetesConfigDir: "/etc/kubernetes",
			Username:            "root",
			Password:            "1234",
		},
		NetworkConfig: NetworkConfig{
			EtcdRepositoryUrl: "https://github.com/coreos/etcd/releases/download",
			EtcdVersion:       "3.3.9",
			EtcdHost:          "0.0.0.0",

			Arch:            profile.Arch,
			OperatingSystem: profile.OperatingSystem,

			Network:     profile.CIDR,
			NetworkType: profile.NetworkType,
		},
		FlannelConfig: FlannelConfig{
			Arch:    profile.Arch,
			Version: profile.FlannelVersion,
			// TODO(stgleb): this should be configurable from user side
			EtcdHost: "0.0.0.0",
		},
		KubeletConfig: KubeletConfig{
			MasterPrivateIP: "localhost",
			ProxyPort:       "8080",
			EtcdClientPort:  "2379",
			K8SVersion:      profile.K8SVersion,
		},
		ManifestConfig: ManifestConfig{
			K8SVersion:          profile.K8SVersion,
			KubernetesConfigDir: "/etc/kubernetes",
			RBACEnabled:         profile.RBACEnabled,
			ProviderString:      "todo",
			MasterHost:          "localhost",
			MasterPort:          "8080",
		},
		PostStartConfig: PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    "root",
			RBACEnabled: profile.RBACEnabled,
			Timeout:     300,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
		},
		SshConfig: SshConfig{
			Port:       "22",
			User:       "root",
			PrivateKey: "",
			Timeout:    10,
		},
		EtcdConfig: EtcdConfig{
			// TODO(stgleb): this field must be changed per node
			Name:           "etcd0",
			Version:        "3.3.9",
			Host:           "0.0.0.0",
			DataDir:        "/etcd-data",
			ServicePort:    "2379",
			ManagementPort: "2380",
			StartTimeout:   "0",
			RestartTimeout: "5",
			DiscoveryUrl:   discoveryUrl,
		},

		MasterNodes: MasterMap{
			internal: make(map[string]*node.Node, len(profile.MasterProfiles)),
		},
		Timeout:          time.Second * 1200,
		CloudAccountName: cloudAccountName,
	}
}

func (c *Config) AddMaster(n *node.Node) {
	c.m.Lock()
	defer c.m.Unlock()
	c.MasterNodes.internal[n.Id] = n
}

func (c *Config) GetMaster() *node.Node {
	c.m.RLock()
	defer c.m.RUnlock()

	if len(c.MasterNodes.internal) == 0 {
		return nil
	}

	for _, value := range c.MasterNodes.internal {
		return value
	}

	return nil
}
