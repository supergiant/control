package repositories

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"
)

var (
	DefaultHome = filepath.Join(os.TempDir(), ".helm")
)

var _ Interface = &Manager{}

// Interface represents an interface for the repositories manager.
type Interface interface {
	GetIndexFile(e *repo.Entry) (*repo.IndexFile, error)
	GetChart(conf repo.Entry, ref string) (*chart.Chart, error)
}

// Manager is responsible for dealing with helm repositories.
type Manager struct {
	helmHome helmpath.Home
}

// New is a constructor for helm Manager.
func New(homePath string) (*Manager, error) {
	m := &Manager{
		helmHome: helmpath.Home(homePath),
	}
	if err := m.ensureCacheDir(); err != nil {
		return nil, errors.Wrap(err, "setup cache dir")
	}
	return m, nil
}

// GetIndexFile retrieves IndexFile for the provided repository.
func (m Manager) GetIndexFile(conf *repo.Entry) (*repo.IndexFile, error) {
	if err := m.ensureCacheDir(); err != nil {
		return nil, err
	}

	cr, err := repo.NewChartRepository(conf, getter.All(environment.EnvSettings{}))
	if err != nil {
		return nil, errors.Wrap(err, "build chart repository")
	}
	if err = cr.DownloadIndexFile(m.helmHome.CacheIndex(conf.Name)); err != nil {
		return nil, errors.Wrap(err, "download index file")
	}
	ind, err := repo.LoadIndexFile(m.helmHome.CacheIndex(conf.Name))
	if err != nil {
		return nil, errors.Wrap(err, "load index file")
	}
	return ind, err
}

// GetChart retrieves a chart to from the remote repository and
// stores it to local cache. If chart exists locally it will be
// read from the cache.
func (m Manager) GetChart(conf repo.Entry, ref string) (*chart.Chart, error) {
	if err := m.ensureCacheDir(); err != nil {
		return nil, err
	}

	chrtPath := path.Join(m.helmHome.Archive(), path.Base(ref))
	chrt, err := chartutil.LoadFile(chrtPath)
	if err == nil {
		return chrt, nil
	}

	g, err := getter.NewHTTPGetter(ref, conf.CertFile, conf.KeyFile, conf.CAFile)
	if err != nil {
		return nil, errors.Wrap(err, "build a http client")
	}
	g.SetCredentials(conf.Username, conf.Password)

	r, err := g.Get(ref)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(chrtPath, r.Bytes(), 0644); err != nil {
		return nil, errors.Wrapf(err, "write %s chart", chrtPath)
	}

	log.Debugf("helm: manager: store %s chart to %s file", path.Base(ref), chrtPath)
	return chartutil.LoadFile(chrtPath)
}

// ensureCacheDir creates a filesystem tree like helm does if it
// doesn't exist. This is used for compatibility with helm libraries.
func (m Manager) ensureCacheDir() error {
	configDirectories := []string{
		m.helmHome.String(),
		m.helmHome.Repository(),
		m.helmHome.Cache(),
		m.helmHome.LocalRepository(),
		m.helmHome.Plugins(),
		m.helmHome.Starters(),
		m.helmHome.Archive(),
	}

	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			if err := os.MkdirAll(p, 0755); err != nil {
				return errors.Wrapf(err, "ould not create %s", p)
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("%s must be a directory", p)
		}
	}

	repoFile := m.helmHome.RepositoryFile()
	if fi, err := os.Stat(repoFile); err != nil {
		f := repo.NewRepoFile()
		if err := f.WriteFile(repoFile, 0644); err != nil {
			return errors.Wrapf(err, "write %s file", repoFile)
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", repoFile)
	}

	return nil
}
