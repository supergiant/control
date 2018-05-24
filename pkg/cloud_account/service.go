package cloud_account

import (
	"context"

	"github.com/pkg/errors"
)

type Service struct {
	Repository Repository
}

func (s *Service) GetAll(ctx context.Context) ([]CloudAccount, error) {
	return s.Repository.GetAll(ctx)
}

//TODO add validations
func (s *Service) Get(ctx context.Context, accountName string) (*CloudAccount, error) {
	if accountName == "" {
		return nil, errors.New("cloud account name can't be empty")
	}
	return s.Repository.Get(ctx, accountName)
}

func (s *Service) Create(ctx context.Context, account *CloudAccount) error {
	if account.Name == "" {
		return errors.New("cloud account name can't be empty")
	}
	return s.Repository.Create(ctx, account)
}
