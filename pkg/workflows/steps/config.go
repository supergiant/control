package steps

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
)

const (
	DefaultK8SAPIPort int64 = 443
)

type DOConfig struct {
	Name string `json:"name" valid:"required"`
	// These come from UI select
	Region string `json:"region" valid:"required"`
	Size   string `json:"size" valid:"required"`
	Image  string `json:"image" valid:"required"`

	// These come from cloud account
	Fingerprint string `json:"fingerprint" valid:"required"`
	AccessToken string `json:"accessToken" valid:"required"`

	ExternalLoadBalancerID string `json:"externalLoadBalancerId"`
	InternalLoadBalancerID string `json:"internalLoadBalancerId"`
}

type ServiceAccount struct {
	// NOTE(stgleb): This comes from cloud account
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`

	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`

	ClientEmail string `json:"client_email"`
	ClientID    string `json:"client_id"`

	AuthURI  string `json:"auth_uri"`
	TokenURI string `json:"token_uri"`

	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type GCEConfig struct {
	ServiceAccount

	// This comes from profile
	ImageFamily      string `json:"imageFamily"`
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone"`
	Size             string `json:"size"`

	NetworkName string `json:"networkName"`
	NetworkLink string `json:"networkLink"`
	SubnetLink  string `json:"subnetLink"`

	AZs map[string]string `json:"subnets"`

	// Mapping AZ -> Instance group
	InstanceGroupNames map[string]string `json:"instanceGroupNames"`
	InstanceGroupLinks map[string]string `json:"instanceGroupLinks"`

	// Target pool acts as a balancer for external traffic https://cloud.google.com/load-balancing/docs/target-pools
	TargetPoolName string `json:"targetPoolName"`
	TargetPoolLink string `json:"targetPoolLink"`

	ExternalAddressName string `json:"externalAddressName"`
	InternalAddressName string `json:"internalAddressName"`

	ExternalIPAddressLink string `json:"externalIpAddressLink"`
	InternalIPAddressLink string `json:"internalIpAddressLink"`

	BackendServiceName string `json:"backendServiceName"`
	BackendServiceLink string `json:"backendServiceLink"`

	HealthCheckName string `json:"healthCheckName"`

	ExternalForwardingRuleName string `json:"externalForwardingRuleName"`
	InternalForwardingRuleName string `json:"externalForwardingRuleName"`
}

type AzureConfig struct {
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`

	Location string `json:"location"`

	VMSize     string `json:"vmSize"`
	VolumeSize string `json:"volumeSize"`
	// TODO: cidr validation?
	VNetCIDR string `json:"vNetCIDR"`
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
	DeviceName             string `json:"deviceName"`
	EbsOptimized           string `json:"ebsOptimized"`
	ImageID                string `json:"image"`
	InstanceType           string `json:"size"`

	ExternalLoadBalancerName string `json:"externalLoadBalancerName"`
	InternalLoadBalancerName string `json:"internalLoadBalancerName"`

	// Map of availability zone to subnet
	Subnets map[string]string `json:"subnets"`
	// Map az to route table association
	RouteTableAssociationIDs map[string]string `json:"routeTableAssociationIds"`
}

type DrainConfig struct {
	PrivateIP string `json:"privateIp"`
}

type ApplyConfig struct {
	Data string `json:"data"`
}

type InstallAppConfig struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	ChartName    string `json:"chartName" valid:"required"`
	ChartVersion string `json:"chartVersion"`
	RepoName     string `json:"repoName" valid:"required"`
	ChartRef     string `json:"chartRef"`
	Values       string `json:"values"`
}

type OpenStackConfig struct {
	AuthURL          string `json:"authUrl"`
	CACert           string `json:"caCert"`
	DomainName       string `json:"domainName"`
	DomainID         string `json:"domainId"`
	TenantID         string `json:"tenantId"`
	TenantName       string `json:"tenantName"`
	UserName         string `json:"userName"`
	Password         string `json:"password"`
	Region           string `json:"region"`
	ImageID          string `json:"imageId"`
	NetworkID        string `json:"networkId"`
	NetworkName      string `json:"networkName"`
	SubnetID         string `json:"subnetId"`
	SubnetIPRange    string `json:"subnetIpRange"`
	RouterID         string `json:"routerId"`
	FlavorName       string `json:"flavorName"`
	FloatingIP       string `json:"floatingIp"`
	FloatingID       string `json:"floatingId"`
	ImageName        string `json:"imageName"`
	LoadBalancerID   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	ListenerID       string `json:"listenerId"`
	PoolID           string `json:"poolId"`
	HealthCheckID    string `json:"healthCheckId"`
	KeyPairName      string `json:"keyPairName"`
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

	DryRun             bool `json:"dryRun"`
	TaskID             string
	IsMaster           bool         `json:"isMaster"`
	IsBootstrap        bool         `json:"IsBootstrap"`
	IsImport           bool         `json:"isImport"`
	DigitalOceanConfig DOConfig     `json:"digitalOceanConfig"`
	AWSConfig          AWSConfig    `json:"awsConfig"`
	GCEConfig          GCEConfig    `json:"gceConfig"`
	AzureConfig        AzureConfig  `json:"azureConfig"`
	OSConfig           OSConfig     `json:"osConfig"`
	PacketConfig       PacketConfig `json:"packetConfig"`

	DrainConfig      DrainConfig      `json:"drainConfig"`
	ConfigMap        ConfigMap        `json:"configMap"`
	ApplyConfig      ApplyConfig      `json:"applyConfig"`
	InstallAppConfig InstallAppConfig `json:"installAppConfig"`
	OpenStackConfig  OpenStackConfig  `json:"openStackConfig"`

	Provider clouds.Name `json:"provider"`

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

	authorizerMux  sync.RWMutex
	azureAthorizer autorest.Authorizer

	nodeChan      chan model.Machine
	kubeStateChan chan model.KubeState
	configChan    chan *Config
}

type ConfigMap struct {
	Data string
}

// NewConfig builds instance of config for provisioning
func NewConfig(clusterName, cloudAccountName string, profile profile.Profile) (*Config, error) {
	if err := validateAddons(profile.Addons); err != nil {
		return nil, err
	}

	var user = "root"

	if profile.Provider == clouds.AWS {
		//on aws default user name on ubuntu images are not root but ubuntu
		//https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html
		// TODO: this should be set by provisioner
		user = "ubuntu"
	} else if profile.Provider == clouds.Azure {
		user = clouds.OSUser
	}

	return &Config{
		Kube: model.Kube{
			Name:       clusterName,
			K8SVersion: profile.K8SVersion,
			SSHConfig: model.SSHConfig{
				Port:      "22",
				User:      user,
				Timeout:   30,
				PublicKey: profile.PublicKey,
			},
			Auth: model.Auth{
				Username:   profile.User,
				Password:   profile.Password,
				StaticAuth: profile.StaticAuth,
			},
			Networking: model.Networking{
				Manager:  profile.NetworkProvider,
				Provider: profile.NetworkProvider,
				Type:     profile.NetworkType,
				CIDR:     profile.CIDR,
			},
			Arch:             profile.Arch,
			OperatingSystem:  profile.OperatingSystem,
			DockerVersion:    profile.DockerVersion,
			HelmVersion:      profile.HelmVersion,
			ExposedAddresses: profile.ExposedAddresses,
			APIServerPort:    ensurePort(profile.K8SAPIPort),
			Provider:         profile.Provider,
			RBACEnabled:      profile.RBACEnabled,
			ServicesCIDR:     profile.K8SServicesCIDR,
			Addons:           profile.Addons,
		},
		Provider: profile.Provider,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		AWSConfig: AWSConfig{
			Region:                 profile.Region,
			AvailabilityZone:       profile.CloudSpecificSettings[clouds.AwsAZ],
			VPCCIDR:                profile.CloudSpecificSettings[clouds.AwsVpcCIDR],
			VPCID:                  profile.CloudSpecificSettings[clouds.AwsVpcID],
			KeyPairName:            profile.CloudSpecificSettings[clouds.AwsKeyPairName],
			MastersSecurityGroupID: profile.CloudSpecificSettings[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:   profile.CloudSpecificSettings[clouds.AwsNodesSecgroupID],
			// TODO(stgleb): Passs this from UI or figure out any better way
			DeviceName: "/dev/sda1",
		},
		GCEConfig: GCEConfig{
			AvailabilityZone:   profile.Zone,
			ImageFamily:        "ubuntu-1604-lts",
			Region:             profile.Region,
			InstanceGroupLinks: make(map[string]string),
			InstanceGroupNames: make(map[string]string),
			AZs:                make(map[string]string),
		},
		AzureConfig: AzureConfig{
			Location: profile.Region,
			VNetCIDR: profile.CloudSpecificSettings[clouds.AzureVNetCIDR],
			// TODO(stgleb): this should be passed from the UI
			VolumeSize: "30",
		},

		Masters: Map{
			internal: make(map[string]*model.Machine, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 60,
		CloudAccountName: cloudAccountName,

		nodeChan:      make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan: make(chan model.KubeState, 2),
		configChan:    make(chan *Config),
	}, nil
}

// TODO(stgleb): Compare that to LoadCloudSpecificDataFromKube
func NewConfigFromKube(profile *profile.Profile, k *model.Kube) (*Config, error) {
	if k == nil {
		return nil, errors.Wrapf(sgerrors.ErrNilEntity, "kube must not be nil")
	}

	var user string

	if profile.Provider == clouds.AWS {
		//on aws default user name on ubuntu images are not root but ubuntu
		//https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html
		// TODO: this should be set by provisioner
		user = "ubuntu"
	} else if profile.Provider == clouds.Azure {
		user = clouds.OSUser
	}

	cfg := &Config{
		Provider: profile.Provider,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		Kube: *k,
		AWSConfig: AWSConfig{
			Region:                   profile.Region,
			AvailabilityZone:         k.CloudSpec[clouds.AwsAZ],
			VPCCIDR:                  k.CloudSpec[clouds.AwsVpcCIDR],
			VPCID:                    k.CloudSpec[clouds.AwsVpcID],
			KeyPairName:              k.CloudSpec[clouds.AwsKeyPairName],
			Subnets:                  k.Subnets,
			MastersSecurityGroupID:   k.CloudSpec[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:     k.CloudSpec[clouds.AwsNodesSecgroupID],
			ImageID:                  k.CloudSpec[clouds.AwsImageID],
			ExternalLoadBalancerName: k.CloudSpec[clouds.AwsExternalLoadBalancerName],
			InternalLoadBalancerName: k.CloudSpec[clouds.AwsInternalLoadBalancerName],
			// TODO(stgleb): Passs this from UI or figure out any better way
			DeviceName: "/dev/sda1",
		},
		GCEConfig: GCEConfig{
			AvailabilityZone: profile.Zone,
		},
		AzureConfig: AzureConfig{
			Location:   profile.Region,
			VNetCIDR:   k.CloudSpec[clouds.AzureVNetCIDR],
			VolumeSize: k.CloudSpec[clouds.AzureVolumeSize],
		},
		Masters: Map{
			internal: make(map[string]*model.Machine, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 60,
		CloudAccountName: k.AccountName,
		nodeChan:         make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan:    make(chan model.KubeState, 5),
		configChan:       make(chan *Config),
	}

	// Restore all masters and workers from kube
	for index := range k.Masters {
		cfg.AddMaster(k.Masters[index])
	}

	for index := range k.Nodes {
		cfg.AddNode(k.Nodes[index])
	}

	cfg.Kube = *k

	cfg.Kube.SSHConfig = model.SSHConfig{
		Port:      "22",
		User:      user,
		Timeout:   10,
		PublicKey: profile.PublicKey,
	}

	return cfg, nil
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
	if c.IsMaster && c.Node.State == model.MachineStateActive {
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

func (c *Config) SetAzureAuthorizer(a autorest.Authorizer) {
	c.authorizerMux.Lock()
	defer c.authorizerMux.Unlock()

	c.azureAthorizer = a
}

func (c *Config) GetAzureAuthorizer() autorest.Authorizer {
	c.authorizerMux.RLock()
	defer c.authorizerMux.RUnlock()

	return c.azureAthorizer
}

func ensurePort(p int64) int64 {
	if p == 0 {
		return DefaultK8SAPIPort
	}
	return p
}

func validateAddons(in []string) error {
	invalid := make([]string, 0)
	for _, addon := range in {
		for _, registered := range []string{"dashboard"} {
			if addon != registered {
				invalid = append(invalid, addon)
			}
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("validate addons: unknown: %v", invalid)
	}
	return nil
}
