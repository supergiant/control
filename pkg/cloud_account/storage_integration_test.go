//+build integration

package cloud_account

import (
	"reflect"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/supergiant/supergiant/pkg/assert"
	"github.com/supergiant/supergiant/pkg/provider"
)

func init() {
	//TODO tests config should be introduced
	assert.EtcdRunning("http://127.0.0.1:2379")
}

func TestStorageE2E(t *testing.T) {
	cl, err := client.New(client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	})
	require.NoError(t, err)

	repo := NewRepository(cl)
	acc := newAccount()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	all, err := repo.GetAll(ctx)
	require.Empty(t, all)
	require.NoError(t, err)

	defer func() {
		err = repo.Delete(ctx, acc.Name)
		require.NoError(t, err)
		cancel()
	}()

	err = repo.Create(ctx, acc)

	require.NoError(t, err)

	acc2, err := repo.Get(ctx, acc.Name)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(acc, acc2))

	acc.Credentials["test"] = "test"
	err = repo.Update(ctx, acc)
	require.NoError(t, err)

	acc2, err = repo.Get(ctx, acc.Name)
	require.Equal(t, "test", acc2.Credentials["test"])

	//should return nil if account not found
	acc3, err := repo.Get(ctx, "DUMMY")
	require.NoError(t, err)
	require.Nil(t, acc3)

	//should not return error if trying to delete non existent account
	err = repo.Delete(ctx, "DUMMY")
	require.NoError(t, err)

	all, err = repo.GetAll(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all)
}

func newAccount() *CloudAccount {
	return &CloudAccount{
		Name:        "Test",
		Credentials: make(map[string]string),
		Provider:    provider.DigitalOcean,
	}
}
