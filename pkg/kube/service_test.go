package kube

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
)

type mockServerResourceGetter struct {
	resources []*metav1.APIResourceList
	err       error
}

func (m *mockServerResourceGetter) ServerResources() ([]*metav1.APIResourceList, error) {
	return m.resources, m.err
}

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
		m.On("Get", context.Background(), prefix, "fake_id").
			Return(testCase.data, testCase.err)

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

func TestService_Delete(t *testing.T) {
	testCases := []struct {
		repoErr error
	}{
		{
			sgerrors.ErrNotFound,
		},
		{
			nil,
		},
	}

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Delete", context.Background(), mock.Anything, mock.Anything).
			Return(testCase.repoErr)

		service := NewService("", m)

		err := service.Delete(context.Background(), "key")

		if err != testCase.repoErr {
			t.Errorf("expected error %v actual %v", testCase.repoErr, err)
		}
	}
}

func TestResourcesGroupInfo(t *testing.T) {
	testCases := []struct {
		discoveryErr       error
		resourceErr        error
		resourcesLists     []*metav1.APIResourceList
		expectedGroupCount int
		expectedErr        error
	}{
		{
			discoveryErr: sgerrors.ErrNotFound,
			expectedErr:  sgerrors.ErrNotFound,
		},
		{
			resourceErr: sgerrors.ErrNotFound,
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			resourcesLists: []*metav1.APIResourceList{
				{
					GroupVersion: "",
					APIResources: []metav1.APIResource{
						{
							Name: "name-1",
							Kind: "kind1",
						},
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
				{
					GroupVersion: "/",
					APIResources: []metav1.APIResource{
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
			},
			expectedGroupCount: 2,
		},
	}

	for _, testCase := range testCases {
		m := &mockServerResourceGetter{
			resources: testCase.resourcesLists,
			err:       testCase.resourceErr,
		}

		svc := Service{
			discoveryClientFn: func(k *model.Kube) (ServerResourceGetter, error) {
				return m, testCase.discoveryErr
			},
		}

		groups, err := svc.resourcesGroupInfo(&model.Kube{})

		if errors.Cause(err) != testCase.expectedErr {
			t.Errorf("expected error %v actual %v",
				testCase.expectedErr, err)
		}

		if len(groups) != testCase.expectedGroupCount {
			t.Errorf("expected group count %d actual %d",
				testCase.expectedGroupCount, len(groups))
		}
	}
}

func TestListKubeResources(t *testing.T) {
	testCases := []struct {
		kubeData           []byte
		getkubeErr         error
		discoveryErr       error
		resourceErr        error
		resourcesLists     []*metav1.APIResourceList
		expectedGroupCount int
		expectedErr        error
	}{
		{
			getkubeErr:  sgerrors.ErrNotFound,
			kubeData:    nil,
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			getkubeErr: nil,
			kubeData:   []byte(`{"name":"kube-name-1234"}`),
			resourcesLists: []*metav1.APIResourceList{
				{
					GroupVersion: "",
					APIResources: []metav1.APIResource{
						{
							Name: "name-1",
							Kind: "kind1",
						},
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
				{
					GroupVersion: "/",
					APIResources: []metav1.APIResource{
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
			},
		},
		{
			getkubeErr:   nil,
			discoveryErr: sgerrors.ErrNotFound,
			kubeData:     []byte(`{"name":"kube-name-1234"}`),
			expectedErr:  sgerrors.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), mock.Anything, mock.Anything).
			Return(testCase.kubeData, testCase.getkubeErr)

		mockResourceGetter := &mockServerResourceGetter{
			resources: testCase.resourcesLists,
			err:       testCase.resourceErr,
		}

		svc := Service{
			storage: m,
			discoveryClientFn: func(k *model.Kube) (ServerResourceGetter, error) {
				return mockResourceGetter, testCase.discoveryErr
			},
		}

		_, err := svc.ListKubeResources(context.Background(), "kube-name-1234")

		if errors.Cause(err) != testCase.expectedErr {
			t.Errorf("expected error %v actual %v",
				testCase.expectedErr, err)
		}
	}
}

func TestService_GetKubeResources(t *testing.T) {
	testCases := []struct {
		kubeData           []byte
		getkubeErr         error
		discoveryErr       error
		resourceErr        error
		resourceName       string
		resourcesLists     []*metav1.APIResourceList
		expectedGroupCount int
		expectedErr        error
	}{
		{
			getkubeErr:  sgerrors.ErrNotFound,
			kubeData:    nil,
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			getkubeErr: nil,
			kubeData:   []byte(`{"name":"kube-name-1234"}`),
			resourcesLists: []*metav1.APIResourceList{
				{
					GroupVersion: "",
					APIResources: []metav1.APIResource{
						{
							Name: "name-1",
							Kind: "kind1",
						},
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
				{
					GroupVersion: "/",
					APIResources: []metav1.APIResource{
						{
							Name: "name-2",
							Kind: "kind2",
						},
					},
				},
			},
			resourceName: "name-3",
			expectedErr:  sgerrors.ErrNotFound,
		},
		{
			getkubeErr:   nil,
			discoveryErr: sgerrors.ErrNotFound,
			kubeData:     []byte(`{"name":"kube-name-1234"}`),
			expectedErr:  sgerrors.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), mock.Anything, mock.Anything).
			Return(testCase.kubeData, testCase.getkubeErr)

		mockResourceGetter := &mockServerResourceGetter{
			resources: testCase.resourcesLists,
			err:       testCase.resourceErr,
		}

		svc := Service{
			storage: m,
			discoveryClientFn: func(k *model.Kube) (ServerResourceGetter, error) {
				return mockResourceGetter, testCase.discoveryErr
			},
		}

		_, err := svc.GetKubeResources(context.Background(),
			"kube-name-1234", testCase.resourceName,
			"namaspace", testCase.resourceName)

		if errors.Cause(err) != testCase.expectedErr {
			t.Errorf("expected error %v actual %v",
				testCase.expectedErr, err)
		}
	}
}

func TestService_GetCerts(t *testing.T) {
	testCases := []struct {
		kname       string
		cname       string
		data        []byte
		getErr      error
		sshErr      error
		expectedErr error
	}{
		{
			data:        nil,
			getErr:      sgerrors.ErrNotFound,
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			kname:       "kube-name-1234",
			data:        []byte(`{"name":"kube-name-1234", "sshUser": "root", "sshKey": ""}`),
			sshErr:      ssh.ErrHostNotSpecified,
			expectedErr: ssh.ErrHostNotSpecified,
		},
	}

	prefix := DefaultStoragePrefix

	for _, testCase := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), prefix, mock.Anything).
			Return(testCase.data, testCase.getErr)

		service := NewService(prefix, m)

		_, err := service.GetCerts(context.Background(),
			testCase.kname, testCase.cname)

		if testCase.expectedErr != errors.Cause(err) {
			t.Errorf("Wrong error expected %v actual %v", testCase.expectedErr, err)
			return
		}
	}
}
