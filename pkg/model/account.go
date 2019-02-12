package model

import (
	"github.com/supergiant/control/pkg/clouds"
)

// CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
// Name should be unique.
type CloudAccount struct {
	Name        string            `json:"name" valid:"required, length(1|32)"`
	Provider    clouds.Name       `json:"provider" valid:"in(aws|digitalocean|gce|azure)"`
	Credentials map[string]string `json:"credentials" valid:"optional"`
}
