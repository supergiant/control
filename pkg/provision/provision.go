package provision

import (
	"github.com/supergiant/supergiant/pkg/model"
	"context"
)

type Interface interface {
	CreateMaster(kube *model.Kube, ips []string, ctx context.Context) error
}

type SSHService struct {
}

func (s *SSHService) CreateMaster(kube *model.Kube, ips []string, ctx context.Context) error {
	return nil
}
