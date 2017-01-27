package model

type HelmRepoList struct {
	BaseList
	Items []*HelmRepo `json:"items"`
}

type HelmRepo struct {
	BaseModel

	Name string `json:"name" validate:"nonzero" gorm:"not null;unique_index" sg:"immutable"`
	URL  string `json:"url" validate:"nonzero" gorm:"not null" sg:"immutable"`

	// has_many Charts
	Charts []*HelmChart `json:"charts,omitempty" gorm:"ForeignKey:RepoName;AssociationForeignKey:Name"`
}
