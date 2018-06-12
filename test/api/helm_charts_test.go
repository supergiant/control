package api

import (
	"testing"

	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelmChartsList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmCharts List works correctly", t, func() {

		table := []struct {
			// Input
			parentRepo     *model.HelmRepo
			existingModels []*model.HelmChart
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				existingModels: []*model.HelmChart{
					{
						RepoName:    "test",
						Name:        "test",
						Version:     "0.1.0",
						Description: "A chart",
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.HelmRepos.Create(item.parentRepo)

			for _, existingModel := range item.existingModels {
				srv.Core.HelmCharts.Create(existingModel)
			}

			list := new(model.HelmChartList)
			err := sg.HelmCharts.List(list)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			So(len(list.Items), ShouldEqual, len(item.existingModels))
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmChartsCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmCharts Create works correctly", t, func() {

		table := []struct {
			// Input
			parentRepo *model.HelmRepo
			model      *model.HelmChart
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				model: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: nil,
			},

			// No RepoName
			{
				model: &model.HelmChart{
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: RepoName: zero value"},
			},

			// Repo does not exist
			{
				model: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'HelmRepoName' on HelmChart"},
			},

			// No Name
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				model: &model.HelmChart{
					RepoName:    "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Name: zero value"},
			},

			// No Version
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				model: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Description: "A chart",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Version: zero value"},
			},

			// No Description
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				model: &model.HelmChart{
					RepoName: "test",
					Name:     "test",
					Version:  "0.1.0",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Description: zero value"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentRepo != nil {
				srv.Core.HelmRepos.Create(item.parentRepo)
			}

			err := sg.HelmCharts.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmChartsGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmCharts Get works correctly", t, func() {

		table := []struct {
			// Input
			parentRepo    *model.HelmRepo
			existingModel *model.HelmChart
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				existingModel: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentRepo != nil {
				srv.Core.HelmRepos.Create(item.parentRepo)
			}

			srv.Core.HelmCharts.Create(item.existingModel)

			err := sg.HelmCharts.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmChartsUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmCharts Update works correctly", t, func() {

		table := []struct {
			// Input
			parentRepo    *model.HelmRepo
			existingModel *model.HelmChart
			modelUpdate   *model.HelmChart
			// Expectations
			err *model.Error
		}{
			// Can't update RepoName
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				existingModel: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				modelUpdate: &model.HelmChart{
					RepoName: "new-repo",
				},
				err: &model.Error{Status: 422, Message: "RepoName cannot be changed"},
			},

			// Can't update Name
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				existingModel: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				modelUpdate: &model.HelmChart{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentRepo != nil {
				srv.Core.HelmRepos.Create(item.parentRepo)
			}

			srv.Core.HelmCharts.Create(item.existingModel)

			err := sg.HelmCharts.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmChartsDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmCharts Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentRepo    *model.HelmRepo
			existingModel *model.HelmChart
			// Expectations
			err *model.Error
		}{
			// Successful example
			{
				parentRepo: &model.HelmRepo{
					Name: "test",
					URL:  "www.test.com",
				},
				existingModel: &model.HelmChart{
					RepoName:    "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "A chart",
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentRepo != nil {
				srv.Core.HelmRepos.Create(item.parentRepo)
			}

			srv.Core.HelmCharts.Create(item.existingModel)

			err := sg.HelmCharts.Delete(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}
