package profile

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
)

type NodeProfileService struct {
	prefix             string
	nodeProfileStorage storage.Interface
}

func (s *NodeProfileService) Get(ctx context.Context, profileId string) (*NodeProfile, error) {
	profileData, err := s.nodeProfileStorage.Get(ctx, s.prefix, profileId)
	profile := &NodeProfile{}

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

func (s *NodeProfileService) Create(ctx context.Context, profile *NodeProfile) error {
	profileData, err := json.Marshal(profile)

	if err != nil {
		return err
	}

	return s.nodeProfileStorage.Put(ctx, s.prefix, profile.Id, profileData)
}

func (s *NodeProfileService) GetAll(ctx context.Context) ([]NodeProfile, error) {
	var (
		profiles []NodeProfile
		profile  NodeProfile
	)

	profilesData, err := s.nodeProfileStorage.GetAll(ctx, s.prefix)

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
