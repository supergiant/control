package model

type LoadBalancerList struct {
	BaseList
	Items []*LoadBalancer `json:"items"`
}

type LoadBalancer struct {
	BaseModel

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" gorm:"not null;unique_index:name_within_kube" validate:"nonzero" sg:"immutable"`

	Name string `json:"name" validate:"nonzero,max=63" gorm:"not null;unique_index:name_within_kube" sg:"immutable"`

	Namespace string `json:"namespace" validate:"nonzero" sg:"immutable"`

	Selector     map[string]string `json:"selector,omitempty" gorm:"-" sg:"store_as_json_in=SelectorJSON"`
	SelectorJSON []byte            `json:"-"`

	Ports     map[int]int `json:"ports,omitempty" gorm:"-" sg:"store_as_json_in=PortsJSON"`
	PortsJSON []byte      `json:"-"`

	Address string `json:"address" sg:"readonly"`
}
