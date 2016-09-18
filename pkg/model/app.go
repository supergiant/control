package model

type AppList struct {
	BaseList
	Items []*App `json:"items"`
}

type App struct {
	BaseModel
	Name string `json:"name" validate:"nonzero,max=24,regexp=^[a-z]([-a-z0-9]*[a-z0-9])?$" gorm:"not null;unique_index:name_within_kube"`

	// belongs_to Kube
	Kube   *Kube  `json:"kube,omitempty"`
	KubeID *int64 `json:"kube_id" gorm:"not null;index;unique_index:name_within_kube"`

	// has_many Components
	Components []*Component `json:"components,omitempty"`
}
