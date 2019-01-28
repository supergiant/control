package machines

import (
	"context"
	"encoding/json"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/storage"
)

// Service contains business logic for node in particular cloud provider
type Service struct {
	prefix     string
	repository storage.Interface
}

// NewService constructs a Service.
func NewService(prefix string, s storage.Interface) *Service {
	return &Service{
		prefix:     prefix,
		repository: s,
	}
}

func (s *Service) Create(ctx context.Context, node *model.Machine) error {
	profileData, err := json.Marshal(node)

	if err != nil {
		return err
	}

	return s.repository.Put(ctx, s.prefix, node.ID, profileData)
}

func (s *Service) Get(ctx context.Context, nodeId string) (*model.Machine, error) {
	profileData, err := s.repository.Get(ctx, s.prefix, nodeId)
	profile := &model.Machine{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(profileData, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) ListAll(ctx context.Context) ([]model.Machine, error) {
	var (
		nodes []model.Machine
		node  model.Machine
	)

	data, err := s.repository.GetAll(ctx, s.prefix)

	if err != nil {
		return nil, err
	}

	for _, nodeData := range data {
		err = json.Unmarshal(nodeData, &node)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (s *Service) Delete(ctx context.Context, nodeId string) error {
	return s.repository.Delete(ctx, s.prefix, nodeId)
}
