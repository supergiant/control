package profile

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
)

const DefaultKubeProfilePreifx = "/supergiant/profile"

type KubeProfileService struct {
	prefix             string
	kubeProfileStorage storage.Interface
}

func NewKubeProfileService(prefix string, s storage.Interface) *KubeProfileService {
	return &KubeProfileService{
		prefix:             prefix,
		kubeProfileStorage: s,
	}
}

func (s *KubeProfileService) Get(ctx context.Context, profileId string) (*Profile, error) {
	profileData, err := s.kubeProfileStorage.Get(ctx, s.prefix, profileId)
	profile := &Profile{}

	if err != nil {
		return nil, err
	}

	// TODO(stgleb): maybe we don't need to unmarshal it in service to marshall again in handler?
	err = json.Unmarshal(profileData, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *KubeProfileService) Create(ctx context.Context, profile *Profile) error {
	profileData, err := json.Marshal(profile)

	if err != nil {
		return err
	}

	return s.kubeProfileStorage.Put(ctx, s.prefix, profile.ID, profileData)
}

func (s *KubeProfileService) GetAll(ctx context.Context) ([]Profile, error) {
	var (
		profiles []Profile
		profile  Profile
	)

	profilesData, err := s.kubeProfileStorage.GetAll(ctx, s.prefix)

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
