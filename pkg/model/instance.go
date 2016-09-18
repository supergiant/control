package model

type InstanceList struct {
	BaseList
	Items []*Instance `json:"items"`
}

type Instance struct {
	BaseModel

	// belongs_to Component (we could simply get by current or target ID off of Component, but this makes it simpler for Component-based lookup)
	Component   *Component `json:"component,omitempty"`
	ComponentID *int64     `json:"component_id" gorm:"not null;index"`

	// belongs_to Release (this can be changed)
	Release   *Release `json:"release,omitempty"`
	ReleaseID *int64   `json:"release_id" gorm:"not null;index"`

	// has_many Volumes (for preloading)
	Volumes []*Volume `json:"volumes,omitempty"`

	Num int `json:"num"`

	Name string `json:"name"`

	Started bool `json:"started" gorm:"index"`

	ResourceMetrics

	Addresses     *Addresses `json:"addresses,omitempty" gorm:"-" sg:"store_as_json_in=AddressesJSON"`
	AddressesJSON []byte     `json:"-"`
}
