package api

import (
	"errors"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

var helmReleaseID int64 = 100

func TestHelmReleasesList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("HelmReleases List works correctly", t, func() {

		table := []struct {
			// Input
			existingModels []*model.HelmRelease
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModels: []*model.HelmRelease{
					{
						KubeName:     kube.Name,
						RepoName:     "stable",
						ChartName:    "redis",
						ChartVersion: "0.1.0",
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			// wipeAndInitialize(srv.Core)

			for _, existingModel := range item.existingModels {
				srv.Core.HelmReleases.Create(existingModel)
			}

			list := new(model.HelmReleaseList)
			err := sg.HelmReleases.List(list)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			So(len(list.Items), ShouldEqual, len(item.existingModels))

			// NOTE we have to clean up HelmReleases manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.HelmRelease{})
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReleasesCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("HelmReleases Create works correctly", t, func() {

		table := []struct {
			// Input
			repos []*model.HelmRepo
			model *model.HelmRelease
			// Mocks
			mockKubeCreateResourceError error
			mockKubeGetResourceFn       func(apiVersion, kind, namespace, name string, out interface{}) error
			// Expectations
			fullCommand string
			err         *model.Error
			asyncErr    string // because Create has a sync and async phase
		}{
			// A successful example
			{
				repos: nil,
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
					Config: map[string]interface{}{
						"nested": map[string]interface{}{
							"key": "value",
						},
						"not-nested": 6,
					},
				},
				mockKubeCreateResourceError: nil,
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					return errors.New("404") // means job finishes successfully
				},
				fullCommand: `/helm init --client-only && /helm install stable/redis --set nested.key='value',not-nested=6 --version 0.1.0 --name test`,
				err:         nil,
				asyncErr:    "",
			},

			// A successful example (with Repos)
			{
				repos: []*model.HelmRepo{
					{
						Name: "supergiant",
						URL:  "www.website.com",
					},
				},
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "supergiant",
					ChartName:    "elasticsearch",
					ChartVersion: "0.7.1",
				},
				mockKubeCreateResourceError: nil,
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					return errors.New("404") // means job finishes successfully
				},
				fullCommand: `/helm init --client-only && /helm repo add supergiant www.website.com && /helm install supergiant/elasticsearch --version 0.7.1 --name test`,
				err:         nil,
				asyncErr:    "",
			},

			// A successful example (when catching Pod phase succeeded)
			{
				repos: nil,
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "supergiant",
					ChartName:    "elasticsearch",
					ChartVersion: "0.7.1",
				},
				mockKubeCreateResourceError: nil,
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					pod := out.(*kubernetes.Pod)
					pod.Status.Phase = "Succeeded"
					return nil
				},
				fullCommand: `/helm init --client-only && /helm install supergiant/elasticsearch --version 0.7.1 --name test`,
				err:         nil,
				asyncErr:    "",
			},

			// Timeout (Pod never starts)
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "supergiant",
					ChartName:    "elasticsearch",
					ChartVersion: "0.7.1",
				},
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					pod := out.(*kubernetes.Pod)
					pod.Status.Phase = "Pending"
					return nil
				},
				fullCommand: `/helm init --client-only && /helm install supergiant/elasticsearch --version 0.7.1 --name test`,
				asyncErr:    "Timed out waiting for Helm cmd 'install supergiant/elasticsearch --version 0.7.1 --name test'",
			},

			// Unexpected error on Kubernetes CreateResource
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "supergiant",
					ChartName:    "elasticsearch",
					ChartVersion: "0.7.1",
				},
				mockKubeCreateResourceError: errors.New("something unexpected"),
				fullCommand:                 `/helm init --client-only && /helm install supergiant/elasticsearch --version 0.7.1 --name test`,
				asyncErr:                    "Error creating Pod: something unexpected",
			},

			// Unexpected error on Kubernetes GetResource
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "supergiant",
					ChartName:    "elasticsearch",
					ChartVersion: "0.7.1",
				},
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					return errors.New("something unexpected")
				},
				fullCommand: `/helm init --client-only && /helm install supergiant/elasticsearch --version 0.7.1 --name test`,
				asyncErr:    "Error GETting Pod: something unexpected",
			},

			// No KubeName
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},

			// Non-existent Kube
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     "crub",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'KubeName' on HelmRelease"},
			},

			// No RepoName
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: RepoName: zero value"},
			},

			// No ChartName
			{
				model: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					RepoName:     "stable",
					ChartVersion: "0.1.0",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: ChartName: zero value"},
			},

			// No ChartVersion
			{
				model: &model.HelmRelease{
					BaseModel: model.BaseModel{ID: &helmReleaseID},
					KubeName:  kube.Name,
					RepoName:  "stable",
					ChartName: "redis",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: ChartVersion: zero value"},
			},
		}

		for _, item := range table {

			// wipeAndInitialize(srv.Core)

			var fullCommand string

			srv.Core.HelmJobStartTimeout = time.Nanosecond

			srv.Core.K8S = func(_ *model.Kube) kubernetes.ClientInterface {
				return &fake_core.KubernetesClient{
					CreateResourceFn: func(apiVersion, kind, namespace string, in, out interface{}) error {
						pod := in.(*kubernetes.Pod)
						fullCommand = pod.Spec.Containers[0].Args[0]
						return item.mockKubeCreateResourceError
					},
					GetResourceFn: item.mockKubeGetResourceFn,
				}
			}

			for _, repo := range item.repos {
				srv.Core.HelmRepos.Create(repo)
			}

			err := sg.HelmReleases.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// Reload to get new status
			// We use fresh to get rid of status (which was set by Create, and is nil at Get)
			freshModel := new(model.HelmRelease)
			sg.HelmReleases.Get(item.model.ID, freshModel)

			if item.asyncErr == "" {
				So(freshModel.Status, ShouldBeNil)
			} else {
				So(freshModel.Status.Error, ShouldEqual, item.asyncErr)
			}

			So(fullCommand, ShouldEqual, item.fullCommand)

			// NOTE we have to clean up HelmReleases manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.HelmRelease{})
			srv.Core.DB.Delete(&model.HelmRepo{})
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReleasesGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("HelmReleases Get works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRelease
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				err: nil,
			},
		}

		for _, item := range table {

			srv.Core.HelmReleases.Create(item.existingModel)

			err := sg.HelmReleases.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// NOTE we have to clean up HelmReleases manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.HelmRelease{})
			srv.Core.DB.Delete(&model.HelmRepo{})
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReleasesUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("HelmReleases Update works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRelease
			modelUpdate   *model.HelmRelease
			// Expectations
			err *model.Error
		}{
			// Can't update KubeName
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				modelUpdate: &model.HelmRelease{
					KubeName: "butt",
				},
				err: &model.Error{Status: 422, Message: "KubeName cannot be changed"},
			},

			// Can't update Name
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				modelUpdate: &model.HelmRelease{
					Name: "butt",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update RepoName
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				modelUpdate: &model.HelmRelease{
					RepoName: "butt",
				},
				err: &model.Error{Status: 422, Message: "RepoName cannot be changed"},
			},

			// Can't update ChartName
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				modelUpdate: &model.HelmRelease{
					ChartName: "butt",
				},
				err: &model.Error{Status: 422, Message: "ChartName cannot be changed"},
			},

			// Can't update ChartVersion
			{
				existingModel: &model.HelmRelease{
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				modelUpdate: &model.HelmRelease{
					ChartVersion: "7.7.7",
				},
				err: &model.Error{Status: 422, Message: "ChartVersion cannot be changed"},
			},
		}

		for _, item := range table {

			srv.Core.HelmReleases.Create(item.existingModel)

			err := sg.HelmReleases.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// NOTE we have to clean up HelmReleases manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.HelmRelease{})
			srv.Core.DB.Delete(&model.HelmRepo{})
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReleasesDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("HelmReleases Delete works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRelease
			// Mocks
			mockKubeCreateResourceError error
			mockKubeGetResourceFn       func(apiVersion, kind, namespace, name string, out interface{}) error
			// Expectations
			fullCommand string
			asyncErr    string
		}{
			// A successful example
			{
				existingModel: &model.HelmRelease{
					BaseModel:    model.BaseModel{ID: &helmReleaseID},
					KubeName:     kube.Name,
					Name:         "test",
					RepoName:     "stable",
					ChartName:    "redis",
					ChartVersion: "0.1.0",
				},
				mockKubeCreateResourceError: nil,
				mockKubeGetResourceFn: func(apiVersion, kind, namespace, name string, out interface{}) error {
					return errors.New("404") // means job finishes successfully
				},
				fullCommand: `/helm init --client-only && /helm delete test --purge`,
				asyncErr:    "",
			},

			// NOTE We don't really need many other scenarios tested, since anything
			// else would test execHelmCmd, which we test in Create.

		}

		for _, item := range table {

			var fullCommand string

			srv.Core.HelmJobStartTimeout = time.Nanosecond

			srv.Core.K8S = func(_ *model.Kube) kubernetes.ClientInterface {
				return &fake_core.KubernetesClient{
					CreateResourceFn: func(apiVersion, kind, namespace string, in, out interface{}) error {
						pod := in.(*kubernetes.Pod)
						fullCommand = pod.Spec.Containers[0].Args[0]
						return item.mockKubeCreateResourceError
					},
					GetResourceFn: item.mockKubeGetResourceFn,
				}
			}

			srv.Core.HelmReleases.Create(item.existingModel)

			// NOTE no need to test this error here
			_ = sg.HelmReleases.Delete(item.existingModel.ID, item.existingModel)

			time.Sleep(10 * time.Millisecond)

			// Reload to get new status
			// We use fresh to get rid of status (which was set by Delete, and is nil at Get)
			freshModel := new(model.HelmRelease)
			getErr := sg.HelmReleases.Get(item.existingModel.ID, freshModel)

			if item.asyncErr == "" {
				So(getErr, ShouldNotBeNil)
			} else {
				So(freshModel.Status.Error, ShouldEqual, item.asyncErr)
			}

			So(fullCommand, ShouldEqual, item.fullCommand)

			// NOTE we have to clean up HelmReleases manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.HelmRelease{})
			srv.Core.DB.Delete(&model.HelmRepo{})
		}
	})
}
