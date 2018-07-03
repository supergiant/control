// +build integration

package profile

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
)

const (
	defaultETCDHost = "http://127.0.0.1:2379"
)

var defaultConfig clientv3.Config

func init() {
	assert.EtcdRunning(defaultETCDHost)
	defaultConfig = clientv3.Config{
		Endpoints: []string{defaultETCDHost},
	}
}

func TestKubeProfileGet(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)

	testCases := []struct {
		expectedId string
		data       []byte
		err        error
	}{
		{
			expectedId: "1234",
			data:       []byte(`{"id":"1234", "nodes":[{},{}]}`),
			err:        nil,
		},
		{
			data: nil,
			err:  storage.ErrKeyNotFound,
		},
	}

	prefix := "/kube/"

	for _, testCase := range testCases {
		if len(testCase.expectedId) > 0 {
			kv.Put(context.Background(), prefix, testCase.expectedId, testCase.data)
		}

		service := profile.NewKubeProfileService(prefix, kv)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		p, err := service.Get(ctx, testCase.expectedId)

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && p.Id != testCase.expectedId {
			t.Errorf("Wrong profile id expected %s actual %s", testCase.expectedId, p.Id)
		}
	}
}

func TestKubeProfileCreate(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)
	prefix := "/kube/"
	key := "key"
	version := "1.8.7"

	kube := &profile.KubeProfile{
		Id:                key,
		KubernetesVersion: version,
	}

	service := profile.NewKubeProfileService(prefix, kv)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := service.Create(ctx, kube)

	if err != nil {
		t.Errorf("Unepexpected error while creating kube profile %v", err)
	}

	kube2, err := service.Get(ctx, kube.Id)

	if err != nil {
		t.Errorf("Unexpected error while getting kube profile %v", err)
	}

	if kube.Id != key || kube.KubernetesVersion != kube2.KubernetesVersion {
		t.Errorf("Wrong data in etcd")
	}
}

func TestKubeProfileGetAll(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)
	prefix := "/kube/"
	key := "key"
	version := "1.8.7"

	service := profile.NewKubeProfileService(prefix, kv)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	kube := &profile.KubeProfile{
		Id:                key,
		KubernetesVersion: version,
	}

	err := service.Create(ctx, kube)

	if err != nil {
		t.Errorf("Unepexpected error while creating kube profile %v", err)
	}

	kubeProfiles, err := service.GetAll(ctx)

	if err != nil {
		t.Errorf("Unexpected error getting kube profiles %v", err)
	}

	if len(kubeProfiles) == 0 {
		t.Error("Kube profiles are empty")
	}
}
