package model

import "regexp"

type ReleaseList struct {
	BaseList
	Items []*Release `json:"items"`
}

// NOTE the word Blueprint is used for Volumes and Containers, since they are
// both "definitions" that create "instances" of the real thing
type Release struct {
	BaseModel

	// belongs_to Component
	Component   *Component `json:"component,omitempty"`
	ComponentID *int64     `json:"component_id" gorm:"not null;index"`

	// TODO
	// InstanceGroup is used as a labeling mechanism for instances. If nil,
	// InstanceGroup is set equal to the release's Timestamp. If a value is
	// supplied by the user, it MUST be the current (previous) Release's
	// Timestamp.
	//
	// The purpose of InstanceGroup is to prevent restarting between Releases.
	// I'm pretty sure the ONLY scenario in which this value makes sense is when
	// changing InstanceCount, and the value supplied in such a scenario must be
	// the previous Release's timestamp.
	//
	// It seems as though we need to break deploys up into:
	//   - config changes
	//	 - adding/removing instances
	//
	// It might still makes sense to have Release as a grouping mechanism, though,
	// because it could allow for grouping metrics recorded per-Release. However,
	// you may be able to separate the operations, but have every operation create
	// a new record that records the config / instance count at that time.
	InstanceGroup *int64 `json:"instance_group,omitempty"`

	InstanceCount int `json:"instance_count" validate:"min=1" sg:"default=1"`

	// Config is stored as JSON (like an embed in Mongo)
	Config     *ComponentConfig `json:"config" gorm:"-" sg:"store_as_json_in=ConfigJSON"`
	ConfigJSON []byte           `json:"-" gorm:"not null"`

	// NOTE these are not a traditional relation. They are pulled out of the Config
	// PrivateImageKeys []*PrivateImageKey `json:"private_image_keys" gorm:"-"`

	InUse bool `json:"in_use" sg:"readonly"`
}

type ComponentConfig struct {
	// These attributes, when changed from last Release, indicate a restart is
	// needed (or just new instances through other means).
	Volumes                []*VolumeBlueprint    `json:"volumes"`
	Containers             []*ContainerBlueprint `json:"containers" validate:"min=1"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"min=0" sg:"default=10"`
}

type VolumeBlueprint struct {
	Name string `json:"name" validate:"nonzero,regexp=^\\w[-\\w\\.]*$/"` // TODO max length
	Type string `json:"type" validate:"regexp=^(gp2)$" sg:"default=gp2"` // TODO support other vol types
	Size int    `json:"size" validate:"min=1"`
}

type ContainerBlueprint struct {
	Image      string      `json:"image" validate:"nonzero,regexp=^[-\\w\\.\\/]+(:[-\\w\\.]+)?$"`
	Name       string      `json:"name,omitempty" validate:"regexp=^[-\\w\\.\\/]+(:[-\\w\\.]+)?$"`
	Command    []string    `json:"command,omitempty"`
	Ports      []*Port     `json:"ports,omitempty"`
	Env        []*EnvVar   `json:"env,omitempty"`
	CPURequest *CoresValue `json:"cpu_request"`
	CPULimit   *CoresValue `json:"cpu_limit"`
	RAMRequest *BytesValue `json:"ram_request"`
	RAMLimit   *BytesValue `json:"ram_limit"`
	Mounts     []*Mount    `json:"mounts,omitempty"`
}

func (c *ContainerBlueprint) NameOrDefault() string {
	if c.Name != "" {
		return c.Name
	}
	return regexp.MustCompile("[^A-Za-z0-9]").ReplaceAllString(c.Image, "-")
}

type EnvVar struct {
	Name  string `json:"name" validate:"nonzero"`
	Value string `json:"value" validate:"nonzero"` // this may be templated, "something_{{ instance_id }}"
}

type Mount struct {
	Volume string `json:"volume" validate:"nonzero"` // TODO should be VolumeName
	Path   string `json:"path" validate:"nonzero"`
}

type Port struct {
	// TODO Kube only accepts TCP|UDP for protocol values, but we accept values
	// like HTTP, which are used to display component addresses. We should either
	// build a map defining the accepted application protocols on top of TCP|UDP,
	// or make a sep. field.
	Protocol string `json:"protocol" validate:"nonzero" sg:"default=TCP"`

	// Number is the port number used by the container. If your application runs
	// on port 80, for example, use that.
	Number int `json:"number" validate:"nonzero,max=40000"`

	// Public determines whether the port can be accessed ONLY from other
	// Components within Supergiant (false), or from BOTH inside and outside of
	// Supergiant (true). When true, the port can be accessed from external Node
	// IPs. When true, and with an EntrypointDomain provided, the port will be
	// exposed on an external load balancer.
	Public bool `json:"public"`

	// PerInstance, when true, provides each Instance of a Component with its own
	// addressable endpoint (in addition to the normal Component-wide endpoints).
	// When false, Instances can not be reached directly, as traffic to the port
	// is load balanced randomly across all Instances.
	PerInstance bool `json:"per_instance"`

	// EntrypointDomain specifies which Entrypoint this Port is added to. Does not
	// apply when Public is false.
	EntrypointID *int64 `json:"entrypoint_id,omitempty"`

	// ExternalNumber instructs the Entrypoint to set the actual Port number
	// specified as the external load balancer port.
	//
	// TODO validation needed just like on Number, but it can't be nonzero since
	// the value provided can be 0.
	//
	// NOTE Does not apply when EntrypointDomain is nil.
	//      Does not apply to PerInstance ports.
	ExternalNumber int `json:"external_number"`
}
