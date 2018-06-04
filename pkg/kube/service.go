package kube

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
)

type KubeService struct {
	prefix      string
	kubeStorage storage.Interface
}

func NewKubeService(prefix string, kubeStorage storage.Interface) *KubeService {
	return &KubeService{
		prefix,
		kubeStorage,
	}
}

func (s *KubeService) Get(ctx context.Context, id string) (*Kube, error) {
	kubeData, err := s.kubeStorage.Get(ctx, s.prefix, id)
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

func (s *KubeService) Create(ctx context.Context, kube *Kube) error {
	kubeData, err := json.Marshal(kube)

	if err != nil {
		return err
	}

	return s.kubeStorage.Put(ctx, s.prefix, kube.Name, kubeData)
}

func (s *KubeService) GetAll(ctx context.Context) ([]Kube, error) {
	var (
		kubes []Kube
		kube  Kube
	)

	kubesData, err := s.kubeStorage.GetAll(ctx, s.prefix)

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

func (s *KubeService) Delete(ctx context.Context, id string) error {
	return s.kubeStorage.Delete(ctx, s.prefix, id)
}
