package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/helm/cmd/helm/search"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/model"
)

type HelmCharts struct {
	Collection
}

func (c *HelmCharts) Populate() error {
	var repos []*model.HelmRepo
	if err := c.Core.DB.Preload("Charts").Find(&repos); err != nil {
		return errors.Wrap(err, "find repositories")
	}

	for _, repoModel := range repos {
		results, err := searchHelmRepo(repoModel)
		if err != nil {
			return errors.Wrapf(err, "search for the %s repo", repoModel.Name)
		}

		chartsToDelete := repoModel.Charts

		for _, result := range results {

			var existingChart *model.HelmChart
			existingIndex := 0
			newChart := &model.HelmChart{
				Repo:        repoModel,
				RepoName:    repoModel.Name,
				Name:        result.Chart.Name,
				Version:     result.Chart.Version,
				Description: result.Chart.Description,
			}

			for i, chart := range repoModel.Charts {
				err := c.Get(chart.ID, new(model.HelmChart)) // TODO: Shuould we only look if null?
				if err != nil {
					fmt.Println("Chart Config Load Error:", err)
				}

				if chart.Version == newChart.Version && chart.Name == newChart.Name {
					existingChart = chart
					existingIndex = i
					break
				}
			}

			if existingChart != nil {
				// remove from chartsToDelete
				chartsToDelete = append(chartsToDelete[:existingIndex], chartsToDelete[existingIndex+1:]...)

				// if !reflect.DeepEqual(existingChart, newChart) {
				// update chart
				// NOTE we're not using the collection's Update method here to avoid immutability constraints
				if err := c.mergeUpdate(existingChart.ID, existingChart, newChart); err != nil {
					return err
				}
				// }
			} else {
				// create new
				c.Core.Log.Info("New Chart Version Detected for Chart:", newChart.Name)
				if err := c.Core.HelmCharts.Create(newChart); err != nil {
					return err
				}
			}
		}

		for _, chartToDelete := range chartsToDelete {
			if err := c.Core.HelmCharts.Delete(chartToDelete.ID, chartToDelete); err != nil {
				return err
			}
		}
	}

	return nil
}

//------------------------------------------------------------------------------

func (c *HelmCharts) Get(id *int64, m *model.HelmChart) error {
	if err := c.Collection.Get(id, m); err != nil {
		return err
	}
	// NOTE we're just letting this fail silently
	if err := c.loadConfig(m); err != nil {
		c.Core.Log.Warnf("Could not load default config for HelmChart '%s': %v", m.Name, err)
	}
	return nil
}

//------------------------------------------------------------------------------

func (c *HelmCharts) loadConfig(m *model.HelmChart) error {
	if len(m.DefaultConfigJSON) > 0 {
		return nil
	}

	if err := updateHelmRepoFile(c.Core); err != nil {
		return errors.Wrap(err, "update repo file")
	}

	nameWithRepo := m.RepoName + "/" + m.Name
	chartPath, err := locateChartPath(nameWithRepo, m.Version)
	if err != nil {
		return errors.Wrap(err, "locate chart path")
	}

	loadedChart, err := chartutil.Load(chartPath)
	if err != nil {
		return errors.Wrap(err, "load chart")
	}

	yamlDefaultConfig := loadedChart.GetValues().Raw

	if yamlDefaultConfig == "" {
		return nil
	}

	if err = yaml.Unmarshal([]byte(yamlDefaultConfig), &m.DefaultConfig); err != nil {
		return err
	}

	return c.Core.DB.Save(m)
}

//------------------------------------------------------------------------------

func helmHome() helmpath.Home {
	return helmpath.Home(filepath.Join(os.TempDir(), ".helm"))
}

func init() {
	home := helmHome()

	configDirectories := []string{
		home.String(),
		home.Repository(),
		home.Cache(),
		home.LocalRepository(),
		home.Plugins(),
		home.Starters(),
		home.Archive(),
	}

	// TODO: review repositories structure, get rid of the panics
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			if err := os.MkdirAll(p, 0755); err != nil {
				panic(fmt.Sprintf("Could not create %s: %s", p, err))
			}
		} else if !fi.IsDir() {
			panic(fmt.Sprintf("%s must be a directory", p))
		}
	}

	repoFile := home.RepositoryFile()
	if fi, err := os.Stat(repoFile); err != nil {
		f := repo.NewRepoFile()
		if err := f.WriteFile(repoFile, 0644); err != nil {
			panic(fmt.Sprintf("write %s file: %v", repoFile, err))
		}
	} else if fi.IsDir() {
		panic(fmt.Sprintf("%s must be a file, not a directory", repoFile))
	}

}

func searchHelmRepo(repoModel *model.HelmRepo) ([]*search.Result, error) {
	home := helmHome()

	f, err := repo.LoadRepositoriesFile(home.RepositoryFile())
	if err != nil {
		return nil, errors.Wrap(err, "load repositories file")
	}

	c := repo.Entry{
		Name:  repoModel.Name,
		URL:   repoModel.URL,
		Cache: home.CacheIndex(repoModel.Name),
	}

	r, err := repo.NewChartRepository(&c, getter.All(environment.EnvSettings{}))
	if err != nil {
		return nil, errors.Wrap(err, "build chart repository")
	}

	if err = r.DownloadIndexFile(home.Cache()); err != nil {
		return nil, errors.Wrap(err, "download index file")
	}

	f.Update(&c)
	err = f.WriteFile(home.RepositoryFile(), 0644)
	if err != nil {
		return nil, errors.Wrap(err, "write repository file")
	}

	ind, err := repo.LoadIndexFile(home.CacheIndex(repoModel.Name))
	if err != nil {
		return nil, errors.Wrap(err, "load index file")
	}

	// TODO: get rid of the index, instead use result of the LoadIndexFile()
	index := search.NewIndex()
	index.AddRepo(repoModel.Name, ind, false)

	return index.All(), nil
}

func updateHelmRepoFile(c *Core) error {
	home := helmHome()
	repoFilepath := home.RepositoryFile()

	r, err := repo.LoadRepositoriesFile(repoFilepath)
	if err != nil {
		r = repo.NewRepoFile()
	}

	var repos []*model.HelmRepo
	if err := c.DB.Find(&repos); err != nil {
		return err
	}

	for _, repoModel := range repos {
		if !r.Has(repoModel.Name) {
			r.Add(&repo.Entry{
				Name:  repoModel.Name,
				URL:   repoModel.URL,
				Cache: repoModel.Name + "-index.yaml",
			})
		}
	}

	return r.WriteFile(repoFilepath, 0644)
}

func locateChartPath(name, version string) (string, error) {
	home := helmHome()

	name = strings.TrimSpace(name)
	version = strings.TrimSpace(version)
	// Verify and return path if passed for name
	if _, err := os.Stat(name); err == nil {
		abs, err := filepath.Abs(name)
		if err != nil {
			return abs, err
		}
		return abs, nil
	}

	// Error out on invalid paths passed
	if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
		return name, fmt.Errorf("path %q not found", name)
	}

	// Look for repo in helm home and return if exists
	crepo := filepath.Join(home.Repository(), name)
	if _, err := os.Stat(crepo); err == nil {
		return filepath.Abs(crepo)
	}

	// Download chart to helm home
	dl := downloader.ChartDownloader{
		HelmHome: home,
		Getters:  getter.All(environment.EnvSettings{}),
	}
	filename, _, err := dl.DownloadTo(name, version, home.Archive())
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	}

	return filename, fmt.Errorf("file %q not found: %s", name, err)
}
