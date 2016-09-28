package model

type EntrypointList struct {
	BaseList
	Items []*Entrypoint `json:"items"`
}

type Entrypoint struct {
	BaseModel

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" gorm:"not null;index"`

	Name string `json:"name" validate:"nonzero,max=21,regexp=^[\\w-]+$" gorm:"not null;unique_index"`

	Listeners     []*EntrypointListener `json:"ports,omitempty" gorm:"-" sg:"store_as_json_in=ListenersJSON"`
	ListenersJSON []byte                `json:"-"`

	ProviderID string `json:"provider_id" sg:"readonly"`
	Address    string `json:"address,omitempty" sg:"readonly"`
}

func (m *Entrypoint) BeforeCreate() error {
	m.ProviderID = "sg-" + m.Name
	return nil
}
