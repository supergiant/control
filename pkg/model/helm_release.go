package model

import "strings"

type HelmReleaseList struct {
	BaseList
	Items []*HelmRelease `json:"items"`
}

type HelmRelease struct {
	BaseModel

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" gorm:"not null;index" validate:"nonzero" sg:"immutable"`

	// NOTE this is just a "soft" belongs_to.
	// We don't do relation since there's no real need, and it complicates things.
	RepoName     string `json:"repo_name" gorm:"not null;index" validate:"nonzero" sg:"immutable"`
	ChartName    string `json:"chart_name" gorm:"not null;index" validate:"nonzero" sg:"immutable"`
	ChartVersion string `json:"chart_version" validate:"nonzero" sg:"immutable"`

	Name string `json:"name" validate:"regexp=^[\\w-\\.]*$" gorm:"index" sg:"immutable"`

	Revision string `json:"revision"`
	// TODO weird naming, but Status is already taken
	StatusValue  string `json:"status_value"`
	UpdatedValue string `json:"updated_value"`

	Config     map[string]interface{} `json:"config" gorm:"-" sg:"store_as_json_in=ConfigJSON,immutable"`
	ConfigJSON []byte                 `json:"-"`
}

func (m *HelmRelease) SetPassiveStatus() {
	m.PassiveStatus = strings.ToLower(m.StatusValue)
	m.PassiveStatusOkay = m.StatusValue == "DEPLOYED"
}
