package model

import (
	"github.com/supergiant/supergiant/pkg/clouds"
)

const (
	KeyFingerprint = "fingerprint"
	KeyPrivateKey  = "privatekey"
)

// CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
// Name should be unique.
type CloudAccount struct {
	Name        string            `json:"name" valid:"required, length(1|32)"`
	Provider    clouds.Name       `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	Credentials map[string]string `json:"credentials" valid:"optional"`
}
