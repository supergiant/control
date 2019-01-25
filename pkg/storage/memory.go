package storage

import (
	"context"
	"sync"
	"strings"

	"github.com/supergiant/control/pkg/sgerrors"
)

type InMemoryRepository struct {
	m sync.RWMutex
	data map[string][]byte
}

func NewInMemoryRepository() Interface {
	return &InMemoryRepository{
		data: make(map[string][]byte),
	}
}

func (i *InMemoryRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	i.m.RLock()
	defer i.m.RUnlock()

	value, ok := i.data[prefix + key]

	if !ok {
		return nil, sgerrors.ErrNotFound
	}

	return value, nil
}

func (i *InMemoryRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	i.m.Lock()
	defer i.m.Unlock()

	i.data[prefix + key] = value
	return nil
}

func (i *InMemoryRepository) Delete(ctx context.Context, prefix string, key string) error {
	i.m.Lock()
	defer i.m.Unlock()


	delete(i.data, prefix + key)
	return nil
}

func (i *InMemoryRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	i.m.RLock()
	defer i.m.RUnlock()

	allKeys := make([][]byte, len(i.data))

	for key := range i.data {
		if strings.Contains(key, prefix) {
			allKeys = append(allKeys, i.data[key])
		}
	}

	return allKeys, nil
}