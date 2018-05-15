package profile

import "github.com/supergiant/supergiant/bindata"

type Interface interface {
	GetUserData(profileName string) ([]byte, error)
}

type Service struct {
}

func (s *Service) GetUserData(profileName string) ([]byte, error) {
	data, err := bindata.Asset("config/providers/" + profileName)
	return data, err
}
