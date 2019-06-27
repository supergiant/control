package api

import (
	"errors"
	"time"
)

const (
	DefaultConfigMapName      = "capacity"
	DefaultConfigMapNamespace = "kube-system"
	DefaultConfigMapKey       = "kubescaler.conf"
)

// Worker is an abstraction used by kubescaler to manage cluster capacity.
// It contains data from a (virtual) machine and a kubernetes node running on it.
type Worker struct {
	// ClusterName is a kubernetes cluster name.
	ClusterName string `json:"clusterName"`
	// MachineID is a unique id of the provider's virtual machine.
	// required: true
	MachineID string `json:"machineID"`
	// MachineName is a human-readable name of virtual machine.
	MachineName string `json:"machineName"`
	// MachineType is type of virtual machine (eg. 't2.micro' for AWS).
	MachineType string `json:"machineType"`
	// MachineState represent a virtual machine state.
	MachineState string `json:"machineState"`
	// CreationTimestamp is a timestamp representing a time when this machine was created.
	CreationTimestamp time.Time `json:"creationTimestamp"`
	// Reserved is a parameter that is used to prevent downscaling of the worker.
	Reserved bool `json:"reserved"`
	// NodeName represents a name of the kubernetes node that runs on top of that machine.
	NodeName string `json:"nodeName"`
	// NodeLabels represents a labels of the kubernetes node that runs on top of that machine.
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`
}

type WorkerList struct {
	Items []*Worker `json:"items"`
}

type Config struct {
	ClusterName     string            `json:"clusterName"`
	ProviderName    string            `json:"providerName"`
	Provider        map[string]string `json:"provider"`
	Paused          *bool             `json:"paused,omitempty"`
	PauseLock       bool              `json:"pauseLock"`
	ScanInterval    string            `json:"scanInterval"`
	WorkersCountMin int               `json:"workersCountMin"`
	WorkersCountMax int               `json:"workersCountMax"`
	MachineTypes    []string          `json:"machineTypes"`
	// TODO: this is hardcoded and isn't used at the moment
	MaxMachineProvisionTime string            `json:"maxMachineProvisionTime"`
	IgnoredNodeLabels       map[string]string `json:"ignoredNodeLabels"`
	NewNodeTimeBuffer       int               `json:"newNodeTimeBuffer"`

	// Userdata is a base64 encoded representation of shell commands or cloud-init directives
	// that applies after the instance starts.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
	Userdata string `json:"userdata"`

	// UserdataTpl is a template to build Userdata dynamically. UserdataVars is used for template
	// configuration.
	// UserdataTpl will be ignored if Userdata is set.
	UserdataTpl string `json:"userdataTpl"`

	// UserdataVars is a configuration used for parsing a UserdataTpl template.
	UserdataVars map[string]string `json:"userdataVars"`

	SupergiantV1Config *SupergiantV1UserdataVars

	// DEPRECATED: moved to the SupergiantV1Config
	MasterPrivateAddr string `json:"masterPrivateAddr"`
	// DEPRECATED: moved to the SupergiantV1Config
	KubeAPIHost string `json:"kubeAPIHost"`
	// DEPRECATED: moved to the SupergiantV1Config
	KubeAPIPort string `json:"kubeAPIPort"`
	// DEPRECATED: moved to the SupergiantV1Config
	KubeAPIUser string `json:"kubeAPIUser"`
	// DEPRECATED: moved to the SupergiantV1Config
	KubeAPIPassword string `json:"kubeAPIPassword"`
	// DEPRECATED: moved to the SupergiantV1Config
	SSHPubKey string `json:"sshPubKey"`
}

type SupergiantV1UserdataVars struct {
	MasterPrivateAddr string `json:"masterPrivateAddr"`
	KubeAPIHost       string `json:"kubeAPIHost"`
	KubeAPIPort       string `json:"kubeAPIPort"`
	KubeAPIUser       string `json:"kubeAPIUser"`
	KubeAPIPassword   string `json:"kubeAPIPassword"`
	SSHPubKey         string `json:"sshPubKey"`
	KubeVersion       string `json:"-"`
}

func (c Config) Validate() error {
	// TODO: pass it with a pointer or use the ConfigRequest struct for patches.
	if c.WorkersCountMin < 0 {
		return errors.New("WorkersCountMin can't be negative")
	}
	if c.WorkersCountMax < 0 {
		return errors.New("WorkersCountMax can't be negative")
	}
	return nil
}
