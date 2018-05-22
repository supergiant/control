package cloud_account

import (
	"context"

	"github.com/supergiant/supergiant/pkg/provider"
)

//Credentials store ssh keys, fingerprints and other creds associated with CloudAccount
type Credentials map[string]string

//CloudAccount is settings of account in public or private cloud (e.g. AWS, vCenter)
type CloudAccount struct {
	Name        string
	Provider    provider.Name
	Credentials Credentials
}

//Repository is used to abstract domain storage from implementation
type Repository interface {
	Create(context.Context, *CloudAccount) (error)
	Get(ctx context.Context, accountName string) (*CloudAccount, error)
	GetAll(ctx context.Context) ([]CloudAccount, error)
	Update(context.Context, *CloudAccount) (error)
	Delete(ctx context.Context, accountName string) (error)
}
