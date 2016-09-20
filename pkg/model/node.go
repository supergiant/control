package model

import "time"

type Node struct {
	BaseModel

	// belongs_to Kube
	Kube   *Kube  `json:"kube,omitempty"`
	KubeID *int64 `json:"kube_id" gorm:"not null;index"`

	// This is the only input for Node
	Size string `json:"size" validate:"nonzero"`

	ProviderID                string    `json:"provider_id" sg:"readonly" gorm:"index"`
	Name                      string    `json:"name" sg:"readonly" gorm:"index"`
	ExternalIP                string    `json:"external_ip" sg:"readonly"`
	ProviderCreationTimestamp time.Time `json:"provider_creation_timestamp" sg:"readonly"`

	OutOfDisk bool `json:"out_of_disk" sg:"readonly"`
	Ready     bool `json:"ready" sg:"readonly"`

	ResourceMetrics
}
