package node

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/storage"
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

func (s *Service) Create(ctx context.Context, node *Node) error {
	profileData, err := json.Marshal(node)

	if err != nil {
		return err
	}

	return s.repository.Put(ctx, s.prefix, node.Id, profileData)
}

func (s *Service) Get(ctx context.Context, nodeId string) (*Node, error) {
	profileData, err := s.repository.Get(ctx, s.prefix, nodeId)
	profile := &Node{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(profileData, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) ListAll(ctx context.Context) ([]Node, error) {
	var (
		nodes []Node
		node  Node
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
