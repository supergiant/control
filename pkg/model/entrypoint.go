package model

type EntrypointList struct {
	BaseList
	Items []*Entrypoint `json:"items"`
}

type Entrypoint struct {
	BaseModel

	// belongs_to Kube
	Kube   *Kube  `json:"kube,omitempty"`
	KubeID *int64 `json:"kube_id" gorm:"not null;index"`

	Name string `json:"name" validate:"nonzero,max=21,regexp=^[\\w-]+$" gorm:"not null;unique_index"`

	ProviderID string `json:"provider_id" sg:"readonly"`

	// the ELB address
	Address string `json:"address,omitempty" sg:"readonly"`
}

func (m *Entrypoint) BeforeCreate() error {
	m.ProviderID = "sg-" + m.Name
	return nil
}
