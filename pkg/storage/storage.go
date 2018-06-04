package storage

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

// Interface is an abstraction over key value storage, gets and returns values serialized as byte slices
// It is up to the services to do data conversion from
type Interface interface {
	GetAll(ctx context.Context, prefix string) ([][]byte, error)
	Get(ctx context.Context, prefix string, key string) ([]byte, error)
	Put(ctx context.Context, prefix string, key string, value []byte) error
	Delete(ctx context.Context, prefix string, key string) error
}

type ETCDRepository struct {
	cfg clientv3.Config
}

var ErrKeyNotFound = errors.New("key not found")

func (e *ETCDRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	cl, err := e.GetClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	res, err := kv.Get(ctx, prefix+key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res.Count == 0 {
		return nil, ErrKeyNotFound
	}
	return res.Kvs[0].Value, nil
}

func (e *ETCDRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	cl, err := e.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	_, err = kv.Put(ctx, prefix+key, string(value))
	return errors.WithStack(err)
}

func (e *ETCDRepository) Delete(ctx context.Context, prefix string, key string) error {
	cl, err := e.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = cl.Delete(ctx, prefix, clientv3.WithPrefix())
	return errors.WithStack(err)
}

func (e *ETCDRepository) GetClient() (*clientv3.Client, error) {
	client, err := clientv3.New(e.cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}

func (e *ETCDRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	result := make([][]byte, 0)

	cl, err := e.GetClient()
	if err != nil {
		return result, errors.WithStack(err)
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	r, err := kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return result, errors.WithStack(err)
	}
	for _, v := range r.Kvs {
		result = append(result, v.Value)
	}
	return result, nil
}

func NewETCDRepository(cfg clientv3.Config) Interface {
	return &ETCDRepository{
		cfg: cfg,
	}
}
