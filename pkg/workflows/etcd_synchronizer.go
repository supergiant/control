package workflows

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

type EtcdSynchronizer struct {
	client *clientv3.Client
}

func NewSynchronizer(cfg clientv3.Config) (*EtcdSynchronizer, error) {
	client, err := clientv3.New(cfg)

	if err != nil {
		return nil, err
	}

	return &EtcdSynchronizer{
		client: client,
	}, nil
}

// TODO(stgleb): Pass context and key to synchronizer
func (s *EtcdSynchronizer) Sync(ctx context.Context, key, data string) error {
	_, err := s.client.Put(ctx, key, data)

	if err != nil {
		return err
	}

	return nil
}
