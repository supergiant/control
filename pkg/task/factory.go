package task

import (
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/provider"
)

func ParseNodeProfile(profile *profile.NodeProfile) {
	switch profile.Provider {
	case provider.AWS:
		//spawn aws tasks
	}
}
