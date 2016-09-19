package model

import (
	"encoding/base64"
	"fmt"
)

type PrivateImageKey struct {
	BaseModel

	// SSL  bool   `json:"ssl"`
	Host string `json:"host" validate:"nonzero,regexp=^[^/]+$"`

	Username string `json:"username" validate:"nonzero" gorm:"not null;index"` // Not unique, because maybe multiple registries

	// NOTE these are not stored in database
	Email    string `json:"email" validate:"nonzero" sg:"private" gorm:"-"`
	Password string `json:"password" validate:"nonzero" sg:"private" gorm:"-"`

	// Used for K8S Secret
	Key string `json:"key" validate:"nonzero" sg:"private,readonly"`
}

var keyConfigTemplate = `{
	"auths": {
		"%s": {
			"auth": "%s",
			"email": "%s"
		}
	}
}`

func (m *PrivateImageKey) RegistryURL() string {
	protocol := "https://"
	return protocol + m.Host + "/v1/"
}

func (m *PrivateImageKey) MakeKey() {
	auth := base64.StdEncoding.EncodeToString([]byte(m.Username + ":" + m.Password))
	config := fmt.Sprintf(keyConfigTemplate, m.RegistryURL(), auth, m.Email)
	m.Key = base64.StdEncoding.EncodeToString([]byte(config))
}
