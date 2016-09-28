package model

import "encoding/json"

type KubeResourceList struct {
	BaseList
	Items []*KubeResource `json:"items"`
}

type KubeResource struct {
	BaseModel

	// NOTE there is a 4-way unique index on kube_name, namespace, kind, and name.

	// belongs_to Kube
	Kube     *Kube  `json:"kube,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`
	KubeName string `json:"kube_name" validate:"nonzero" gorm:"not null;unique_index:kube_namespace_kind_name"`

	// TODO get actual Kubernetes name regex for validation

	// Kind corresponds directly to the Kind of Kubernetes resource (e.g. Pod, Service, etc.)
	Kind string `json:"kind" validate:"nonzero" gorm:"not null;unique_index:kube_namespace_kind_name"`

	// Namespace corresponds directly to the name of the Kubernetes namespace.
	Namespace string `json:"namespace" validate:"nonzero,max=24" gorm:"not null;unique_index:kube_namespace_kind_name"`

	// Name corresponds directly to the name of the resource in Kubernetes.
	Name string `json:"name" validate:"nonzero,max=24" gorm:"not null;unique_index:kube_namespace_kind_name"`

	// Template is where Kubernetes resources can be defined with special
	// Supergiant "modifiers". If the User is wanting to use the Pod or Service
	// Provisioner (to have Supergiant resources created and applied), Template
	// must be defined in place of Definition.
	Template     *json.RawMessage `json:"template" gorm:"-" sg:"store_as_json_in=TemplateJSON"`
	TemplateJSON []byte           `json:"-"`

	// Definition is where the finalized Kubernetes resource is stored, used to
	// actually create the entity. The User can provide this directly, which will
	// use the DefaultProvisioner.
	Definition     *json.RawMessage `json:"definition" gorm:"-"` // We aren't actually storing this anymore
	DefinitionJSON []byte           `json:"-"`

	// Artifact is where the full resource response from Kubernetes is stored.
	Artifact     *json.RawMessage `json:"artifact" gorm:"-" sg:"store_as_json_in=ArtifactJSON,readonly"`
	ArtifactJSON []byte           `json:"-"`

	// Started represents whether the resource exists in Kubernetes or not. If it
	// is a Pod, it also means the Pod is running.
	Started bool `json:"started" sg:"readonly"`

	// This is used to store unstructured data such as metrics from Heapster.
	ExtraData     map[string]interface{} `json:"extra_data" gorm:"-" sg:"store_as_json_in=ExtraDataJSON,readonly"`
	ExtraDataJSON []byte                 `json:"-"`
}

func (m *KubeResource) SetPassiveStatus() {
	m.PassiveStatusOkay = m.Started
	if m.Started {
		m.PassiveStatus = "started"
		return
	}
	m.PassiveStatus = "stopped"
}
