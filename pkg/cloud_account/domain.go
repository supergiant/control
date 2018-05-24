package cloud_account

import (
	"github.com/supergiant/supergiant/pkg/provider"
)

//Credentials store ssh keys, fingerprints and other creds associated with CloudAccount
type Credentials map[string]string

//CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
type CloudAccount struct {
	Name        string
	Provider    provider.Name
	Credentials Credentials
}
