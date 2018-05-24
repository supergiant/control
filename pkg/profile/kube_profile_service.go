package profile

import (
	"github.com/supergiant/supergiant/pkg/storage"
	"context"
)

type KubeProfileService struct{
	prefix string
	profileStorage storage.Interface
}

func NewKubeProfileService(prefix string, s storage.Interface) KubeProfileService {
	return KubeProfileService{
		prefix:prefix,
		profileStorage: s,
	}
}

func (s *KubeProfileService) Get(ctx context.Context, profileId string) (*KubeProfile, error) {
	profiles, err := s.profileStorage.GetAll(ctx, s.prefix)

	if err != nil {
		return nil, err
	}


	return nil, nil
}

func (s *KubeProfileService) Create(ctx context.Context, profile *KubeProfile) error {
	return nil
}

func (s *KubeProfileService) GetAll(ctx context.Context) ([]KubeProfile, error) {
	return nil, nil
}
