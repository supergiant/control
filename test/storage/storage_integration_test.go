// +build integration

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/assert"
	"github.com/supergiant/supergiant/pkg/storage"
)

const (
	defaultETCDHost = "http://127.0.0.1:2379"
	testPrefix      = "/test/"
)

var defaultConfig clientv3.Config

func init() {
	assert.EtcdRunning(defaultETCDHost)
	defaultConfig = clientv3.Config{
		Endpoints: []string{defaultETCDHost},
	}
}

func TestStorageE2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	kv := storage.NewETCDRepository(defaultConfig)

	err := kv.Delete(ctx, testPrefix, "")
	require.NoError(t, err)

	res, err := kv.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.Empty(t, res)

	err = kv.Put(ctx, testPrefix, "1", []byte("test"))
	require.NoError(t, err)

	getResult, err := kv.Get(ctx, testPrefix, "1")
	require.Equal(t, "test", string(getResult))

	err = kv.Put(ctx, testPrefix, "2", []byte("test"))
	require.NoError(t, err)

	err = kv.Put(ctx, testPrefix, "2", []byte("test222"))
	require.NoError(t, err)

	getResult, err = kv.Get(ctx, testPrefix, "2")
	require.Equal(t, "test222", string(getResult))

	res, err = kv.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.True(t, len(res) == 2)

	err = kv.Delete(ctx, testPrefix, "")
	require.NoError(t, err)

	res, err = kv.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.Empty(t, res)

	x, err := kv.Get(ctx, testPrefix, "NO_SUCH_KEY")
	require.NoError(t, err)
	require.Nil(t, x)

	resultSlice, err := kv.GetAll(ctx, "NO_SUCH_PREFIX")
	require.Empty(t, resultSlice)
}
