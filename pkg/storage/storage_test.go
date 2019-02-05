package storage

import (
	"reflect"
	"testing"

	"github.com/supergiant/control/pkg/storage/etcd"
	"github.com/supergiant/control/pkg/storage/file"
	"github.com/supergiant/control/pkg/storage/memory"
)

func TestGetStorage(t *testing.T) {
	testCases := []struct {
		uri         string
		storageType string
		t           reflect.Type
	}{
		{
			"",
			memoryStorageType,
			reflect.TypeOf(&memory.InMemoryRepository{}),
		},
		{
			"file.db",
			fileStorageType,
			reflect.TypeOf(&file.FileRepository{}),
		},
		{
			"file.db",
			etcdStorageType,
			reflect.TypeOf(&etcd.ETCDRepository{}),
		},
	}

	for _, testCase := range testCases {
		storage, err := GetStorage(testCase.storageType, testCase.uri)

		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if testCase.t != reflect.TypeOf(storage) {
			t.Errorf("Wrong type expected %v actual %v",
				testCase.t, reflect.TypeOf(storage))
		}

	}
}
