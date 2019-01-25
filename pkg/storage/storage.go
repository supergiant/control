package storage

import (
	"context"

)

// Interface is an abstraction over key value storage, gets and returns values serialized as byte slices
// It is up to the services to do data conversion from
type Interface interface {
	GetAll(ctx context.Context, prefix string) ([][]byte, error)
	Get(ctx context.Context, prefix string, key string) ([]byte, error)
	Put(ctx context.Context, prefix string, key string, value []byte) error
	Delete(ctx context.Context, prefix string, key string) error
}
