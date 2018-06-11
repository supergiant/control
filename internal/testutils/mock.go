package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Method names for MockStorage
const (
	StoragePut    = "Put"
	StorageGet    = "Get"
	StorageGetAll = "GetAll"
	StorageDelete = "Delete"
)

// MockStorage is a reusable mock of storage.Interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Put(ctx context.Context, prefix string, key string, value []byte) error {
	args := m.Called(ctx, prefix, key, value)
	return args.Error(0)
}

func (m *MockStorage) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	args := m.Called(ctx, prefix, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorage) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	args := m.Called(ctx, prefix)
	return args.Get(0).([][]byte), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, prefix string, key string) error {
	args := m.Called(ctx, prefix, key)
	return args.Error(0)
}
