package storage

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/supergiant/supergiant/pkg/assert"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/user"
	"github.com/stretchr/testify/require"
	"reflect"
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
	assert.EtcdRunning(defaultETCDHost)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	kv := storage.NewETCDRepository(defaultConfig)
	tt := []struct {
		key   string
		value *user.User
	}{
		{
			key: "valid_object",
			value: &user.User{
				Password: "test",
				Login:    "test",
			},
		},
	}

	for _, tc := range tt {
		err := kv.Delete(ctx, testPrefix, "")
		require.NoError(t, err)

		err = kv.Put(ctx, testPrefix, tc.key, tc.value)
		require.NoError(t, err)

		obj, err := kv.Get(ctx, testPrefix, tc.key, reflect.TypeOf(tc.value).Elem())
		require.NoError(t, err)

		u, ok := obj.(*user.User)
		require.True(t, ok)

		require.Equal(t, tc.value.Login, u.Login, )
		require.Equal(t, tc.value.Password, u.Password)
	}
}
