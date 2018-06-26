// +build integration

package profile

import (
	"context"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
)

func TestNodeProfileGet(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)

	testCases := []struct {
		expectedId string
		data       []byte
		err        error
	}{
		{
			expectedId: "1234",
			data:       []byte(`{"id": "1234", "size": "2gb", "image": "ubuntu"}`),
			err:        nil,
		},
		{
			data: nil,
			err:  storage.ErrKeyNotFound,
		},
	}

	prefix := "/node/"

	for _, testCase := range testCases {
		if len(testCase.expectedId) > 0 {
			kv.Put(context.Background(), prefix, testCase.expectedId, testCase.data)
		}

		service := profile.NewNodeProfileService(prefix, kv)
		p, err := service.Get(context.Background(), testCase.expectedId)

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && p.ID != testCase.expectedId {
			t.Errorf("Wrong profile id expected %s actual %s", testCase.expectedId, p.ID)
		}
	}
}

func TestNodeProfileCreate(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)
	prefix := "/node/"
	key := "key"

	node := &profile.NodeProfile{
		ID:    key,
		Size:  "2gb",
		Image: "ubuntu-16.04",
	}

	service := profile.NewNodeProfileService(prefix, kv)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := service.Create(ctx, node)

	if err != nil {
		t.Errorf("Unepexpected error while creating node profile %v", err)
	}

	node2, err := service.Get(ctx, node.ID)

	if err != nil {
		t.Errorf("Unexpected error while getting node profile %v", err)
	}

	if node.ID != key || node.Image != node2.Image || node.Size != node2.Size {
		t.Errorf("Wrong data in etcd")
	}
}

func TestNodeProfileGetAll(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)
	prefix := "/node/"
	key := "key"

	service := profile.NewNodeProfileService(prefix, kv)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	node := &profile.NodeProfile{
		ID:    key,
		Size:  "2gb",
		Image: "ubuntu-16.04",
	}

	err := service.Create(ctx, node)

	if err != nil {
		t.Errorf("Unepexpected error while creating node profile %v", err)
	}

	nodeProfiles, err := service.GetAll(ctx)

	if err != nil {
		t.Errorf("Unexpected error getting node profiles %v", err)
	}

	if len(nodeProfiles) == 0 {
		t.Error("Node profiles are empty")
	}
}
