package account

import (
	"github.com/supergiant/supergiant/pkg/provider"
)

//Credentials store ssh keys, fingerprints and other creds associated with CloudAccount
type Credentials map[string]string

//CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
type CloudAccount struct {
	Name        string        `json:"name" valid:"required, length(1|32)"`
	Provider    provider.Name `json:"provider" valid:"in(aws|digitalocean|packet|gce|openstack)"`
	Credentials Credentials   `json:"credentials" valid:"optional"`
}
