package helm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"sort"
)

const (
	repoPrefix       = "/helm/repositories/"
	chrtPrefixFormat = repoPrefix + "%s/charts/"

	cacheDir = "/tmp/helm/index.yaml"
)

// Service manages helm repositories.
type Service struct {
	storage storage.Interface
}

// NewService constructs a Service for helm repository.
func NewService(s storage.Interface) *Service {
	return &Service{
		storage: s,
	}
}

// Create stores a helm repository in the provided storage.
func (s *Service) Create(ctx context.Context, e *repo.Entry) error {
	if e == nil {
		return sgerrors.ErrNotFound
	}

	r, err := s.Get(ctx, e.Name)
	if err == nil && r != nil {
		return sgerrors.ErrAlreadyExists
	}

	//e.Cache = helmHome().CacheIndex(e.Name)
	cr, err := repo.NewChartRepository(e, getter.All(environment.EnvSettings{}))
	if err != nil {
		return errors.Wrap(err, "build chart repository")
	}
	if err = cr.DownloadIndexFile(""); err != nil {
		return errors.Wrap(err, "download index file")
	}
	if err = cr.Load(); err != nil {
		return errors.Wrap(err, "load charts")
	}

	// store the index file
	rawJSON, err := json.Marshal(helm.Repository{Config: *e, Index: *cr.IndexFile})
	if err != nil {
		return errors.Wrap(err, "marshal index file")
	}
	if err = s.storage.Put(ctx, repoPrefix, e.Name, rawJSON); err != nil {
		return errors.Wrap(err, "storage")
	}

	return nil
}

// Get retrieves the repository index file for provided nam.
func (s *Service) Get(ctx context.Context, repoName string) (*helm.Repository, error) {
	res, err := s.storage.Get(ctx, repoPrefix, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}
	// not found
	if res == nil {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "repo not found")
	}

	repo := &helm.Repository{}
	if err = json.Unmarshal(res, repo); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return repo, nil
}

// GetAll retrieves all helm repositories from the storage.
func (s *Service) GetAll(ctx context.Context) ([]helm.Repository, error) {
	rawRepos, err := s.storage.GetAll(ctx, repoPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "storage")
	}

	repos := make([]helm.Repository, len(rawRepos))
	for i, raw := range rawRepos {
		repo := &helm.Repository{}
		err = json.Unmarshal(raw, repo)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}
		repos[i] = *repo
	}

	return repos, nil
}

// Delete removes a helm repository from the storage by its name.
func (s *Service) Delete(ctx context.Context, repoName string) error {
	return s.storage.Delete(ctx, repoPrefix, repoName)
}

func toRepo(e repo.Entry, index repo.IndexFile) *helm.Repository {
	r := &helm.Repository{
		Config: e,
		Charts: make([]helm.Chart, 0, len(index.Entries)),
	}
	for name, entry := range index.Entries {
		if len(entry) == 0 {
			continue
		}

		sort.Sort(entry)
		if entry[0].Deprecated {
			continue
		}

		r.Charts = append(r.Charts, helm.Chart{
			Name:          name,
			Repo:          e.Name,
			Description:   entry[0].Description,
			Home:          entry[0].Home,
			Keywords:      entry[0].Keywords,
			Maintainers:   toMaintainers(entry[0].Maintainers),
			Sources:       entry[0].Sources,
			Icon:          entry[0].Icon,
			ChartVersions: toChartVersions(entry),
		})
	}
	return r
}

func toChartVersions(cvs repo.ChartVersions) []helm.ChartVersion {
	chartVersions := make([]helm.ChartVersion, 0, len(cvs))
	for _, cv := range cvs {
		chartVersions = append(chartVersions, helm.ChartVersion{
			Version:    cv.Version,
			AppVersion: cv.AppVersion,
			Created:    cv.Created,
			Digest:     cv.Digest,
			URLs:       cv.URLs,
		})
	}
	return chartVersions
}

func toMaintainers(maintainers []*chart.Maintainer) []helm.Maintainer {
	list := make([]helm.Maintainer, 0, len(maintainers))
	for _, m := range maintainers {
		list = append(list, helm.Maintainer{
			Name:  m.Name,
			Email: m.Email,
		})
	}
	return list
}
