package kube

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"

	"github.com/supergiant/supergiant/pkg/model"
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
		kubeData, _ := json.Marshal(testCase.kube)

		m.On("Put",
			context.Background(),
			prefix,
			testCase.kube.Name,
			kubeData).
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

		clusterName string
		rlsInput    *ReleaseInput

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
					err: fakeErr,
				},
			},
			expectedErr: fakeErr,
		},
		{ // TC#3
			rlsInput: &ReleaseInput{
				Name: "fake",
			},
			svc: Service{
				chrtGetter: &fakeChartGetter{},
				storage: &storage.Fake{
					GetErr: fakeErr,
				},
			},
			expectedErr: fakeErr,
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
					return nil, fakeErr
				},
			},
			expectedErr: fakeErr,
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
						err: fakeErr,
					}, nil
				},
			},
			expectedErr: fakeErr,
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
		rls, err := tc.svc.InstallRelease(context.Background(), tc.clusterName, tc.rlsInput)
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
					GetErr: fakeErr,
				},
			},
			expectedErr: fakeErr,
		},
		{ // TC#2
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return nil, fakeErr
				},
			},
			expectedErr: fakeErr,
		},
		{ // TC#3
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						err: fakeErr,
					}, nil
				},
			},
			expectedErr: fakeErr,
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
					Status:       fakeRls.GetInfo().Status.String(),
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
					GetErr: fakeErr,
				},
			},
			expectedErr: fakeErr,
		},
		{ // TC#2
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return nil, fakeErr
				},
			},
			expectedErr: fakeErr,
		},
		{ // TC#3
			svc: Service{
				storage: &storage.Fake{
					Item: []byte("{}"),
				},
				newHelmProxyFn: func(kube *model.Kube) (proxy.Interface, error) {
					return &fakeHelmProxy{
						err: fakeErr,
					}, nil
				},
			},
			expectedErr: fakeErr,
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
				Status:       fakeRls.GetInfo().Status.String(),
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
