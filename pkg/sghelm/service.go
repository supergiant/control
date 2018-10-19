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
	repoPrefix = "/helm/repositories/"
)

var _ Servicer = &Service{}

// Servicer is an interface for the helm service.
type Servicer interface {
	CreateRepo(ctx context.Context, e *repo.Entry) (*model.Repository, error)
	GetRepo(ctx context.Context, repoName string) (*model.Repository, error)
	ListRepos(ctx context.Context) ([]model.Repository, error)
	DeleteRepo(ctx context.Context, repoName string) (*model.Repository, error)
	GetChartInfo(ctx context.Context, repoName, chartName string) (*model.Chart, error)
	ListChartInfos(ctx context.Context, repoName string) ([]model.Chart, error)
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
func (s Service) CreateRepo(ctx context.Context, e *repo.Entry) (*model.Repository, error) {
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
	r = toRepo(e, ind)
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
func (s Service) GetRepo(ctx context.Context, repoName string) (*model.Repository, error) {
	res, err := s.storage.Get(ctx, repoPrefix, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}
	// not found
	if res == nil {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "repo not found")
	}

	r := &model.Repository{}
	if err = json.Unmarshal(res, r); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return r, nil
}

// ListRepos retrieves all helm repositories from the storage.
func (s Service) ListRepos(ctx context.Context) ([]model.Repository, error) {
	rawRepos, err := s.storage.GetAll(ctx, repoPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	repos := make([]model.Repository, len(rawRepos))
	for i, raw := range rawRepos {
		r := &model.Repository{}
		err = json.Unmarshal(raw, r)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}
		repos[i] = *r
	}

	return repos, nil
}

// DeleteRepo removes a helm repository from the storage by its name.
func (s Service) DeleteRepo(ctx context.Context, repoName string) (*model.Repository, error) {
	hrepo, err := s.GetRepo(ctx, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "get repository")
	}
	return hrepo, s.storage.Delete(ctx, repoPrefix, repoName)
}

func (s Service) GetChartInfo(ctx context.Context, repoName, chartName string) (*model.Chart, error) {
	charts, err := s.ListChartInfos(ctx, repoName)
	if err != nil {
		return nil, err
	}

	for _, chrt := range charts {
		if chrt.Name == chartName {
			return &chrt, nil
		}
	}

	return nil, sgerrors.ErrNotFound
}

func (s Service) ListChartInfos(ctx context.Context, repoName string) ([]model.Chart, error) {
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

func findChartURL(charts []model.Chart, chartName, chartVersion string) (string, error) {
	for _, chrt := range charts {
		if chrt.Name != chartName {
			continue
		}
		if len(chrt.Versions) == 0 {
			break
		}
		chrtVer := findChartVersion(chrt.Versions, chartVersion)
		if len(chrtVer.URLs) != 0 {
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

func toRepo(e *repo.Entry, index *repo.IndexFile) *model.Repository {
	if e == nil {
		return nil
	}

	r := &model.Repository{
		Config: *e,
	}
	if index == nil {
		return r
	}

	r.Charts = make([]model.Chart, 0, len(index.Entries))
	for name, entry := range index.Entries {
		if len(entry) == 0 {
			continue
		}

		sort.Sort(entry)
		if entry[0].Deprecated {
			continue
		}

		r.Charts = append(r.Charts, model.Chart{
			Name:        name,
			Repo:        e.Name,
			Description: entry[0].Description,
			Home:        entry[0].Home,
			Keywords:    entry[0].Keywords,
			Maintainers: toMaintainers(entry[0].Maintainers),
			Sources:     entry[0].Sources,
			Icon:        entry[0].Icon,
			Versions:    toChartVersions(entry),
		})
	}
	return r
}

func toChartVersions(cvs repo.ChartVersions) []model.ChartVersion {
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

func toMaintainers(maintainers []*chart.Maintainer) []model.Maintainer {
	list := make([]model.Maintainer, 0, len(maintainers))
	for _, m := range maintainers {
		list = append(list, model.Maintainer{
			Name:  m.Name,
			Email: m.Email,
		})
	}
	return list
}
