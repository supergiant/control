package model

type EntrypointListenerList struct {
	BaseList
	Items []*EntrypointListener `json:"items"`
}

type EntrypointListener struct {
	BaseModel

	// belongs_to Entrypoint
	Entrypoint     *Entrypoint `json:"entrypoint,omitempty" gorm:"ForeignKey:EntrypointName;AssociationForeignKey:Name"`
	EntrypointName string      `json:"entrypoint_name" validate:"nonzero" gorm:"not null;index;unique_index:entrypoint_port;unique_index:port_name_within_entrypoint" sg:"immutable"`

	// belongs_to KubeResource
	KubeResource   *KubeResource `json:"kube_resource,omitempty"`
	KubeResourceID *int64        `json:"kube_resource_id" gorm:"index" sg:"readonly"`

	// Name is required to give the port identity
	Name string `json:"name" validate:"nonzero,max=24,regexp=^[a-z]([-a-z0-9]*[a-z0-9])?$" gorm:"not null;index;unique_index:port_name_within_entrypoint" sg:"immutable"`

	// EntrypointPort is the external port the user connects to
	EntrypointPort     int64  `json:"entrypoint_port" validate:"nonzero" gorm:"not null;unique_index:entrypoint_port" sg:"immutable"`
	EntrypointProtocol string `json:"entrypoint_protocol" validate:"nonzero" sg:"default=TCP,immutable"`

	// NodePort is the target port, what EntrypointPort maps to
	NodePort     int64  `json:"node_port" validate:"nonzero" sg:"immutable"`
	NodeProtocol string `json:"node_protocol" validate:"nonzero" sg:"default=TCP,immutable"`
}
