package model

type VolumeList struct {
	BaseList
	Items []*Volume `json:"items"`
}

type Volume struct {
	BaseModel

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" validate:"nonzero" gorm:"not null;index" sg:"immutable"`

	// belongs_to KubeResource (this can change, or even be temporarily nil)
	KubeResource   *KubeResource `json:"kube_resource,omitempty"`
	KubeResourceID *int64        `json:"kube_resource_id" gorm:"index" sg:"readonly"`

	Name string `json:"name" validate:"nonzero,max=24,regexp=^[\\w-]+$" gorm:"not null;unique_index" sg:"immutable"`
	Type string `json:"type" sg:"immutable"`
	Size int    `json:"size" validate:"nonzero"`

	ProviderID string `json:"provider_id" sg:"readonly"`
}
