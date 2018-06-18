package task

import (
	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/provider"
)

//ParseNodeProfile return chain of tasks to provision node as specified to node profile
func ParseNodeProfile(profile *profile.NodeProfile) (*tasks.Chain, error) {
	switch profile.Provider {
	case provider.AWS:
		//spawn aws tasks
	}
	return nil, nil
}
