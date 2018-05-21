package provision

import (
	"context"

	"github.com/supergiant/supergiant/pkg/profile"
)

//Settings of the provision
type Settings struct {
	IPS        []string
	K8SVersion string
	//TODO Add credentials handling
}

//Interface should be used to provision a node to kubernetes cluster with given settings
type Interface interface {
	Provision(context.Context, *profile.NodeProfile, *Settings) error
}
