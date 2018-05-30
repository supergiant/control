package repository

import (
	"context"

	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/model/helm"
)

// Service manages helm repositories.
type Service struct {
	prefix string
	s      storage.Interface
}

// NewService constructs a Service for helm repository.
func NewService(s storage.Interface) (*Service, error) {
	return &Service{
		prefix: "/helm/repositories/",
		s: s,
	}, nil
}

// Create stores a helm repository to the provided storage.
func (s *Service) Create(ctx context.Context, r *helm.Repository) error {
	return nil
}

// Get retrieves a helm repository from the storage by its name.
func (s *Service) Get(ctx context.Context, key string) (*helm.Repository, error) {
	return nil, nil
}

// GetAll retrieves all helm repositories from the storage.
func (s *Service) GetAll(ctx context.Context) ([]helm.Repository, error) {
	return nil, nil
}

// Delete removes a helm repository from the storage by its name.
func (s *Service) Delete(ctx context.Context, key string) error {
	return nil
}
