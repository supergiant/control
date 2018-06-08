package kube

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
)

type Service struct {
	prefix  string
	storage storage.Interface
}

func NewKubeService(prefix string, kubeStorage storage.Interface) *Service {
	return &Service{
		prefix,
		kubeStorage,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*Kube, error) {
	kubeData, err := s.storage.Get(ctx, s.prefix, id)
	kube := &Kube{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(kubeData, kube)

	if err != nil {
		return nil, err
	}

	return kube, nil
}

func (s *Service) Create(ctx context.Context, kube *Kube) error {
	kubeData, err := json.Marshal(kube)

	if err != nil {
		return err
	}

	return s.storage.Put(ctx, s.prefix, kube.Name, kubeData)
}

func (s *Service) GetAll(ctx context.Context) ([]Kube, error) {
	var (
		kubes []Kube
		kube  Kube
	)

	kubesData, err := s.storage.GetAll(ctx, s.prefix)

	if err != nil {
		return nil, err
	}

	for _, data := range kubesData {
		err = json.Unmarshal(data, &kube)

		if err != nil {
			return nil, err
		}

		kubes = append(kubes, kube)
	}

	return kubes, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.storage.Delete(ctx, s.prefix, id)
}
