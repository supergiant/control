package profile

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
)

type KubeProfileService struct {
	prefix         string
	profileStorage storage.Interface
}

func NewKubeProfileService(prefix string, s storage.Interface) *KubeProfileService {
	return &KubeProfileService{
		prefix:         prefix,
		profileStorage: s,
	}
}

func (s *KubeProfileService) Get(ctx context.Context, profileId string) (*KubeProfile, error) {
	profileData, err := s.profileStorage.Get(ctx, s.prefix, profileId)
	profile := &KubeProfile{}

	if err != nil {
		return nil, err
	}

	// TODO(stgleb): maybe we don't need to unmarshall it in service to marshall again in handler?
	err = json.Unmarshal(profileData, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *KubeProfileService) Create(ctx context.Context, profile *KubeProfile) error {
	profileData, err := json.Marshal(profile)

	if err != nil {
		return err
	}

	return s.profileStorage.Create(ctx, s.prefix, profile.Id, profileData)
}

func (s *KubeProfileService) GetAll(ctx context.Context) ([]KubeProfile, error) {
	var (
		profiles []KubeProfile
		profile  KubeProfile
	)

	profilesData, err := s.profileStorage.GetAll(ctx, s.prefix)

	if err != nil {
		return nil, err
	}

	for _, profileData := range profilesData {
		err = json.Unmarshal(profileData, &profile)

		if err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
