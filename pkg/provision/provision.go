package provision

import (
	"context"
	"github.com/supergiant/supergiant/pkg/profile"
)

type Settings struct {
	IPS        []string
	K8SVersion string
	//TODO Add credentials handling
}

type Interface interface {
	Provision(context.Context, *profile.NodeProfile, *Settings) (error)
}
