package core

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/technosophos/moniker"
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
					return err
				}
			}
		}

		for _, oldRelease := range oldReleases {
			if err := c.Core.DB.Delete(oldRelease); err != nil {
				return err
			}
		}
	}

	return nil
}

//------------------------------------------------------------------------------

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
			cmd := fmt.Sprintf("install %s/%s", m.RepoName, m.ChartName)
			if len(m.Config) > 0 {
				cmd += fmt.Sprintf(" --set %s", strings.Replace(releaseConfigAsFlagValue(m.Config, ""), ",,", ",", -1)) // TODO: This will remove double , but what if we have 3?
			}
			if m.ChartVersion != "" {
				cmd += " --version " + m.ChartVersion
			}
			if m.Name != "" {
				cmd += " --name " + m.Name
			}
			if m.Namespace != "" {
				cmd += " --namespace " + m.Namespace
			}

			c.Core.Log.Debug("DEBUG - Chart install cmd:", cmd)
			_, err := execHelmCmd(c.Core, m.Kube, cmd)
			return err
		},
	}
	return action.Async()
}

//------------------------------------------------------------------------------

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
				cmd := fmt.Sprintf("delete %s --purge", m.Name)
				_, err := execHelmCmd(c.Core, m.Kube, cmd)
				if err != nil && !strings.Contains(err.Error(), "Error: release: not found") {
					return err
				}
			}
			return c.Collection.Delete(id, m)
		},
	}
}

//------------------------------------------------------------------------------

func getHelmReleases(c *Core, kube *model.Kube) ([]*model.HelmRelease, error) {
	log, err := execHelmCmd(c, kube, "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(log, "\n")

	var headerExists bool
	var headerLineNo int

	for lineNo, line := range lines {
		if regexp.MustCompile(`^NAME`).MatchString(line) {
			headerExists = true
			headerLineNo = lineNo
			break
		}
	}

	if !headerExists {
		return nil, nil
	}

	var releases []*model.HelmRelease
	for _, line := range lines[headerLineNo+1:] {
		cols := regexp.MustCompile(" *\t").Split(line, -1)

		// NOTE this should really only be the last line causing this to trigger
		if len(cols) < 5 {
			continue
		}

		chartNameSplit := strings.Split(cols[4], "-")

		release := &model.HelmRelease{
			KubeName:     kube.Name,
			Name:         cols[0],
			Revision:     cols[1],
			UpdatedValue: cols[2],
			StatusValue:  cols[3],
			// TODO this is not full ChartName (does not include Repo)
			ChartName:    chartNameSplit[0],
			ChartVersion: chartNameSplit[1],
		}
		releases = append(releases, release)
	}

	return releases, nil
}

//------------------------------------------------------------------------------

// TODO really there should be one for each Kube, but it is hard to prevent race condition in mutex creation
var globalHelmCmdMutex = new(sync.Mutex)

func execHelmCmd(c *Core, kube *model.Kube, cmd string) (out string, err error) {

	globalHelmCmdMutex.Lock()
	defer globalHelmCmdMutex.Unlock()

	var repos []*model.HelmRepo
	if err = c.DB.Find(&repos); err != nil {
		return
	}
	var repoAddCmds []string
	for _, repo := range repos {
		repoAddCmd := fmt.Sprintf("/helm repo add %s %s", repo.Name, repo.URL)
		repoAddCmds = append(repoAddCmds, repoAddCmd)
	}
	repoAddCmd := strings.Join(repoAddCmds, " && ")

	fullCmd := "/helm init --client-only"
	if repoAddCmd != "" {
		fullCmd += " && " + repoAddCmd
	}
	fullCmd += " && /helm " + cmd

	jobLabel := util.RandomString(16)

	pod := &kubernetes.Pod{
		Metadata: kubernetes.Metadata{
			Name: "supergiant-helm-job",
			Labels: map[string]string{
				"app": "supergiant-helm-job",
				"job": jobLabel,
			},
		},
		Spec: kubernetes.PodSpec{
			NodeSelector: map[string]string{
				"beta.kubernetes.io/arch": "amd64",
			},
			Containers: []kubernetes.Container{
				{
					Name:  "helm-worker",
					Image: "supergiant/helm-worker:v2.8.2",
					// ImagePullPolicy: "Always",
					Command: []string{"/bin/sh", "-c"},
					Args:    []string{fullCmd},
				},
			},
			RestartPolicy: "Never",
		},
	}

	if err = c.K8S(kube).CreateResource("api/v1", "Pod", "default", pod, pod); err != nil {
		err = errors.New("Error creating Pod: " + err.Error())
		return
	}

	podName := pod.Metadata.Name

	defer c.K8S(kube).DeleteResource("api/v1", "Pod", "default", podName)

	// TODO(stgleb): Context should be inherited from higher level context
	ctx, cancel := context.WithTimeout(context.Background(), c.HelmJobStartTimeout)
	defer cancel()

	waitErr := util.WaitFor(ctx, fmt.Sprintf("Helm cmd '%s'", cmd), 1*time.Second, func() (bool, error) {
		if err = c.K8S(kube).GetResource("api/v1", "Pod", "default", podName, pod); err != nil {
			if strings.Contains(err.Error(), "404") {
				// This or the Phase == "Succeeded" line may fire, but this one is much
				// less likely. The pod seems to linger for a while as we capture the Status
				return true, nil // done
			}
			err = errors.New("Error GETting Pod: " + err.Error())
			return false, err
		}

		// Get log
		out, _ = c.K8S(kube).GetPodLog("default", podName)

		if pod.Status.Phase == "Failed" {
			return false, fmt.Errorf("Helm cmd failed: %s\n\n%s", cmd, out)
		}
		if pod.Status.Phase == "Succeeded" {
			return true, nil // good to go
		}

		return false, nil // pod still pending / running
	})

	return out, waitErr
}

//------------------------------------------------------------------------------

func releaseConfigAsFlagValue(config map[string]interface{}, parent string) (fv string) {

	if parent != "" {
		parent += "."
	}

	for key, value := range config {

		if value == nil {
			value = ""
		}

		fullKey := parent + key

		if fv != "" {
			fv += ","
		}

		switch reflect.TypeOf(value).Kind() {

		case reflect.Map:
			fv += releaseConfigAsFlagValue(value.(map[string]interface{}), fullKey)

		case reflect.String:
			fv += fmt.Sprintf(`%s='%v'`, fullKey, value)

		default:
			fv += fmt.Sprintf(`%s=%v`, fullKey, value)
		}
	}
	return
}
