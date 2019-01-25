package storage

import (
	"testing"
	"context"

	"github.com/supergiant/control/pkg/sgerrors"
)

func TestNewInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()

	if repo == nil {
		t.Errorf("repo must not be nil")
	}
}

func TestInMemoryRepository_GetSuccess(t *testing.T) {
	repo := &InMemoryRepository{
		data: map[string][]byte{
			"prefixkey": []byte(`value`),
		},
	}

	if _, err := repo.Get(context.Background(), "prefix", "key"); err != nil {
		t.Errorf("Unexpected error for get key: key")
	}
}

func TestInMemoryRepository_GetErr(t *testing.T) {
	repo := &InMemoryRepository{
		data: map[string][]byte{
			"prefixkey": []byte(`value`),
		},
	}

	if _, err := repo.Get(context.Background(), "prefix", "notfound"); err != sgerrors.ErrNotFound {
		t.Errorf("Unexpected error value %v", err)
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo := &InMemoryRepository{
		data: map[string][]byte{
			"prefixkey": []byte(`value`),
		},
	}

	if err := repo.Delete(context.Background(), "prefix", "key"); err != nil {
		t.Errorf("Unexpected error when delete key")
	}
}

func TestInMemoryRepository_Put(t *testing.T) {
	repo := &InMemoryRepository{
		data: map[string][]byte{
			"prefixkey": []byte(`value`),
		},
	}

	if err := repo.Put(context.Background(), "prefix", "key", []byte(`value`)); err != nil {
		t.Errorf("Unexpected error when put key")
	}

	if _, err := repo.Get(context.Background(), "prefix", "key"); err != nil {
		t.Errorf("Unexpected error for get key: key")
	}
}

func TestInMemoryRepository_GetAll(t *testing.T) {
	repo := &InMemoryRepository{
		data: map[string][]byte{
			"prefixkeyone": []byte(`value1`),
			"prefixkeytwo": []byte(`value2`),
			"prefixkeythree": []byte(`value3`),
		},
	}

	if keys, err := repo.GetAll(context.Background(), "prefix"); err != nil {
		t.Errorf("Unexpected error for get all keys")

		if len(keys) != 3 {
			t.Errorf("Wrong key count expected 3 actual %d", len(keys))
		}
	}
}
