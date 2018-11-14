package kube

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/sghelm/proxy"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/testutils/storage"
)

var (
	fakeRls = &release.Release{
		Name: "fakeRelease",
		Info: &release.Info{
			FirstDeployed: &timestamp.Timestamp{},
			LastDeployed:  &timestamp.Timestamp{},
			Status: &release.Status{
				Code: release.Status_UNKNOWN,
			},
		},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{},
		},
	}
)

type fakeChartGetter struct {
	chrt *chart.Chart
	err  error
}

func (f fakeChartGetter) GetChart(ctx context.Context, repoName, chartName, chartVersion string) (*chart.Chart, error) {
	return f.chrt, f.err
}

type fakeHelmProxy struct {
	proxy.Interface

	err               error
	installRlsResp    *services.InstallReleaseResponse
	listReleaseResp   *services.ListReleasesResponse
	uninstReleaseResp *services.UninstallReleaseResponse
}

func (p *fakeHelmProxy) InstallReleaseFromChart(chart *chart.Chart, namespace string, opts ...helm.InstallOption) (*services.InstallReleaseResponse, error) {
	return p.installRlsResp, p.err
}
func (p *fakeHelmProxy) ListReleases(opts ...helm.ReleaseListOption) (*services.ListReleasesResponse, error) {
	return p.listReleaseResp, p.err
}
func (p *fakeHelmProxy) DeleteRelease(rlsName string, opts ...helm.DeleteOption) (*services.UninstallReleaseResponse, error) {
	return p.uninstReleaseResp, p.err
}

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

		service := NewService(prefix, m, nil)

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

		m.On("Put",
			context.Background(),
			prefix,
			mock.Anything,
			mock.Anything).
			Return(testCase.err)

		service := NewService(prefix, m, nil)
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
			data: [][]byte{[]byte(`{"id":"kube-name-1234"}`), []byte(`{"id":"56kube-name-5678"}`)},
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

		service := NewService(prefix, m, nil)

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

func TestService_InstallRelease(t *testing.T) {
	tcs := []struct {
		svc Service

		clusterID string
		rlsInput  *ReleaseInput

		expectedRes *release.Release
		expectedErr error
	}{
		{ // TC#1
			expectedErr: sgerrors.ErrNilEntity,
		},
		{ // TC#2
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: fakeChartGetter{
					err: errFake,
				},
			},
			expectedErr: errFake,
		},
		{ // TC#3
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: &fakeChartGetter{},
				storage: &storage.Fake{
					GetErr: errFake,
				},
			},
			expectedErr: errFake,
		},
		{ // TC#4
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: &fakeChartGetter{},
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return nil, errFake
				},
			},
			expectedErr: errFake,
		},
		{ // TC#5
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: &fakeChartGetter{},
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						err: errFake,
					}, nil
				},
			},
			expectedErr: errFake,
		},
		{ // TC#6
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: &fakeChartGetter{},
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						installRlsResp: &services.InstallReleaseResponse{
							Release: fakeRls,
						},
					}, nil
				},
			},
			expectedRes: fakeRls,
		},
	}

	for i, tc := range tcs {
		rls, err := tc.svc.InstallRelease(context.Background(), tc.clusterID, tc.rlsInput)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)

		if err == nil {
			require.Equalf(t, tc.expectedRes, rls, "TC#%d: check results", i+1)
		}
	}
}

func TestService_ListReleases(t *testing.T) {
	tcs := []struct {
		svc Service

		expectedRes []*model.ReleaseInfo
		expectedErr error
	}{
		{ // TC#1
			svc: Service{
				storage: &storage.Fake{
					GetErr: errFake,
				},
			},
			expectedErr: errFake,
		},
		{ // TC#2
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return nil, errFake
				},
			},
			expectedErr: errFake,
		},
		{ // TC#3
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						err: errFake,
					}, nil
				},
			},
			expectedErr: errFake,
		},
		{ // TC#4
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						listReleaseResp: &services.ListReleasesResponse{
							Releases: []*release.Release{fakeRls, nil},
						},
					}, nil
				},
			},
			expectedRes: []*model.ReleaseInfo{
				{
					Name:         fakeRls.GetName(),
					Namespace:    fakeRls.GetNamespace(),
					Version:      fakeRls.GetVersion(),
					CreatedAt:    timeconv.String(fakeRls.GetInfo().GetFirstDeployed()),
					LastDeployed: timeconv.String(fakeRls.GetInfo().GetLastDeployed()),
					Chart:        fakeRls.GetChart().Metadata.Name,
					ChartVersion: fakeRls.GetChart().Metadata.Version,
					Status:       fakeRls.GetInfo().Status.Code.String(),
				},
			},
		},
	}

	for i, tc := range tcs {
		rls, err := tc.svc.ListReleases(context.Background(), "testCluster", "", "", 0)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)

		if err == nil {
			require.Equalf(t, tc.expectedRes, rls, "TC#%d: check results", i+1)
		}
	}
}

func TestService_DeleteRelease(t *testing.T) {
	tcs := []struct {
		svc Service

		expectedRes *model.ReleaseInfo
		expectedErr error
	}{
		{ // TC#1
			svc: Service{
				storage: &storage.Fake{
					GetErr: errFake,
				},
			},
			expectedErr: errFake,
		},
		{ // TC#2
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return nil, errFake
				},
			},
			expectedErr: errFake,
		},
		{ // TC#3
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						err: errFake,
					}, nil
				},
			},
			expectedErr: errFake,
		},
		{ // TC#4
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						uninstReleaseResp: &services.UninstallReleaseResponse{
							Release: fakeRls,
						},
					}, nil
				},
			},
			expectedRes: &model.ReleaseInfo{
				Name:         fakeRls.GetName(),
				Namespace:    fakeRls.GetNamespace(),
				Version:      fakeRls.GetVersion(),
				CreatedAt:    timeconv.String(fakeRls.GetInfo().GetFirstDeployed()),
				LastDeployed: timeconv.String(fakeRls.GetInfo().GetLastDeployed()),
				Chart:        fakeRls.GetChart().Metadata.Name,
				ChartVersion: fakeRls.GetChart().Metadata.Version,
				Status:       fakeRls.GetInfo().Status.Code.String(),
			},
		},
	}

	for i, tc := range tcs {
		rls, err := tc.svc.DeleteRelease(context.Background(), "testCluster", "", true)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)

		if err == nil {
			require.Equalf(t, tc.expectedRes, rls, "TC#%d: check results", i+1)
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

		service := NewService("", m, nil)

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

func TestService_KubeConfigFor(t *testing.T) {
	testCases := []struct {
		user string

		kubeData   []byte
		getkubeErr error

		expectedErr error
	}{
		{
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			user:        "user",
			expectedErr: sgerrors.ErrNotFound,
		},
		{
			user:        KubernetesAdminUser,
			getkubeErr:  fakeErrFileNotFound,
			expectedErr: fakeErrFileNotFound,
		},
		{
			user:     KubernetesAdminUser,
			kubeData: []byte(`{"masters":{"m":{"publicIp":"1.2.3.4"}}}`),
		},
	}

	for i, tc := range testCases {
		m := new(testutils.MockStorage)
		m.On("Get", context.Background(), mock.Anything, mock.Anything).
			Return(tc.kubeData, tc.getkubeErr)

		svc := Service{
			storage: m,
		}

		data, err := svc.KubeConfigFor(context.Background(), "kname", tc.user)
		require.Equal(t, tc.expectedErr, errors.Cause(err), "TC#%d", i+1)

		if err == nil {
			require.NotNilf(t, data, "TC#%d", i+1)
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

		service := NewService(prefix, m, nil)

		_, err := service.GetCerts(context.Background(),
			testCase.kname, testCase.cname)

		if testCase.expectedErr != errors.Cause(err) {
			t.Errorf("Wrong error expected %v actual %v", testCase.expectedErr, err)
			return
		}
	}
}
