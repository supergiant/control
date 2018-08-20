package account

import (
	"github.com/supergiant/supergiant/pkg/clouds"
)

// Credentials store ssh keys, fingerprints and other creds associated with CloudAccount
type Credentials map[string]string

// CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
// Name should be unique.
type CloudAccount struct {
	Name        string      `json:"name" valid:"required, length(1|32)"`
	Provider    clouds.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	Credentials Credentials `json:"credentials" valid:"optional"`
}
