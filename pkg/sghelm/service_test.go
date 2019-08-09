package sghelm

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
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

		expectedRepo *model.RepositoryInfo
		expectedErr  error
	}{
		{ // TC#1
			repoConf:    nil,
			expectedErr: sgerrors.ErrNilEntity,
		},
		{ // TC#2
			repoConf: &repo.Entry{
				Name: "storageError",
			},
			storage: fakeStorage{
				getErr: errFake,
			},
			expectedErr: errFake,
		},
		{ // TC#3
			repoConf: &repo.Entry{
				Name: "alreadyExists",
			},
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"alreadyExists"}}`),
			},
			expectedErr: sgerrors.ErrAlreadyExists,
		},
		{ // TC#4
			repoConf: &repo.Entry{
				Name: "getIndexFileError",
			},
			repos: fakeRepoManager{
				err: errFake,
			},
			expectedErr: errFake,
		},
		{ // TC#5
			repoConf: &repo.Entry{
				Name: "putError",
			},
			storage: fakeStorage{
				putErr: errFake,
			},
			expectedErr: errFake,
		},
		{ // TC#6
			repoConf: &repo.Entry{
				Name: "emptyIndex",
			},
			repos: fakeRepoManager{
				index: &repo.IndexFile{},
			},
			expectedRepo: &model.RepositoryInfo{
				Config: repo.Entry{
					Name: "emptyIndex",
				},
			},
		},
		{ // TC#7
			repoConf: &repo.Entry{
				Name: "success",
			},
			repos: fakeRepoManager{
				index: &repo.IndexFile{
					Entries: map[string]repo.ChartVersions{
						"chartDeprecated": {
							&repo.ChartVersion{
								Metadata: &chart.Metadata{
									Deprecated: true,
								},
							},
						},
						"chartNoMetadata": nil,
						"chartVersions": {
							&repo.ChartVersion{
								Metadata: &chart.Metadata{
									Name:        "chartVersions",
									Icon:        "chartVersions icon url",
									Description: "chartVersions description",
									Version:     "0.2.0",
								},
							},
							&repo.ChartVersion{
								Metadata: &chart.Metadata{
									Name:    "chartVersions",
									Version: "1.1.0",
								},
							},
						},
						"chartFake": {
							&repo.ChartVersion{
								Metadata: &chart.Metadata{
									Name:        "chartFake",
									Icon:        "chartFake icon url",
									Description: "chartFake description",
									Version:     "1.0.0",
									AppVersion:  "1.0.1",
								},
							},
						},
					},
				},
			},
			expectedRepo: &model.RepositoryInfo{
				Config: repo.Entry{
					Name: "success",
				},
				Charts: []model.ChartInfo{
					{
						Name:        "chartFake",
						Repo:        "success",
						Icon:        "chartFake icon url",
						Description: "chartFake description",
						Versions: []model.ChartVersion{
							{
								Version:    "1.0.0",
								AppVersion: "1.0.1",
							},
						},
					},
					{
						Name: "chartVersions",
						Repo: "success",
						Versions: []model.ChartVersion{
							{
								Version: "1.1.0",
							},
							{
								Version: "0.2.0",
							},
						},
					},
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

func TestService_UpdateRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	repoIndexFile := repo.IndexFile{
		Entries: map[string]repo.ChartVersions{
			"chartDeprecated": {
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Deprecated: true,
					},
				},
			},
			"chartNoMetadata": nil,
			"chartVersions": {
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Name:        "chartVersions",
						Icon:        "chartVersions icon url",
						Description: "chartVersions description",
						Version:     "0.2.0",
					},
				},
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Name:    "chartVersions",
						Version: "1.1.0",
					},
				},
			},
			"chartFake": {
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Name:        "chartFake",
						Icon:        "chartFake icon url",
						Description: "chartFake description",
						Version:     "1.0.0",
						AppVersion:  "1.0.1",
					},
				},
			},
		},
	}

	tcs := []struct {
		name string

		repoName string
		repoConf *repo.Entry

		storage fakeStorage
		repos   fakeRepoManager

		expectedRepo *model.RepositoryInfo
		expectedErr  error
	}{
		{
			name:     "storage_error",
			repoName: "s",
			storage: fakeStorage{
				getErr: errFake,
			},
			expectedErr: errFake,
		},
		{
			name:     "get_repo_index_error",
			repoName: "getRepo",
			repoConf: nil,
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"getRepo"}}`),
			},
			repos: fakeRepoManager{
				err: errFake,
			},
			expectedErr: errFake,
		},
		{
			name:     "put_repo_error",
			repoName: "getRepo",
			repoConf: nil,
			storage: fakeStorage{
				item:   []byte(`{"config":{"name":"getRepo"}}`),
				putErr: errFake,
			},
			repos:       fakeRepoManager{},
			expectedErr: errFake,
		},
		{
			name:     "update_repo_index",
			repoName: "updateRepoIndex",
			repoConf: nil,
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"updateRepoIndex"}}`),
			},
			repos: fakeRepoManager{
				index: &repoIndexFile,
			},
			expectedRepo: &model.RepositoryInfo{
				Config: repo.Entry{
					Name: "updateRepoIndex",
				},
				Charts: []model.ChartInfo{
					{
						Name:        "chartFake",
						Repo:        "updateRepoIndex",
						Icon:        "chartFake icon url",
						Description: "chartFake description",
						Versions: []model.ChartVersion{
							{
								Version:    "1.0.0",
								AppVersion: "1.0.1",
							},
						},
					},
					{
						Name: "chartVersions",
						Repo: "updateRepoIndex",
						Versions: []model.ChartVersion{
							{
								Version: "1.1.0",
							},
							{
								Version: "0.2.0",
							},
						},
					},
				},
			},
		},
		{
			name:     "update_repo",
			repoName: "updateRepoURL",
			repoConf: &repo.Entry{
				Name: "ignoreNewRepoName",
				URL:  "url",
			},
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"updateRepoURL"}}`),
			},
			expectedRepo: &model.RepositoryInfo{
				Config: repo.Entry{
					Name: "updateRepoURL",
					URL:  "url",
				},
			},
		},
	}

	for _, tc := range tcs {
		svc := Service{
			storage: &tc.storage,
			repos:   &tc.repos,
		}

		hrepo, err := svc.UpdateRepo(context.Background(), tc.repoName, tc.repoConf)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%s: check errors", tc.name)

		if err == nil {
			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%s: check results", tc.name)
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

		expectedRepo *model.RepositoryInfo
		expectedErr  error
	}{
		{ // TC#1
			repoName: "getError",
			storage: fakeStorage{
				getErr: errFake,
			},
			expectedErr: errFake,
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
			expectedRepo: &model.RepositoryInfo{
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

		expectedRepos []model.RepositoryInfo
		expectedErr   error
	}{
		{ // TC#1
			storage: fakeStorage{
				listErr: errFake,
			},
			expectedErr: errFake,
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
			expectedRepos: []model.RepositoryInfo{
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

		expectedRepo *model.RepositoryInfo
		expectedErr  error
	}{
		{ // TC#1
			repoName: "getError",
			storage: fakeStorage{
				getErr: errFake,
			},
			expectedErr: errFake,
		},
		{ // TC#2
			repoName: "deleteError",
			storage: fakeStorage{
				item:      []byte(`{"config":{"name":"success"}}`),
				deleteErr: errFake,
			},
			expectedErr: errFake,
		},
		{ // TC#4
			repoName: "success",
			storage: fakeStorage{
				item: []byte(`{"config":{"name":"success"}}`),
			},
			expectedRepo: &model.RepositoryInfo{
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

func Test_iconFrom(t *testing.T) {
	tcs := []struct {
		in       repo.ChartVersions
		expected string
	}{
		{},
		{
			in: repo.ChartVersions{
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Icon: "icon ref",
					},
				},
			},
			expected: "icon ref",
		},
	}

	for i, tc := range tcs {
		res := iconFrom(tc.in)
		require.Equalf(t, tc.expected, res, "TC#%d: check results", i+1)
	}
}

func Test_descriptionFrom(t *testing.T) {
	tcs := []struct {
		in       repo.ChartVersions
		expected string
	}{
		{},
		{
			in: repo.ChartVersions{
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Description: "description",
					},
				},
			},
			expected: "description",
		},
	}

	for i, tc := range tcs {
		res := descriptionFrom(tc.in)
		require.Equalf(t, tc.expected, res, "TC#%d: check results", i+1)
	}
}

func Test_toChartVersions(t *testing.T) {
	tcs := []struct {
		in       repo.ChartVersions
		expected []model.ChartVersion
	}{
		{},
		{
			in: repo.ChartVersions{},
		},
		{
			in: repo.ChartVersions{
				&repo.ChartVersion{
					Metadata: &chart.Metadata{
						Version:     "0.0.1",
						Description: "description",
					},
				},
			},
			expected: []model.ChartVersion{
				{
					Version: "0.0.1",
				},
			},
		},
	}

	for i, tc := range tcs {
		res := toChartVersions(tc.in)
		require.Equalf(t, tc.expected, res, "TC#%d: check results", i+1)
	}
}
