package storage

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"encoding/json"
	"reflect"
)

// Interface is an abstraction over key value storage, gets and returns values serialized as byte slices
// It is up to the services to do data conversion from
type Interface interface {
	GetAll(ctx context.Context, prefix string, resultType reflect.Type) ([]interface{}, error)
	Get(ctx context.Context, prefix string, key string, resultType reflect.Type) (interface{}, error)
	Put(ctx context.Context, prefix string, key string, value interface{}) error
	Delete(ctx context.Context, prefix string, key string) error
}

type ETCDRepository struct {
	cfg clientv3.Config
}

func (e *ETCDRepository) Get(ctx context.Context, prefix string, key string, resultType reflect.Type) (interface{}, error) {
	cl, err := e.GetClient()
	if err != nil {
		return nil, err
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	res, err := kv.Get(ctx, prefix+key)
	if err != nil {
		return nil, err
	}
	if res.Count == 0 {
		return nil, nil
	}
	obj := reflect.New(resultType).Interface()
	err = json.Unmarshal(res.Kvs[0].Value, obj)
	if err != nil {
		return nil, err
	}
	return obj, err
}

func (e *ETCDRepository) Put(ctx context.Context, prefix string, key string, value interface{}) error {
	if value == nil {
		return errors.Errorf("can't save nil to key %s", key)
	}
	cl, err := e.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}
	defer cl.Close()
	kv := clientv3.NewKV(cl)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = kv.Put(ctx, prefix+key, string(data))
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

func (e *ETCDRepository) GetAll(ctx context.Context, prefix string, resultType reflect.Type) ([]interface{}, error) {
	result := make([]interface{}, 0)

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
		obj := reflect.New(resultType)
		err := json.Unmarshal(v.Value, obj)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

func NewETCDRepository(cfg clientv3.Config) *ETCDRepository {
	return &ETCDRepository{
		cfg: cfg,
	}
}
