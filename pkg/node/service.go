package node

import (
	"context"

	"github.com/supergiant/supergiant/pkg/storage"
)

// Service contains business logic for node in particular cloud provider
type Service struct {
	repository storage.Interface
}

func (s *Service) Create(ctx context.Context, node *Node) error {
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (*Node, error) {
	return nil, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*Node, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return nil
}
