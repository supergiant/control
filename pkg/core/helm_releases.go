package core

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/technosophos/moniker"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/timeconv"

	"github.com/supergiant/supergiant/pkg/model"
)

const (
	// HelmInstallTimeout is time in seconds to wait for any individual Kubernetes operation
	HelmInstallTimeout int64 = 300
)

type HelmReleases struct {
	Collection
}

func (c *HelmReleases) Populate() error {
	var kubes []*model.Kube
	if err := c.Core.DB.Preload("HelmReleases").Where("ready = ?", true).Find(&kubes); err != nil {
		return err
	}

	for _, kube := range kubes {
		newReleases, err := getHelmReleases(c.Core, kube)
		if err != nil {
			return err
		}

		oldReleases := kube.HelmReleases

		for _, newRelease := range newReleases {

			var oldRelease *model.HelmRelease
			oldIndex := 0

			for i, release := range oldReleases {
				if release.Name == newRelease.Name {
					oldRelease = release
					oldIndex = i
					break
				}
			}

			if oldRelease != nil {
				// remove from oldReleases
				oldReleases = append(oldReleases[:oldIndex], oldReleases[oldIndex+1:]...)

				// update chart if changed
				// if !reflect.DeepEqual(oldRelease, newRelease) {
				// NOTE we're not using the collection's Update method here to avoid immutability constraints
				if err := c.mergeUpdate(oldRelease.ID, oldRelease, newRelease); err != nil {
					return err
				}
				// }
			} else {
				// create new
				if err := c.Collection.Create(newRelease); err != nil {
					return errors.Wrapf(err, "create %s release", newRelease.Name)
				}
			}
		}

		for _, oldRelease := range oldReleases {
			if err := c.Core.DB.Delete(oldRelease); err != nil {
				return errors.Wrapf(err, "delete %s release", oldRelease.Name)
			}
		}
	}

	return nil
}

func (c *HelmReleases) Create(m *model.HelmRelease) error {
	// Generate Release name just like Helm does. We want to do this for our on
	// DB storage purposes -- relying on Helm for name can create issue with how
	// we sync Release records.
	if m.Name == "" {
		m.Name = moniker.New().NameSep("-")
	}

	if err := c.Collection.Create(m); err != nil {
		return err
	}

	action := &Action{
		Status: &model.ActionStatus{
			Description: "deploying",
			MaxRetries:  0,
		},
		Core: c.Core,
		// Nodes are needed to register with ELB on AWS
		Scope: c.Core.DB.Preload("Kube"),
		Model: m,
		ID:    m.ID,
		Fn: func(a *Action) error {
			hclient, err := c.Core.HelmClient(m.Kube)
			if err != nil {
				return errors.Wrap(err, "build helm client")
			}

			name := m.RepoName + "/" + m.ChartName
			chartPath, err := locateChartPath(name, m.ChartVersion)
			if err != nil {
				return errors.Wrap(err, "locate chart file")
			}

			chartRequested, err := chartutil.Load(chartPath)
			if err != nil {
				return errors.Wrap(err, "load chart config")
			}

			resp, err := hclient.InstallReleaseFromChart(
				chartRequested,
				m.Namespace,
				helm.ReleaseName(m.Name),
				helm.InstallWait(false),
				helm.InstallTimeout(HelmInstallTimeout),
			)
			if err != nil {
				return errors.Wrapf(err, "install %s chart", m.ChartName)
			}

			c.Core.Log.Debugf("DEBUG - Install %s chart: %s", m.ChartName, resp.Release.Info.Status.Code.String())
			return err
		},
	}
	return action.Async()
}

func (c *HelmReleases) Delete(id *int64, m *model.HelmRelease) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			if m.Name != "" {
				hclient, err := c.Core.HelmClient(m.Kube)
				if err != nil {
					return errors.Wrap(err, "build helm client")
				}

				resp, err := hclient.DeleteRelease(m.Name, helm.DeletePurge(true))
				if err != nil {
					return errors.Wrapf(err, "delete %s release", m.Name)
				}

				c.Core.Log.Debugf("DEBUG - Delete %s release: %s", m.Name, resp.Info)
			}
			return c.Collection.Delete(id, m)
		},
	}
}

func getHelmReleases(c *Core, kube *model.Kube) ([]*model.HelmRelease, error) {
	hclient, err := c.HelmClient(kube)
	if err != nil {
		return nil, errors.Wrap(err, "build helm client")
	}

	resp, err := hclient.ListReleases()
	if err != nil {
		return nil, errors.Wrap(err, "list releases")
	}

	var releases []*model.HelmRelease
	for _, r := range resp.GetReleases() {
		releases = append(releases, &model.HelmRelease{
			KubeName:     kube.Name,
			Name:         r.Name,
			Revision:     strconv.Itoa(int(r.Version)),
			UpdatedValue: timeconv.String(r.Info.LastDeployed),
			StatusValue:  r.Info.Status.Code.String(),
			// TODO this is not full ChartName (does not include Repo)
			ChartName:    r.Chart.GetMetadata().Name,
			ChartVersion: r.Chart.GetMetadata().Version,
		})
	}

	return releases, nil
}
