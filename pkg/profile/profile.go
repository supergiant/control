package profile

import "github.com/supergiant/control/pkg/clouds"

type Profile struct {
	ID string `json:"id" valid:"required"`

	MasterProfiles []NodeProfile `json:"masterProfiles" valid:"-"`
	NodesProfiles  []NodeProfile `json:"nodesProfiles" valid:"-"`

	// StaticAuth represents tokens and basic authentication credentials that
	// would be set to kube-apiserver on start.
	StaticAuth StaticAuth `json:"staticAuth" valid:"-"`

	// TODO: get rid of its usage
	// DEPRECATED: should be a part of the static auth
	User string `json:"user" valid:"-"`
	// DEPRECATED: should be a part of the static auth
	Password string `json:"password" valid:"-"`

	// TODO(stgleb): In future releases arch will probably migrate to node profile
	// to allow user create heterogeneous cluster of machine with different arch
	Provider        clouds.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)" valid:"-"`
	Region          string      `json:"region" valid:"-"`
	Zone            string      `json:"zone" valid:"-"`
	Arch            string      `json:"arch" valid:"-"`
	OperatingSystem string      `json:"operatingSystem" valid:"-"`
	UbuntuVersion   string      `json:"ubuntuVersion" valid:"-"`
	DockerVersion   string      `json:"dockerVersion" valid:"-"`
	K8SVersion      string      `json:"K8SVersion" valid:"-"`
	K8SServicesCIDR string      `json:"k8sServicesCIDR" valid:"-"`
	K8SAPIPort      int       `json:"k8sApiPort" valid:"-"`
	NetworkProvider string      `json:"networkProvider" valid:"-"`
	FlannelVersion  string      `json:"flannelVersion" valid:"-"`
	NetworkType     string      `json:"networkType" valid:"-"`
	CIDR            string      `json:"cidr" valid:"-"`
	HelmVersion     string      `json:"helmVersion" valid:"-"`
	RBACEnabled     bool        `json:"rbacEnabled" valid:"-"`
	// This field is AWS specific, mapping AZ -> subnet
	Subnets               map[string]string     `json:"subnets" valid:"-"`
	CloudSpecificSettings CloudSpecificSettings `json:"cloudSpecificSettings" valid:"-"`
	PublicKey             string                `json:"publicKey" valid:"-"`

	// ExposedAddresses is a list of cidr/port pairs that will be exposes
	// by cloud provider security groups.
	ExposedAddresses []Addresses `json:"exposedAddresses" valid:"-"`
	Addons           []string    `json:"addons,omitempty" valid:"-"`
}

type NodeProfile map[string]string
type CloudSpecificSettings map[string]string

// Addresses uses cidr to define an ip list.
type Addresses struct {
	CIDR string `json:"cidr"`
}

// StaticAuth represents tokens and basic authentication credentials.
type StaticAuth struct {
	BasicAuth []BasicAuthUser `json:"basicAuth"`
	Tokens    []TokenAuthUser `json:"tokens"`
}

// BasicAuthUser represents an entry of the static password file.
// https://kubernetes.io/docs/reference/access-authn-authz/authentication/#static-password-file
type BasicAuthUser struct {
	Password string   `json:"password"`
	Name     string   `json:"name"`
	ID       string   `json:"id"`
	Groups   []string `json:"groups"`
}

// BasicAuthUser represents an entry of the static token file.
// https://kubernetes.io/docs/reference/access-authn-authz/authentication/#static-token-file
type TokenAuthUser struct {
	Token  string   `json:"token"`
	Name   string   `json:"name"`
	ID     string   `json:"id"`
	Groups []string `json:"groups"`
}
