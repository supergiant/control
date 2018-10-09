package repositories

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"
)

type Manager struct {
	helmHome helmpath.Home
}

func New(homePath string) (*Manager, error) {
	m := &Manager{
		helmHome: helmpath.Home(homePath),
	}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m Manager) GetIndexFile(e *repo.Entry) (*repo.IndexFile, error) {
	cr, err := repo.NewChartRepository(e, getter.All(environment.EnvSettings{}))
	if err != nil {
		return nil, errors.Wrap(err, "build chart repository")
	}
	if err = cr.DownloadIndexFile(m.helmHome.CacheIndex(e.Name)); err != nil {
		return nil, errors.Wrap(err, "download index file")
	}
	ind, err := repo.LoadIndexFile(m.helmHome.CacheIndex(e.Name))
	if err != nil {
		return nil, errors.Wrap(err, "load index file")
	}
	return ind, err
}

func (m Manager) GetChart(conf repo.Entry, ref string) (*chart.Chart, error) {
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
