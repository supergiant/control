package repos

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
)

type Manager struct {
	helmHome helmpath.Home
}

func New(homePath string) *Manager {
	return &Manager{
		helmHome: helmpath.Home(homePath),
	}
}

func (m Manager) GetIndexFile(e *repo.Entry) (*repo.IndexFile, error) {
	cr, err := repo.NewChartRepository(e, getter.All(environment.EnvSettings{}))
	if err != nil {
		return nil, errors.Wrap(err, "build chart repository")
	}
	if err = cr.DownloadIndexFile(m.helmHome.Cache()); err != nil {
		return nil, errors.Wrap(err, "download index file")
	}
	ind, err := repo.LoadIndexFile(m.helmHome.CacheIndex(e.Name))
	if err != nil {
		return nil, errors.Wrap(err, "load index file")
	}
	return ind, err
}

// init creates a filesystem tree like helm does.
// This is used for compatibility with helm libraries.
func (m Manager) init() error {
	configDirectories := []string{
		m.helmHome.String(),
		m.helmHome.Repository(),
		m.helmHome.Cache(),
		m.helmHome.LocalRepository(),
		m.helmHome.Plugins(),
		m.helmHome.Starters(),
		m.helmHome.Archive(),
	}

	// TODO: review repositories structure, get rid of the panics
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			if err := os.MkdirAll(p, 0755); err != nil {
				return errors.Wrapf(err, "ould not create %s", p)
			}
		} else if !fi.IsDir() {
			return errors.New(fmt.Sprintf("%s must be a directory", p))
		}
	}

	repoFile := m.helmHome.RepositoryFile()
	if fi, err := os.Stat(repoFile); err != nil {
		f := repo.NewRepoFile()
		if err := f.WriteFile(repoFile, 0644); err != nil {
			return errors.Wrapf(err, "write %s file", repoFile)
		}
	} else if fi.IsDir() {
		return errors.New(fmt.Sprintf("%s must be a file, not a directory", repoFile))
	}

	return nil
}
