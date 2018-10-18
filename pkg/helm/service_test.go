package helm

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type fakeRepoManager struct {
	index *repo.IndexFile
	chrt  *chart.Chart
	err   error
}

func (m fakeRepoManager) GetIndexFile(e *repo.Entry) (*repo.IndexFile, error) {
	return m.index, m.err
}
func (m fakeRepoManager) GetChart(conf repo.Entry, ref string) (*chart.Chart, error) {
	return m.chrt, m.err
}

type fakeStorage struct {
	item      []byte
	items     [][]byte
	putErr    error
	getErr    error
	listErr   error
	deleteErr error
}

func (s fakeStorage) Put(ctx context.Context, prefix string, key string, value []byte) error {
	return s.putErr
}

func (s fakeStorage) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	return s.item, s.getErr
}

func (s fakeStorage) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	return s.items, s.listErr
}

func (s fakeStorage) Delete(ctx context.Context, prefix string, key string) error {
	return s.deleteErr
}

func TestService_CreateRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		repoConf *repo.Entry

		storage fakeStorage
		repos   fakeRepoManager

		expectedRepo *helm.Repository
		expectedErr  error
	}{
		{ // TC#1
			repoConf:    nil,
			expectedErr: sgerrors.ErrNotFound,
		},
		{ // TC#2
			repoConf: &repo.Entry{
				Name: "alreadyExists",
			},
			storage: fakeStorage{
				item: []byte(`{"name":"alreadyExists"}`),
			},
			expectedErr: sgerrors.ErrAlreadyExists,
		},
		{ // TC#3
			repoConf: &repo.Entry{
				Name: "getIndexFileError",
			},
			repos: fakeRepoManager{
				err: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#4
			repoConf: &repo.Entry{
				Name: "putError",
			},
			storage: fakeStorage{
				putErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#5
			repoConf: &repo.Entry{
				Name: "success",
			},
			repos: fakeRepoManager{},
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "success",
				},
			},
		},
	}

	for i, tc := range tcs {
		svc := Service{
			storage: &tc.storage,
			repos:   &tc.repos,
		}

		hrepo, err := svc.CreateRepo(context.Background(), tc.repoConf)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)

		if err == nil {
			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check results", i+1)
		}
	}
}

func TestService_GetRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		repoName string
		storage  fakeStorage

		expectedRepo *helm.Repository
		expectedErr  error
	}{
		{ // TC#1
			repoName: "getError",
			storage: fakeStorage{
				getErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#2
			repoName:    "notFound",
			expectedErr: sgerrors.ErrNotFound,
		},
		{ // TC#3
			repoName: "decodeError",
			storage: fakeStorage{
				item: []byte(`{}}`),
			},
		},
		{ // TC#4
			repoName: "success",
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"success"}}`),
			},
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "success",
				},
			},
		},
	}

	for i, tc := range tcs {
		svc := Service{
			storage: &tc.storage,
		}

		hrepo, err := svc.GetRepo(context.Background(), tc.repoName)
		if tc.expectedErr != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)
		} else {
			require.Nilf(t, tc.expectedErr, "TC#%d: no errors", i+1)
		}

		if err == nil {
			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check results", i+1)
		}
	}
}

func TestService_ListRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		storage fakeStorage

		expectedRepos []helm.Repository
		expectedErr   error
	}{
		{ // TC#1
			storage: fakeStorage{
				listErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#2: unmarshal error
			storage: fakeStorage{
				items: [][]byte{[]byte(`{}}`)},
			},
		},
		{ // TC#3
			storage: fakeStorage{
				items: [][]byte{[]byte(`{"config":{"name":"success"}}`)},
			},
			expectedRepos: []helm.Repository{
				{
					Config: repo.Entry{
						Name: "success",
					},
				},
			},
		},
	}

	for i, tc := range tcs {
		svc := Service{
			storage: &tc.storage,
		}

		hrepo, err := svc.ListRepos(context.Background())
		if tc.expectedErr != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)
		} else {
			require.Nilf(t, tc.expectedErr, "TC#%d: no errors", i+1)
		}

		if err == nil {
			require.Equalf(t, tc.expectedRepos, hrepo, "TC#%d: check results", i+1)
		}
	}
}

func TestService_DeleteRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		repoName string
		storage  fakeStorage

		expectedRepo *helm.Repository
		expectedErr  error
	}{
		{ // TC#1
			repoName: "getError",
			storage: fakeStorage{
				getErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#2
			repoName: "deleteError",
			storage: fakeStorage{
				item:      []byte(`{"config":{"name":"success"}}`),
				deleteErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{ // TC#4
			repoName: "success",
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"success"}}`),
			},
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "success",
				},
			},
		},
	}

	for i, tc := range tcs {
		svc := Service{
			storage: &tc.storage,
		}

		hrepo, err := svc.DeleteRepo(context.Background(), tc.repoName)
		if tc.expectedErr != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: check errors", i+1)
		} else {
			require.Nilf(t, tc.expectedErr, "TC#%d: no errors", i+1)
		}

		if err == nil {
			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check results", i+1)
		}
	}
}
