package steps

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/pborman/uuid"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/storage"
)

type CertificatesConfig struct {
	KubernetesConfigDir string `json:"kubernetesConfigDir"`
	PublicIP            string `json:"publicIp"`
	PrivateIP           string `json:"privateIp"`
	IsMaster            bool   `json:"isMaster"`
	// TODO: this shouldn't be a part of SANs
	// https://kubernetes.io/docs/setup/certificates/#all-certificates
	KubernetesSvcIP string `json:"kubernetesSvcIp"`

	StaticAuth profile.StaticAuth `json:"staticAuth"`

	// DEPRECATED: it's a part of staticAuth
	Username string
	// DEPRECATED: it's a part of staticAuth
	Password string

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
	// NOTE(stgleb): This comes from cloud account
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
	ProjectID   string `json:"project_id"`

	// This comes from profile
	ImageFamily      string `json:"imageFamily"`
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone"`
	Size             string `json:"size"`
	InstanceGroup    string `json:"instanceGroup"`
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
	RouteTableID           string `json:"routeTableId"`
	InternetGatewayID      string `json:"internetGatewayId"`
	NodesSecurityGroupID   string `json:"nodesSecurityGroupID"`
	MastersSecurityGroupID string `json:"mastersSecurityGroupID"`
	MastersInstanceProfile string `json:"mastersInstanceProfile"`
	NodesInstanceProfile   string `json:"nodesInstanceProfile"`
	VolumeSize             string `json:"volumeSize"`
	EbsOptimized           string `json:"ebsOptimized"`
	ImageID                string `json:"image"`
	InstanceType           string `json:"size"`
	HasPublicAddr          bool   `json:"hasPublicAddr"`
	// Map of availability zone to subnet
	Subnets map[string]string `json:"subnets"`
	// Map az to route table association
	RouteTableAssociationIDs map[string]string `json:"routeTableAssociationIds"`
}

type FlannelConfig struct {
	IsMaster bool   `json:"isMaster"`
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
	IsMaster       bool   `json:"isMaster"`
	NodeLabels     string `json:"nodeLabels"`
	ProxyPort      string `json:"proxyPort"`
	K8SVersion     string `json:"k8sVersion"`
	ProviderString string `json:"ProviderString"`
}

type ManifestConfig struct {
	IsMaster            bool   `json:"isMaster"`
	K8SVersion          string `json:"k8sVersion"`
	KubernetesConfigDir string `json:"kubernetesConfigDir"`
	RBACEnabled         bool   `json:"rbacEnabled"`
	ProviderString      string `json:"ProviderString"`
	ServicesCIDR        string `json:"servicesCIDR"`
	ClusterDNSIP        string `json:"clusterDNSIp"`
	MasterHost          string `json:"masterHost"`
	MasterPort          string `json:"masterPort"`
	Password            string `json:"password"`
}

type PostStartConfig struct {
	IsMaster    bool          `json:"isMaster"`
	Host        string        `json:"host"`
	Port        string        `json:"port"`
	Username    string        `json:"username"`
	RBACEnabled bool          `json:"rbacEnabled"`
	Timeout     time.Duration `json:"timeout"`
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
	AdvertiseHost  string        `json:"advertiseHost"`
	Host           string        `json:"host"`
	DataDir        string        `json:"dataDir"`
	ServicePort    string        `json:"servicePort"`
	ManagementPort string        `json:"managementPort"`
	Timeout        time.Duration `json:"timeout"`
	StartTimeout   string        `json:"startTimeout"`
	RestartTimeout string        `json:"restartTimeout"`
	ClusterToken   string        `json:"clusterToken"`
}

type ClusterCheckConfig struct {
	MachineCount int
}

type PrometheusConfig struct {
	Port        string `json:"port"`
	RBACEnabled bool   `json:"rbacEnabled"`
}

type DrainConfig struct {
	PrivateIP string `json:"privateIp"`
}

type Map struct {
	internal map[string]*model.Machine
}

func (m *Map) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.internal)
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.internal)
}

func NewMap(m map[string]*model.Machine) Map {
	return Map{
		internal: m,
	}
}

type Config struct {
	Kube model.Kube `json:"kube"`

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
	PrometheusConfig   PrometheusConfig   `json:"prometheusConfig"`
	DrainConfig        DrainConfig        `json:"drainConfig"`

	ClusterCheckConfig ClusterCheckConfig `json:"clusterCheckConfig"`

	Node             model.Machine `json:"node"`
	CloudAccountID   string        `json:"cloudAccountId" valid:"required, length(1|32)"`
	CloudAccountName string        `json:"cloudAccountName" valid:"required, length(1|32)"`
	Timeout          time.Duration `json:"timeout"`
	Runner           runner.Runner `json:"-"`

	repository storage.Interface `json:"-"`

	m1      sync.RWMutex
	Masters Map `json:"masters"`

	m2    sync.RWMutex
	Nodes Map `json:"nodes"`

	nodeChan      chan model.Machine
	kubeStateChan chan model.KubeState
	configChan    chan *Config

	ReadyForBootstrapLatch *sync.WaitGroup
}

// NewConfig builds instance of config for provisioning
func NewConfig(clusterName, clusterToken, cloudAccountName string, profile profile.Profile) *Config {
	return &Config{
		Kube: model.Kube{
			SSHConfig: model.SSHConfig{
				Port:      "22",
				User:      "root",
				Timeout:   10,
				PublicKey: profile.PublicKey,
			},
		},
		Provider:    profile.Provider,
		ClusterName: clusterName,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		LogBootstrapPrivateKey: profile.LogBootstrapPrivateKey,
		AWSConfig: AWSConfig{
			Region:                 profile.Region,
			AvailabilityZone:       profile.CloudSpecificSettings[clouds.AwsAZ],
			VPCCIDR:                profile.CloudSpecificSettings[clouds.AwsVpcCIDR],
			VPCID:                  profile.CloudSpecificSettings[clouds.AwsVpcID],
			KeyPairName:            profile.CloudSpecificSettings[clouds.AwsKeyPairName],
			MastersSecurityGroupID: profile.CloudSpecificSettings[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:   profile.CloudSpecificSettings[clouds.AwsNodesSecgroupID],
			HasPublicAddr:          true,
		},
		GCEConfig: GCEConfig{
			AvailabilityZone: profile.Zone,
			ImageFamily:      "ubuntu-1604-lts",
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
			Username:            profile.User,
			Password:            profile.Password,
			StaticAuth:          profile.StaticAuth,
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
			// NOTE(stgleb): this is any host by default works on master nodes
			// on worker node this host is changed by any master ip address
			EtcdHost: "0.0.0.0",
		},
		KubeletConfig: KubeletConfig{
			ProxyPort:      "8080",
			K8SVersion:     profile.K8SVersion,
			ProviderString: toCloudProviderOpt(profile.Provider),
		},
		ManifestConfig: ManifestConfig{
			K8SVersion:          profile.K8SVersion,
			KubernetesConfigDir: "/etc/kubernetes",
			RBACEnabled:         profile.RBACEnabled,
			ServicesCIDR:        profile.K8SServicesCIDR,
			ProviderString:      toCloudProviderOpt(profile.Provider),
			MasterHost:          "localhost",
			MasterPort:          "8080",
			Password:            profile.Password,
		},
		PostStartConfig: PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    profile.User,
			RBACEnabled: profile.RBACEnabled,
			Timeout:     time.Minute * 20,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
			RBACEnabled:     profile.RBACEnabled,
		},
		EtcdConfig: EtcdConfig{
			// TODO(stgleb): this field must be changed per node
			Name:           "etcd0",
			Version:        "3.3.10",
			Host:           "0.0.0.0",
			DataDir:        "/var/supergiant/etcd-data",
			ServicePort:    "2379",
			ManagementPort: "2380",
			Timeout:        time.Minute * 20,
			StartTimeout:   "0",
			RestartTimeout: "5",
			ClusterToken:   clusterToken,
		},
		ClusterCheckConfig: ClusterCheckConfig{
			MachineCount: len(profile.NodesProfiles) + len(profile.MasterProfiles),
		},
		PrometheusConfig: PrometheusConfig{
			Port:        "30900",
			RBACEnabled: profile.RBACEnabled,
		},

		Masters: Map{
			internal: make(map[string]*model.Machine, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 30,
		CloudAccountName: cloudAccountName,

		nodeChan:      make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan: make(chan model.KubeState, 2),
		configChan:    make(chan *Config),
	}
}

func NewConfigFromKube(profile *profile.Profile, k *model.Kube) *Config {
	clusterToken := uuid.New()

	cfg := &Config{
		ClusterID:   k.ID,
		Provider:    profile.Provider,
		ClusterName: k.Name,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		LogBootstrapPrivateKey: profile.LogBootstrapPrivateKey,
		AWSConfig: AWSConfig{
			Region:                 profile.Region,
			AvailabilityZone:       k.CloudSpec[clouds.AwsAZ],
			VPCCIDR:                k.CloudSpec[clouds.AwsVpcCIDR],
			VPCID:                  k.CloudSpec[clouds.AwsVpcID],
			KeyPairName:            k.CloudSpec[clouds.AwsKeyPairName],
			Subnets:                k.Subnets,
			MastersSecurityGroupID: k.CloudSpec[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:   k.CloudSpec[clouds.AwsNodesSecgroupID],
			ImageID:                k.CloudSpec[clouds.AwsImageID],
		},
		GCEConfig: GCEConfig{
			AvailabilityZone: profile.Zone,
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
			Username:            profile.User,
			Password:            profile.Password,
			StaticAuth:          profile.StaticAuth,
			CAKey:               k.Auth.CAKey,
			CACert:              k.Auth.CACert,
			AdminCert:           k.Auth.AdminCert,
			AdminKey:            k.Auth.AdminKey,
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
			// NOTE(stgleb): this is any host by default works on master nodes
			// on worker node this host is changed by any master ip address
			EtcdHost: "0.0.0.0",
		},
		KubeletConfig: KubeletConfig{
			ProxyPort:      "8080",
			K8SVersion:     profile.K8SVersion,
			ProviderString: toCloudProviderOpt(profile.Provider),
		},
		ManifestConfig: ManifestConfig{
			K8SVersion:          profile.K8SVersion,
			KubernetesConfigDir: "/etc/kubernetes",
			RBACEnabled:         profile.RBACEnabled,
			ProviderString:      toCloudProviderOpt(profile.Provider),
			MasterHost:          "localhost",
			MasterPort:          "8080",
			Password:            profile.Password,
		},
		PostStartConfig: PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    profile.User,
			RBACEnabled: profile.RBACEnabled,
			Timeout:     time.Minute * 20,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
			RBACEnabled:     profile.RBACEnabled,
		},
		EtcdConfig: EtcdConfig{
			// TODO(stgleb): this field must be changed per node
			Name:           "etcd0",
			Version:        "3.3.10",
			Host:           "0.0.0.0",
			DataDir:        "/var/supergiant/etcd-data",
			ServicePort:    "2379",
			ManagementPort: "2380",
			Timeout:        time.Minute * 20,
			StartTimeout:   "0",
			RestartTimeout: "5",
			ClusterToken:   clusterToken,
		},
		ClusterCheckConfig: ClusterCheckConfig{
			MachineCount: len(profile.NodesProfiles) + len(profile.MasterProfiles),
		},
		PrometheusConfig: PrometheusConfig{
			Port:        "30900",
			RBACEnabled: profile.RBACEnabled,
		},

		Masters: Map{
			internal: make(map[string]*model.Machine, len(k.Masters)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(k.Nodes)),
		},
		Timeout:          time.Minute * 30,
		CloudAccountName: k.AccountName,
		nodeChan:         make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan:    make(chan model.KubeState, 5),
		configChan:       make(chan *Config),
	}

	if k != nil {
		cfg.Kube = *k

		cfg.Kube.SSHConfig = model.SSHConfig{
			Port:      "22",
			User:      "root",
			Timeout:   10,
			PublicKey: profile.PublicKey,
		}
	}

	return cfg
}

// AddMaster to map of master, map is used because it is reference and can be shared among
// goroutines that run multiple tasks of cluster deployment
func (c *Config) AddMaster(n *model.Machine) {
	c.m1.Lock()
	defer c.m1.Unlock()
	c.Masters.internal[n.ID] = n
}

// AddNode to map of nodes in cluster
func (c *Config) AddNode(n *model.Machine) {
	c.m2.Lock()
	defer c.m2.Unlock()
	c.Nodes.internal[n.ID] = n
}

// GetMaster returns first master in master map or nil
func (c *Config) GetMaster() *model.Machine {
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
		if c.Masters.internal[key] != nil && c.Masters.internal[key].State == model.MachineStateActive {
			return c.Masters.internal[key]
		}
	}

	return nil
}

func (c *Config) GetMasters() map[string]*model.Machine {
	c.m1.RLock()
	defer c.m1.RUnlock()

	m := make(map[string]*model.Machine, len(c.Masters.internal))

	for key := range c.Masters.internal {
		m[c.Masters.internal[key].Name] = c.Masters.internal[key]
	}

	return m
}

func (c *Config) GetNodes() map[string]*model.Machine {
	c.m2.RLock()
	defer c.m2.RUnlock()

	m := make(map[string]*model.Machine, len(c.Nodes.internal))

	for key := range c.Nodes.internal {
		m[c.Nodes.internal[key].Name] = c.Nodes.internal[key]
	}

	return m
}

// GetMaster returns first master in master map or nil
func (c *Config) GetNode() *model.Machine {
	c.m2.RLock()
	defer c.m2.RUnlock()

	if len(c.Nodes.internal) == 0 {
		return nil
	}

	for key := range c.Nodes.internal {
		// Skip inactive nodes for selecting
		if c.Nodes.internal[key] != nil && c.Nodes.internal[key].State == model.MachineStateActive {
			return c.Nodes.internal[key]
		}
	}

	return nil
}

func (c *Config) NodeChan() chan model.Machine {
	return c.nodeChan
}

func (c *Config) KubeStateChan() chan model.KubeState {
	return c.kubeStateChan
}

func (c *Config) ConfigChan() chan *Config {
	return c.configChan
}

func (c *Config) SetNodeChan(nodeChan chan model.Machine) {
	c.nodeChan = nodeChan
}

func (c *Config) SetKubeStateChan(kubeStateChan chan model.KubeState) {
	c.kubeStateChan = kubeStateChan
}

func (c *Config) SetConfigChan(configChan chan *Config) {
	c.configChan = configChan
}

// TODO: cloud profiles is deprecated by kubernetes, use controller-managers
func toCloudProviderOpt(cloudName clouds.Name) string {
	switch cloudName {
	case clouds.AWS:
		return "aws"
	case clouds.GCE:
		return "gce"
	}
	return ""
}
