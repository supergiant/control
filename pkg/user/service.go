package user

import (
	"context"

	"github.com/pkg/errors"
)

type Service struct {
	Repository Repository
}

func (s *Service) GetByToken(ctx context.Context, apiToken string) (*User, error) {
	users, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.APIToken == apiToken {
			return &u, nil
		}
	}
	return nil, errors.Errorf("user with api token %s not found", apiToken)
}
func (s *Service) GetBySession(ctx context.Context, session string) error {
	//TODO implement
	panic("Implement me")
}
