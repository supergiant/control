package etcd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/supergiant/control/pkg/sgerrors"
)

type ETCDRepository struct {
	cfg clientv3.Config
}

func NewETCDRepository(uri string) *ETCDRepository {
	return &ETCDRepository{
		cfg: clientv3.Config{
			Endpoints: []string{uri},
		},
	}
}

func (e *ETCDRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	cl, err := e.GetClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the etcd")
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	res, err := kv.Get(ctx, prefix+key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from the etcd")
	}
	if res.Count == 0 {
		return nil, sgerrors.ErrNotFound
	}
	return res.Kvs[0].Value, nil
}

func (e *ETCDRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	cl, err := e.GetClient()
	if err != nil {
		return errors.Wrap(err, "failed to connect to the etcd")
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	_, err = kv.Put(ctx, prefix+key, string(value))
	return errors.Wrap(err, "failed to write to the etcd")
}

func (e *ETCDRepository) Delete(ctx context.Context, prefix string, key string) error {
	cl, err := e.GetClient()
	if err != nil {
		return errors.Wrap(err, "failed to connect to the etcd")
	}
	_, err = cl.Delete(ctx, prefix+key, clientv3.WithPrefix())
	return errors.Wrap(err, "failed to read from the etcd")
}

func (e *ETCDRepository) GetClient() (*clientv3.Client, error) {
	client, err := clientv3.New(e.cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (e *ETCDRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	result := make([][]byte, 0)

	cl, err := e.GetClient()
	if err != nil {
		return result, errors.Wrap(err, "failed to connect to the etcd")
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	r, err := kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return result, errors.Wrap(err, "failed to read from the etcd")
	}
	for _, v := range r.Kvs {
		result = append(result, v.Value)
	}
	return result, nil
}
