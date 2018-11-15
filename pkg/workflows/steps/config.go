package steps

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/storage"
)

type CertificatesConfig struct {
	KubernetesConfigDir string `json:"kubernetesConfigDir"`
	PublicIP            string `json:"publicIp"`
	PrivateIP           string `json:"privateIp"`

	IsMaster        bool   `json:"isMaster"`
	MasterPrivateIP string `json:"masterPrivateIp"`
	MasterPublicIP  string `json:"masterPublicIp"`

	Username string `json:"username"`
	Password string `json:"password"`

	AdminCert string `json:"adminCert"`
	AdminKey  string `json:"adminKey"`

	ParenCert []byte `json:"parenCert"`
	CACert    string `json:"caCert"`
	CAKey     string `json:"caKey"`
}

type DOConfig struct {
	Name string `json:"name" valid:"required"`
	// These come from UI select
	Region string `json:"region" valid:"required"`
	Size   string `json:"size" valid:"required"`
	Image  string `json:"image" valid:"required"`

	// These come from cloud account
	Fingerprint string `json:"fingerprint" valid:"required"`
	AccessToken string `json:"accessToken" valid:"required"`
}

// TODO(stgleb): Fill struct with fields when provisioning on other providers is done

type GCEConfig struct {
	PrivateKey    string `json:"privateKey"`
	ImageFamily   string `json:"imageFamily"`
	ProjectID     string `json:"projectId"`
	Zone          string `json:"zone"`
	Size          string `json:"size"`
	InstanceGroup string `json:"instanceGroup"`
	ClientEmail   string `json:"clientEmail"`
	TokenURI      string `json:"tokenURI"`
	AuthURI       string `json:"authURI"`
}

type PacketConfig struct{}

type OSConfig struct{}

type AWSConfig struct {
	KeyID                  string `json:"access_key"`
	Secret                 string `json:"secret_key"`
	Region                 string `json:"region"`
	AvailabilityZone       string `json:"availabilityZone"`
	KeyPairName            string `json:"keyPairName"`
	VPCID                  string `json:"vpcid"`
	VPCCIDR                string `json:"vpccidr"`
	SubnetID               string `json:"subnetID"`
	NodesSecurityGroupID   string `json:"nodesSecurityGroupID"`
	MastersSecurityGroupID string `json:"mastersSecurityGroupID"`
	VolumeSize             string `json:"volumeSize"`
	EbsOptimized           string `json:"ebsOptimized"`
	ImageID                string `json:"image"`
	InstanceType           string `json:"size"`
	HasPublicAddr          string `json:"hasPublicAddr"`
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
	IsMaster   bool   `json:"isMaster"`
	ProxyPort  string `json:"proxyPort"`
	K8SVersion string `json:"k8sVersion"`
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
	IsMaster    bool   `json:"isMaster"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	RBACEnabled bool   `json:"rbacEnabled"`
	Timeout     int    `json:"timeout"`
}

type TillerConfig struct {
	HelmVersion     string `json:"helmVersion"`
	RBACEnabled     bool   `json:"rbacEnabled"`
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
	Name           string        `json:"name"`
	Version        string        `json:"version"`
	DiscoveryUrl   string        `json:"discoveryUrl"`
	AdvertiseHost  string        `json:"advertiseHost"`
	Host           string        `json:"host"`
	DataDir        string        `json:"dataDir"`
	ServicePort    string        `json:"servicePort"`
	ManagementPort string        `json:"managementPort"`
	Timeout        time.Duration `json:"timeout"`
	StartTimeout   string        `json:"startTimeout"`
	RestartTimeout string        `json:"restartTimeout"`
}

type SshConfig struct {
	User                string `json:"user"`
	Port                string `json:"port"`
	BootstrapPrivateKey string `json:"bootstrapPrivateKey"`
	BootstrapPublicKey  string `json:"bootstrapPublicKey"`
	PublicKey           string `json:"publicKey"`
	Timeout             int    `json:"timeout"`
}

type ClusterCheckConfig struct {
	MachineCount int
}

type PrometheusConfig struct {
	Port        string `json:"port"`
	RBACEnabled bool   `json:"rbacEnabled"`
}

type Map struct {
	internal map[string]*node.Node
}

func (m *Map) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.internal)
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.internal)
}

type Config struct {
	TaskID                 string
	Provider               clouds.Name  `json:"provider"`
	IsMaster               bool         `json:"isMaster"`
	ClusterID              string       `json:"clusterId"`
	ClusterName            string       `json:"clusterName"`
	LogBootstrapPrivateKey bool         `json:"logBootstrapPrivateKey"`
	DigitalOceanConfig     DOConfig     `json:"digitalOceanConfig"`
	AWSConfig              AWSConfig    `json:"awsConfig"`
	GCEConfig              GCEConfig    `json:"gceConfig"`
	OSConfig               OSConfig     `json:"osConfig"`
	PacketConfig           PacketConfig `json:"packetConfig"`

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
	PrometheusConfig   PrometheusConfig   `json:"prometheusConfig"`

	ClusterCheckConfig ClusterCheckConfig `json:"clusterCheckConfig"`

	Node             node.Node     `json:"node"`
	CloudAccountID   string        `json:"cloudAccountId" valid:"required, length(1|32)"`
	CloudAccountName string        `json:"cloudAccountName" valid:"required, length(1|32)"`
	Timeout          time.Duration `json:"timeout"`
	Runner           runner.Runner `json:"-"`

	repository storage.Interface `json:"-"`

	m1      sync.RWMutex
	Masters Map `json:"masters"`

	m2    sync.RWMutex
	Nodes Map `json:"nodes"`

	nodeChan      chan node.Node
	kubeStateChan chan model.KubeState
}

// NewConfig builds instance of config for provisioning
func NewConfig(clusterName, discoveryUrl, cloudAccountName string, profile profile.Profile) *Config {
	cfg := &Config{
		Provider:    profile.Provider,
		ClusterName: clusterName,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		LogBootstrapPrivateKey: profile.LogBootstrapPrivateKey,
		AWSConfig: AWSConfig{
			Region: profile.Region,
			AvailabilityZone:       profile.CloudSpecificSettings[clouds.AwsAZ],
			VPCCIDR:                profile.CloudSpecificSettings[clouds.AwsVpcCIDR],
			VPCID:                  profile.CloudSpecificSettings[clouds.AwsVpcID],
			KeyPairName:            profile.CloudSpecificSettings[clouds.AwsKeyPairName],
			SubnetID:               profile.CloudSpecificSettings[clouds.AwsSubnetID],
			MastersSecurityGroupID: profile.CloudSpecificSettings[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:   profile.CloudSpecificSettings[clouds.AwsNodesSecgroupID],
		},
		GCEConfig: GCEConfig{
			Zone:          profile.CloudSpecificSettings["gceZone"],
			InstanceGroup: profile.CloudSpecificSettings["gceInstanceGroup"],
		},
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
			ProxyPort:  "8080",
			K8SVersion: profile.K8SVersion,
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
			Timeout:     600,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
			RBACEnabled:     profile.RBACEnabled,
		},
		SshConfig: SshConfig{
			Port:      "22",
			User:      "root",
			Timeout:   10,
			PublicKey: profile.PublicKey,
		},
		EtcdConfig: EtcdConfig{
			// TODO(stgleb): this field must be changed per node
			Name:           "etcd0",
			Version:        "3.3.9",
			Host:           "0.0.0.0",
			DataDir:        "/tmp/etcd-data",
			ServicePort:    "2379",
			ManagementPort: "2380",
			Timeout:        time.Minute * 10,
			StartTimeout:   "0",
			RestartTimeout: "5",
			DiscoveryUrl:   discoveryUrl,
		},
		ClusterCheckConfig: ClusterCheckConfig{
			MachineCount: len(profile.NodesProfiles) + len(profile.MasterProfiles),
		},
		PrometheusConfig: PrometheusConfig{
			Port:        "30900",
			RBACEnabled: profile.RBACEnabled,
		},

		Masters: Map{
			internal: make(map[string]*node.Node, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*node.Node, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 30,
		CloudAccountName: cloudAccountName,

		nodeChan:      make(chan node.Node, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan: make(chan model.KubeState, 2),
	}

	return cfg
}

// AddMaster to map of master, map is used because it is reference and can be shared among
// goroutines that run multiple tasks of cluster deployment
func (c *Config) AddMaster(n *node.Node) {
	c.m1.Lock()
	defer c.m1.Unlock()
	c.Masters.internal[n.ID] = n
}

// AddNode to map of nodes in cluster
func (c *Config) AddNode(n *node.Node) {
	c.m2.Lock()
	defer c.m2.Unlock()
	c.Nodes.internal[n.ID] = n
}

// GetMaster returns first master in master map or nil
func (c *Config) GetMaster() *node.Node {
	// non-blocking fast path for master nodes
	if c.IsMaster {
		return &c.Node
	}

	c.m1.RLock()
	defer c.m1.RUnlock()

	if len(c.Masters.internal) == 0 {
		return nil
	}

	for key := range c.Masters.internal {
		// Skip inactive nodes for selecting
		if c.Masters.internal[key] != nil && c.Masters.internal[key].State == node.StateActive {
			return c.Masters.internal[key]
		}
	}

	return nil
}

func (c *Config) GetMasters() map[string]*node.Node {
	c.m1.RLock()
	defer c.m1.RUnlock()

	m := make(map[string]*node.Node, len(c.Masters.internal))

	for key := range c.Masters.internal {
		m[c.Masters.internal[key].Name] = c.Masters.internal[key]
	}

	return m
}

func (c *Config) GetNodes() map[string]*node.Node {
	c.m2.RLock()
	defer c.m2.RUnlock()

	m := make(map[string]*node.Node, len(c.Nodes.internal))

	for key := range c.Nodes.internal {
		m[c.Nodes.internal[key].Name] = c.Nodes.internal[key]
	}

	return m
}

// GetMaster returns first master in master map or nil
func (c *Config) GetNode() *node.Node {
	c.m2.RLock()
	defer c.m2.RUnlock()

	if len(c.Nodes.internal) == 0 {
		return nil
	}

	for key := range c.Nodes.internal {
		// Skip inactive nodes for selecting
		if c.Nodes.internal[key] != nil && c.Nodes.internal[key].State == node.StateActive {
			return c.Nodes.internal[key]
		}
	}

	return nil
}

func (c *Config) NodeChan() chan node.Node {
	return c.nodeChan
}

func (c *Config) KubeStateChan() chan model.KubeState {
	return c.kubeStateChan
}
