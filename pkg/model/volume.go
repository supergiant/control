package model

type VolumeList struct {
	Pagination
	Items []*Volume `json:"items"`
}

type Volume struct {
	BaseModel

	// belongs_to Instance
	Instance   *Instance `json:"instance"`
	InstanceID *int64    `json:"instance_id" gorm:"not null;index"`

	// belongs_to Kube
	Kube   *Kube  `json:"kube"`
	KubeID *int64 `json:"kube_id" gorm:"not null;index"`

	// NOTE these are the same as VolumeBlueprint (we may want to repeat valiations)
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`

	ProviderID string `json:"provider_id" sg:"readonly"`
}
