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
	// NodeState represents a kubernetes node state.
	NodeState string `json:"nodeState"`
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

	// Userdata is a base64 encoded representation of shell commands or cloud-init directives
	// that applies after the instance starts.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
	Userdata string `json:"userdata"`

	SupergiantV1Config *SupergiantV1UserdataVars
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
