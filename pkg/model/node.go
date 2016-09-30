package model

import "time"

type NodeList struct {
	BaseList
	Items []*Node `json:"items"`
}

type Node struct {
	BaseModel

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" gorm:"not null;index" validate:"nonzero" sg:"immutable"`

	// This is the only input for Node
	Size string `json:"size" validate:"nonzero" sg:"immutable"`

	ProviderID                string    `json:"provider_id" sg:"readonly" gorm:"index"`
	Name                      string    `json:"name" sg:"readonly" gorm:"index"`
	ExternalIP                string    `json:"external_ip" sg:"readonly"`
	ProviderCreationTimestamp time.Time `json:"provider_creation_timestamp" sg:"readonly"`

	OutOfDisk bool `json:"out_of_disk" sg:"readonly"`

	ResourceMetrics
}
