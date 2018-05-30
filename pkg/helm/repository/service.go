package repository

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/storage"
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
		s:      s,
	}, nil
}

// Create stores a helm repository in the provided storage.
func (s *Service) Create(ctx context.Context, r *helm.Repository) error {
	rawJSON, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	err = s.s.Put(ctx, s.prefix, r.Name, rawJSON)
	if err != nil {
		return errors.Wrap(err, "storage")
	}

	return nil
}

// Get retrieves a helm repository from the storage by its name.
func (s *Service) Get(ctx context.Context, repoName string) (*helm.Repository, error) {
	res, err := s.s.Get(ctx, s.prefix, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	repo := new(helm.Repository)
	err = json.Unmarshal(res, repo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return nil, nil
}

// GetAll retrieves all helm repositories from the storage.
func (s *Service) GetAll(ctx context.Context) ([]helm.Repository, error) {
	rawRepos, err := s.s.GetAll(ctx, s.prefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	repos := make([]helm.Repository, len(rawRepos))
	for i, raw := range rawRepos {
		repo := new(helm.Repository)
		err = json.Unmarshal(raw, repo)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}
		repos[i] = *repo
	}

	return repos, nil
}

// Delete removes a helm repository from the storage by its name.
func (s *Service) Delete(ctx context.Context, repoName string) error {
	return s.s.Delete(ctx, s.prefix, repoName)
}
