package api

import (
	"testing"

	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelmReposList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmRepos List works correctly", t, func() {

		table := []struct {
			// Input
			existingModels []*model.HelmRepo
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModels: []*model.HelmRepo{
					&model.HelmRepo{
						Name: "test",
						URL:  "www.website.com",
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			for _, existingModel := range item.existingModels {
				srv.Core.HelmRepos.Create(existingModel)
			}

			list := new(model.HelmRepoList)
			err := sg.HelmRepos.List(list)

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

func TestHelmReposCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmRepos Create works correctly", t, func() {

		table := []struct {
			// Input
			model *model.HelmRepo
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				model: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				err: nil,
			},

			// No name
			{
				model: &model.HelmRepo{
					URL: "www.website.com",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Name: zero value"},
			},

			// No URL
			{
				model: &model.HelmRepo{
					Name: "test",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: URL: zero value"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			err := sg.HelmRepos.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReposGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmRepos Get works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRepo
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModel: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.HelmRepos.Create(item.existingModel)

			err := sg.HelmRepos.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReposUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmRepos Update works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRepo
			modelUpdate   *model.HelmRepo
			// Expectations
			err *model.Error
		}{
			// Can't update Name
			{
				existingModel: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				modelUpdate: &model.HelmRepo{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update URL
			{
				existingModel: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				modelUpdate: &model.HelmRepo{
					URL: "www.website.biz",
				},
				err: &model.Error{Status: 422, Message: "URL cannot be changed"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.HelmRepos.Create(item.existingModel)

			err := sg.HelmRepos.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReposDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("HelmRepos Delete works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.HelmRepo
			// Expectations
			err *model.Error
		}{
			// Successful example
			{
				existingModel: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.HelmRepos.Create(item.existingModel)

			err := sg.HelmRepos.Delete(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}
