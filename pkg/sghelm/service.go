package sghelm

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/sghelm/repositories"
	"github.com/supergiant/supergiant/pkg/storage"
)

const (
	readmeFileName = "readme.md"

	repoPrefix = "/helm/repositories/"
)

var _ Servicer = &Service{}

// Servicer is an interface for the helm service.
type Servicer interface {
	CreateRepo(ctx context.Context, e *repo.Entry) (*model.RepositoryInfo, error)
	GetRepo(ctx context.Context, repoName string) (*model.RepositoryInfo, error)
	ListRepos(ctx context.Context) ([]model.RepositoryInfo, error)
	DeleteRepo(ctx context.Context, repoName string) (*model.RepositoryInfo, error)
	GetChartData(ctx context.Context, repoName, chartName, chartVersion string) (*model.ChartData, error)
	ListCharts(ctx context.Context, repoName string) ([]model.ChartInfo, error)
	GetChart(ctx context.Context, repoName, chartName, chartVersion string) (*chart.Chart, error)
}

// Service manages helm repositories.
type Service struct {
	storage storage.Interface
	repos   repositories.Interface
}

// NewService constructs a Service for helm repository.
func NewService(s storage.Interface) (*Service, error) {
	repos, err := repositories.New(repositories.DefaultHome)
	if err != nil {
		return nil, errors.Wrap(err, "setup repositories manager")
	}

	return &Service{
		storage: s,
		repos:   repos,
	}, nil
}

// CreateRepo stores a helm repository in the provided storage.
func (s Service) CreateRepo(ctx context.Context, e *repo.Entry) (*model.RepositoryInfo, error) {
	if e == nil {
		return nil, sgerrors.ErrNotFound
	}

	r, err := s.GetRepo(ctx, e.Name)
	if err == nil && r != nil {
		return nil, sgerrors.ErrAlreadyExists
	}

	ind, err := s.repos.GetIndexFile(e)
	if err != nil {
		return nil, errors.Wrap(err, "get repository index")
	}

	// store the index file
	r = toRepoInfo(e, ind)
	rawJSON, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "marshal index file")
	}
	if err = s.storage.Put(ctx, repoPrefix, e.Name, rawJSON); err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	return r, nil
}

// GetRepo retrieves the repository index file for provided nam.
func (s Service) GetRepo(ctx context.Context, repoName string) (*model.RepositoryInfo, error) {
	res, err := s.storage.Get(ctx, repoPrefix, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}
	// not found
	if res == nil {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "repo not found")
	}

	r := &model.RepositoryInfo{}
	if err = json.Unmarshal(res, r); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return r, nil
}

// ListRepos retrieves all helm repositories from the storage.
func (s Service) ListRepos(ctx context.Context) ([]model.RepositoryInfo, error) {
	rawRepos, err := s.storage.GetAll(ctx, repoPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	repos := make([]model.RepositoryInfo, len(rawRepos))
	for i, raw := range rawRepos {
		r := &model.RepositoryInfo{}
		err = json.Unmarshal(raw, r)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}
		repos[i] = *r
	}

	return repos, nil
}

// DeleteRepo removes a helm repository from the storage by its name.
func (s Service) DeleteRepo(ctx context.Context, repoName string) (*model.RepositoryInfo, error) {
	hrepo, err := s.GetRepo(ctx, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "get repository")
	}
	return hrepo, s.storage.Delete(ctx, repoPrefix, repoName)
}

func (s Service) GetChartData(ctx context.Context, repoName, chartName, chartVersion string) (*model.ChartData, error) {
	chrt, err := s.GetChart(ctx, repoName, chartName, chartVersion)
	if err != nil {
		return nil, err
	}
	return toChartData(chrt), nil
}

func (s Service) ListCharts(ctx context.Context, repoName string) ([]model.ChartInfo, error) {
	hrepo, err := s.GetRepo(ctx, repoName)
	if err != nil {
		return nil, errors.Wrapf(err, "get %s repository info", repoName)
	}

	return hrepo.Charts, nil
}

func (s Service) GetChart(ctx context.Context, repoName, chartName, chartVersion string) (*chart.Chart, error) {
	hrepo, err := s.GetRepo(ctx, repoName)
	if err != nil {
		return nil, errors.Wrapf(err, "get %s repository info", repoName)
	}
	ref, err := findChartURL(hrepo.Charts, chartName, chartVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "get %s(%s) chart", chartName, chartVersion)
	}

	chrt, err := s.repos.GetChart(hrepo.Config, ref)
	if err != nil {
		return nil, errors.Wrapf(err, "get %s chart", ref)
	}

	return chrt, nil
}

func toChartData(chrt *chart.Chart) *model.ChartData {
	if chrt == nil {
		return nil
	}

	out := &model.ChartData{
		Metadata: chrt.Metadata,
	}
	if chrt.Values != nil {
		out.Values = chrt.Values.Raw
	}
	if chrt.Files != nil {
		for _, f := range chrt.Files {
			if f != nil && strings.ToLower(f.TypeUrl) == readmeFileName {
				out.Readme = string(f.Value)
			}
		}
	}
	return out
}

func findChartURL(charts []model.ChartInfo, chartName, chartVersion string) (string, error) {
	for _, chrt := range charts {
		if chrt.Name != chartName {
			continue
		}
		if len(chrt.Versions) == 0 {
			break
		}
		chrtVer := findChartVersion(chrt.Versions, chartVersion)
		if len(chrtVer.URLs) != 0 {
			// charts are sorted
			return chrtVer.URLs[0], nil
		}
	}
	return "", sgerrors.ErrNotFound
}

func findChartVersion(chrtVers []model.ChartVersion, version string) model.ChartVersion {
	version = strings.TrimSpace(version)
	if len(chrtVers) > 0 && version == "" {
		return chrtVers[len(chrtVers)-1]
	}
	for _, v := range chrtVers {
		if v.Version == version {
			return v
		}
	}
	return model.ChartVersion{}
}

func toRepoInfo(e *repo.Entry, index *repo.IndexFile) *model.RepositoryInfo {
	if e == nil {
		return nil
	}

	r := &model.RepositoryInfo{
		Config: *e,
	}
	if index == nil {
		return r
	}

	r.Charts = make([]model.ChartInfo, 0, len(index.Entries))
	for name, entry := range index.Entries {
		if len(entry) == 0 {
			continue
		}

		sort.Sort(entry)
		if entry[0].Deprecated {
			continue
		}

		r.Charts = append(r.Charts, model.ChartInfo{
			Name:        name,
			Repo:        e.Name,
			Icon:        iconFrom(entry),
			Description: descriptionFrom(entry),
			Versions:    toChartVersions(entry),
		})
	}
	return r
}

func iconFrom(cvs repo.ChartVersions) string {
	// chartVersions are sorted, use the latest one
	if len(cvs) > 0 {
		return cvs[0].Icon
	}
	return ""
}

func descriptionFrom(cvs repo.ChartVersions) string {
	// chartVersions are sorted, use the latest one
	if len(cvs) > 0 {
		return cvs[0].Description
	}
	return ""
}

func toChartVersions(cvs repo.ChartVersions) []model.ChartVersion {
	if cvs == nil {
		return nil
	}
	chartVersions := make([]model.ChartVersion, 0, len(cvs))
	for _, cv := range cvs {
		chartVersions = append(chartVersions, model.ChartVersion{
			Version:    cv.Version,
			AppVersion: cv.AppVersion,
			Created:    cv.Created,
			Digest:     cv.Digest,
			URLs:       cv.URLs,
		})
	}
	return chartVersions
}
