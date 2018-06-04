package storage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	ret := m.Called(ctx, prefix, key)
	r1 := ret.Get(0).([]byte)

	if r2 := ret.Get(1); r2 == nil {
		return r1, nil
	} else {
		return r1, r2.(error)
	}
}

func (m *MockStorage) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	ret := m.Called(ctx, prefix)
	r1 := ret.Get(0).([][]byte)

	if r2 := ret.Get(1); r2 == nil {
		return r1, nil
	} else {
		return r1, r2.(error)
	}
}

func (m *MockStorage) Put(ctx context.Context, prefix string, key string, value []byte) error {
	ret := m.Called(ctx, prefix, key, value)

	if r := ret.Get(0); r == nil {
		return nil
	} else {
		return r.(error)
	}
}

func (m *MockStorage) Delete(ctx context.Context, prefix string, key string) error {
	return nil
}
