package storage

import (
	"context"
)

type Fake struct {
	Item      []byte
	Items     [][]byte
	PutErr    error
	GetErr    error
	ListErr   error
	DeleteErr error
}

func (s Fake) Put(ctx context.Context, prefix string, key string, value []byte) error {
	return s.PutErr
}

func (s Fake) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	return s.Item, s.GetErr
}

func (s Fake) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	return s.Items, s.ListErr
}

func (s Fake) Delete(ctx context.Context, prefix string, key string) error {
	return s.DeleteErr
}
