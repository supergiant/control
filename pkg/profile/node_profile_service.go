package profile

import (
	"github.com/supergiant/supergiant/pkg/storage"
)

type NodeProfileService struct {
	storage.Interface
}

func (s *NodeProfileService) Get() (*NodeProfile, error) {
	return nil, nil
}

func (s *NodeProfileService) Create() error {
	return nil
}

func (s *NodeProfileService) GetAll() ([]NodeProfile, error) {
	return nil, nil
}
