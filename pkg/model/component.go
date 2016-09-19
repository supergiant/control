package model

type Component struct {
	BaseModel

	// belongs_to App
	App   *App   `json:"app,omitempty"`
	AppID *int64 `json:"app_id" gorm:"not null;index;unique_index:name_within_app"`

	Name string `json:"name" validate:"nonzero,max=17,regexp=^[a-z]([-a-z0-9]*[a-z0-9])?$" gorm:"not null;unique_index:name_within_app"`

	// CustomDeployScript is stored as JSON (like an embed in Mongo)
	CustomDeployScript     *CustomDeployScript `json:"custom_deploy_script,omitempty" gorm:"-" sg:"store_as_json_in=CustomDeployScriptJSON"`
	CustomDeployScriptJSON []byte              `json:"-"`

	// has_many Releases (for preloading)
	Releases []*Release `json:"releases,omitempty"`

	// has_many Instances (for preloading)
	Instances []*Instance `json:"instances,omitempty"`

	// has_one CurrentRelease
	CurrentRelease   *Release `json:"current_release,omitempty" gorm:"ForeignKey:CurrentReleaseID"`
	CurrentReleaseID *int64   `json:"current_release_id" sg:"readonly"`

	// has_one TargetRelease
	TargetRelease   *Release `json:"target_release,omitempty" gorm:"ForeignKey:TargetReleaseID"`
	TargetReleaseID *int64   `json:"target_release_id" sg:"readonly"`

	Addresses     *Addresses `json:"addresses,omitempty" gorm:"-" sg:"store_as_json_in=AddressesJSON"`
	AddressesJSON []byte     `json:"-"`

	// has_many ComponentPrivateImageKeys (really a many2many with PrivateImageKeys)
	PrivateImageKeys []*ComponentPrivateImageKey `json:"private_image_keys"`
}

// Join model
type ComponentPrivateImageKey struct {
	BaseModel
	// belongs_to Component
	Component   *Component `json:"component,omitempty"`
	ComponentID *int64     `json:"component_id" gorm:"not null;index"`
	// belongs_to PrivateImageKey
	Key   *PrivateImageKey `json:"key,omitempty" gorm:"ForeignKey:KeyID"`
	KeyID *int64           `json:"key_id" gorm:"not null;index"`
}

type CustomDeployScript struct {
	Image   string   `json:"image" validate:"nonzero,regexp=^[-\\w\\.\\/]+(:[-\\w\\.]+)?$"`
	Command []string `json:"command"` // TODO need validation here, I think we need to reqire command
	Timeout int      `json:"timeout" sg:"default=1800"`
}

// returns max of the 2 releases
func (m *Component) InstanceCount() int {
	if m.CurrentRelease != nil && m.CurrentRelease.InstanceCount > m.TargetRelease.InstanceCount {
		return m.CurrentRelease.InstanceCount
	}
	return m.TargetRelease.InstanceCount
}

func (m *Component) InstanceByNum(num int) *Instance {
	for _, instance := range m.Instances {
		if instance.Num == num {
			return instance
		}
	}
	return nil
}
