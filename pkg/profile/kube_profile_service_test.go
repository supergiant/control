package profile

import (
	"context"
	"testing"

	"encoding/json"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/testutils"
)

func TestKubeProfileServiceGet(t *testing.T) {
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
			err:  errors.New("test err"),
		},
	}

	prefix := "/kube/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), prefix, "fake_id").Return(testCase.data, testCase.err)

		service := KubeProfileService{
			prefix,
			m,
		}

		profile, err := service.Get(context.Background(), "fake_id")

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && profile.Id != testCase.expectedId {
			t.Errorf("Wrong profile id expected %s actual %s", testCase.expectedId, profile.Id)
		}
	}
}

func TestKubeProfileServiceCreate(t *testing.T) {
	testCases := []struct {
		kube *KubeProfile
		err  error
	}{
		{
			kube: &KubeProfile{},
			err:  nil,
		},
		{
			kube: &KubeProfile{},
			err:  errors.New("test err"),
		},
	}

	prefix := "/kube/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		kubeData, _ := json.Marshal(testCase.kube)

		m.On("Put",
			context.Background(),
			prefix,
			testCase.kube.Id,
			kubeData).
			Return(testCase.err)

		service := KubeProfileService{
			prefix,
			m,
		}

		err := service.Create(context.Background(), testCase.kube)

		if testCase.err != err {
			t.Errorf("Unexpected error when create node %v", err)
		}
	}
}

func TestKubeProfileServiceGetAll(t *testing.T) {
	testCases := []struct {
		data [][]byte
		err  error
	}{
		{
			data: [][]byte{[]byte(`{"id":"1234", "nodes":[{},{}]}`), []byte(`{"id":"5678", "nodes":[{},{}]}`)},
			err:  nil,
		},
		{
			data: nil,
			err:  errors.New("test err"),
		},
	}

	prefix := "/kube/"

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("GetAll", context.Background(), prefix).Return(testCase.data, testCase.err)

		service := KubeProfileService{
			prefix,
			m,
		}

		profiles, err := service.GetAll(context.Background())

		if testCase.err != err {
			t.Errorf("Wrong error expected %v actual %v", testCase.err, err)
			return
		}

		if testCase.err == nil && len(profiles) != 2 {
			t.Errorf("Wrong len of profiles expected 2 actual %d", len(profiles))
		}
	}
}
