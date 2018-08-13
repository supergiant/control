package provisioner

import "github.com/supergiant/supergiant/pkg/profile"

// Provisioner gets kube profile and returns list of task ids of provision tasks
type Provisioner interface {
	Provision(profile profile.KubeProfile) []string
}
