package model

type HelmChartList struct {
	BaseList
	Items []*HelmChart `json:"items"`
}

type HelmChart struct {
	BaseModel

	// belongs_to Repo
	Repo     *HelmRepo `json:"repo,omitempty" gorm:"ForeignKey:RepoName;AssociationForeignKey:Name"`
	RepoName string    `json:"repo_name" validate:"nonzero" gorm:"not null;index" sg:"immutable"`

	Name        string `json:"name" validate:"nonzero" gorm:"not null;index" sg:"immutable"`
	Version     string `json:"version" validate:"nonzero"`
	Description string `json:"description" validate:"nonzero"`

	DefaultConfig     map[string]interface{} `json:"default_config" gorm:"-" sg:"store_as_json_in=DefaultConfigJSON,immutable"`
	DefaultConfigJSON []byte                 `json:"-"`
}
