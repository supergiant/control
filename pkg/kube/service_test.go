package kube

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/testutils"
)

func TestKubeServiceGet(t *testing.T) {
	testCases := []struct {
		expectedName string
		data         []byte
		err          error
	}{
		{
			expectedName: "kube-name-1234",
			data:         []byte(`{"name":"kube-name-1234"}`),
			err:          nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := DefaultStoragePrefix

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), prefix, "fake_id").Return(testCase.data, testCase.err)

		service := NewService(prefix, m)

		kube, err := service.Get(context.Background(), "fake_id")

		if testCase.err != errors.Cause(err) {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && kube.Name != testCase.expectedName {
			t.Errorf("Wrong kube name expected %s actual %s", testCase.expectedName, kube.Name)
		}
	}
}

func TestKubeServiceCreate(t *testing.T) {
	testCases := []struct {
		kube *model.Kube
		err  error
	}{
		{
			kube: &model.Kube{},
			err:  nil,
		},
		{
			kube: &model.Kube{},
			err:  errors.New("test err"),
		},
	}

	prefix := DefaultStoragePrefix

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		kubeData, _ := json.Marshal(testCase.kube)

		m.On("Put",
			context.Background(),
			prefix,
			testCase.kube.Name,
			kubeData).
			Return(testCase.err)

		service := NewService(prefix, m)

		err := service.Create(context.Background(), testCase.kube)

		if testCase.err != errors.Cause(err) {
			t.Errorf("Unexpected error when create node %v", err)
		}
	}
}

func TestKubeServiceGetAll(t *testing.T) {
	testCases := []struct {
		data [][]byte
		err  error
	}{
		{
			data: [][]byte{[]byte(`{"name":"kube-name-1234"}`), []byte(`{"id":"56kube-name-5678"}`)},
			err:  nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := DefaultStoragePrefix

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("GetAll", context.Background(), prefix).Return(testCase.data, testCase.err)

		service := NewService(prefix, m)

		kubes, err := service.ListAll(context.Background())

		if testCase.err != errors.Cause(err) {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && len(kubes) != 2 {
			t.Errorf("Wrong len of kubes expected 2 actual %d", len(kubes))
		}
	}
}
