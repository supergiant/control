package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage/file"
)

const (
	testPrefix = "/test/"
)

func TestStorageE2E(t *testing.T) {
	s, err := file.NewFileRepository(fmt.Sprintf("/tmp/sg-storage-%d", time.Now().UnixNano()))
	require.Nil(t, err, "setup file storage provider")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = s.Delete(ctx, testPrefix, "")
	require.NoError(t, err)

	res, err := s.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.Empty(t, res)

	err = s.Put(ctx, testPrefix, "1", []byte("test"))
	require.NoError(t, err)

	getResult, err := s.Get(ctx, testPrefix, "1")
	require.Equal(t, "test", string(getResult))

	err = s.Put(ctx, testPrefix, "2", []byte("test"))
	require.NoError(t, err)

	err = s.Put(ctx, testPrefix, "2", []byte("test222"))
	require.NoError(t, err)

	getResult, err = s.Get(ctx, testPrefix, "2")
	require.Equal(t, "test222", string(getResult))

	res, err = s.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.True(t, len(res) == 2)

	err = s.Delete(ctx, testPrefix, "1")
	require.NoError(t, err)
	err = s.Delete(ctx, testPrefix, "2")
	require.NoError(t, err)

	res, err = s.GetAll(ctx, testPrefix)
	require.NoError(t, err)
	require.Empty(t, res)

	x, err := s.Get(ctx, testPrefix, "NO_SUCH_KEY")
	require.EqualError(t, sgerrors.ErrNotFound, err.Error())
	require.Nil(t, x)

	resultSlice, err := s.GetAll(ctx, "NO_SUCH_PREFIX")
	require.Empty(t, resultSlice)
}
