package pki

import (
	"context"
	"net"

	"github.com/satori/go.uuid"

	"github.com/supergiant/supergiant/pkg/storage"
)

type Service struct {
	storagePrefix string
	repository    storage.Interface
}

func NewService(storagePrefix string, repository storage.Interface) *Service {
	return &Service{
		storagePrefix: storagePrefix,
		repository:    repository,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*PKI, error) {
	rawData, err := s.repository.GetAll(ctx, s.storagePrefix)
	if err != nil {
		return nil, err
	}

	pkis := make([]*PKI, len(rawData))
	for _, v := range rawData {
		pki, err := Unmarshall(v)
		if err != nil {
			return nil, err
		}

		pkis = append(pkis, pki)
	}

	return pkis, nil
}

func (s *Service) Get(ctx context.Context, ID string) (*PKI, error) {
	rawData, err := s.repository.Get(ctx, s.storagePrefix, ID)
	if err != nil {
		return nil, err
	}
	return Unmarshall(rawData)
}

func (s *Service) Delete(ctx context.Context, ID string) error {
	return s.repository.Delete(ctx, s.storagePrefix, ID)
}

func (s *Service) GenerateFromCA(ctx context.Context, CA *PairPEM, dnsDomain string, masterIPs []net.IP) (*PKI, error) {
	p, err := NewPKI(CA, dnsDomain, masterIPs)
	if err != nil {
		return nil, err
	}
	ID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	p.ID = ID.String()
	if err := s.repository.Put(ctx, s.storagePrefix, p.ID, p.Marshall()); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GenerateSelfSigned(ctx context.Context, dnsDomain string, masterIPs []net.IP) (*PKI, error) {
	p, err := NewPKI(nil, dnsDomain, masterIPs)
	if err != nil {
		return nil, err
	}

	ID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	p.ID = ID.String()

	if err := s.repository.Put(ctx, s.storagePrefix, p.ID, p.Marshall()); err != nil {
		return nil, err
	}
	return p, nil
}
