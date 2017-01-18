package kubernetes

import "time"

type HeapsterStats struct {
	Name     string `json:"name"`
	CPUUsage int64  `json:"cpuUsage"`
	RAMUsage int64  `json:"memUsage"`
}

type HeapsterMetric struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

type HeapsterMetrics struct {
	Metrics []*HeapsterMetric `json:"metrics"`
}

type Metadata struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
}

//------------------------------------------------------------------------------
type NamespaceList struct {
	Items []*Namespace `json:"items"`
}

type Namespace struct {
	Metadata Metadata `json:"metadata"`
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
type NodeList struct {
	Items []*Node `json:"items"`
}

type Node struct {
	Metadata Metadata   `json:"metadata"`
	Spec     NodeSpec   `json:"spec"`
	Status   NodeStatus `json:"status"`
}

type NodeSpec struct {
	ExternalID string `json:"externalID"`
}

type NodeStatusCapacity struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type NodeStatusCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type NodeStatus struct {
	Capacity   NodeStatusCapacity    `json:"capacity"`
	Conditions []NodeStatusCondition `json:"conditions"`
	Addresses  []NodeAddress         `json:"addresses"`
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
type PodList struct {
	Items []*Pod `json:"items"`
}

type Pod struct {
	Metadata Metadata  `json:"metadata"`
	Spec     PodSpec   `json:"spec"`
	Status   PodStatus `json:"status"`
}

type AwsElasticBlockStore struct {
	VolumeID string `json:"volumeID"`
	FSType   string `json:"fsType"`
}

type FlexVolume struct {
	Driver  string            `json:"driver"`
	FSType  string            `json:"fsType"`
	Options map[string]string `json:"options"`
}

type Cinder struct {
	VolumeID string `json:"volumeID"`
	FSType   string `json:"fsType"`
}

type GcePersistentDisk struct {
	PDName string `json:"pdName"`
	FSType string `json:"fsType"`
}

type Volume struct {
	Name                 string                `json:"name"`
	AwsElasticBlockStore *AwsElasticBlockStore `json:"awsElasticBlockStore,omitempty"`
	FlexVolume           *FlexVolume           `json:"flexVolume,omitempty"`
	Cinder               *Cinder               `json:"cinder,omitempty"`
	GcePersistentDisk    *GcePersistentDisk    `json:"gcePersistentDisk,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type ResourceValues struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

type Resources struct {
	Limits   ResourceValues `json:"limits"`
	Requests ResourceValues `json:"requests"`
}

type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SecurityContext struct {
	Privileged bool `json:"privileged"`
}

type Container struct {
	Name            string          `json:"name"`
	Image           string          `json:"image"`
	Command         []string        `json:"command"`
	Resources       Resources       `json:"resources"`
	Ports           []ContainerPort `json:"ports"`
	VolumeMounts    []VolumeMount   `json:"volumeMounts"`
	Env             []EnvVar        `json:"env"`
	SecurityContext SecurityContext `json:"securityContext"`
	ImagePullPolicy string          `json:"imagePullPolicy"`
}

type ImagePullSecret struct {
	Name string `json:"name"`
}

type PodSpec struct {
	Volumes                       []Volume          `json:"volumes"`
	Containers                    []Container       `json:"containers"`
	ImagePullSecrets              []ImagePullSecret `json:"imagePullSecrets"`
	TerminationGracePeriodSeconds int               `json:"terminationGracePeriodSeconds"`
	RestartPolicy                 string            `json:"restartPolicy"`
	NodeName                      string            `json:"nodeName"`
}

type ContainerStateRunning struct {
	StartedAt string `json:"startedAt"` // TODO should be time type
}

type ContainerStateTerminated struct {
	ExitCode   int    `json:"exitcode"`
	StartedAt  string `json:"startedAt"`  // TODO should be time type
	FinishedAt string `json:"finishedAt"` // TODO should be time type
	Reason     string `json:"reason"`
}

type ContainerState struct {
	Running    ContainerStateRunning    `json:"running"`
	Terminated ContainerStateTerminated `json:"terminated"`
}

type ContainerStatus struct {
	ContainerID  string         `json:"containerID"`
	Image        string         `json:"image"`
	ImageID      string         `json:"imageID"`
	Name         string         `json:"name"`
	Ready        bool           `json:"ready"`
	RestartCount int            `json:"restartCount"`
	State        ContainerState `json:"state"`
	LastState    ContainerState `json:"state"`
}

type PodStatusCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type PodStatus struct {
	Phase             string               `json:"phase"`
	Conditions        []PodStatusCondition `json:"conditions"`
	ContainerStatuses []ContainerStatus    `json:"containerStatuses"`
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
type EventList struct {
	Items []*Event `json:"items"`
}

type Event struct {
	Metadata Metadata `json:"metadata"`
	Message  string   `json:"message"`
	Count    int      `json:"count"`
	Source   Source   `json:"source"`
}

type Source struct {
	Host string `json:"host"`
}
