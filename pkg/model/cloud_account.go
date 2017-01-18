package model

type CloudAccountList struct {
	BaseList
	Items []*CloudAccount `json:"items"`
}

type CloudAccount struct {
	BaseModel

	// has_many Kubes
	Kubes []*Kube `json:"kubes,omitempty" gorm:"ForeignKey:CloudAccountName;AssociationForeignKey:Name"`

	Name string `json:"name" validate:"nonzero" gorm:"not null;unique_index" sg:"immutable"`

	Provider string `json:"provider" validate:"regexp=^(aws|digitalocean|openstack|gce)$" gorm:"not null" sg:"immutable"`

	// NOTE this is loose map to allow for multiple clouds (eventually)
	Credentials     map[string]string `json:"credentials,omitempty" validate:"nonzero" gorm:"-" sg:"store_as_json_in=CredentialsJSON,private,immutable"`
	CredentialsJSON []byte            `json:"-" gorm:"not null"`
}
