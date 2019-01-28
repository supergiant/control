package storage

import (
	"context"

	"github.com/pkg/errors"
)

const (
	memoryStorageType = "memory"
	fileStorageType   = "file"
	etcdStorageType   = "etcd"
)

// Interface is an abstraction over key value storage, gets and returns values serialized as byte slices
// It is up to the services to do data conversion from
type Interface interface {
	GetAll(ctx context.Context, prefix string) ([][]byte, error)
	Get(ctx context.Context, prefix string, key string) ([]byte, error)
	Put(ctx context.Context, prefix string, key string, value []byte) error
	Delete(ctx context.Context, prefix string, key string) error
}

func GetStorage(storageType, uri string) (Interface, error) {
	switch storageType {
	case memoryStorageType:
		return NewInMemoryRepository(), nil
	case fileStorageType:
		return NewFileRepository(uri)
	case etcdStorageType:
		return NewETCDRepository(uri), nil
	}

	return nil, errors.New("wrong storage type" + storageType)
}
