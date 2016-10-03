package api

import (
	"errors"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/test/fake_core"
	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCloudAccountsList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("CloudAccounts List works correctly", t, func() {

		table := []struct {
			// Input
			existingModels []*model.CloudAccount
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModels: []*model.CloudAccount{
					&model.CloudAccount{
						Name:        "test",
						Provider:    "aws",
						Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			// For ValidateAccount on Create
			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			for _, existingModel := range item.existingModels {
				srv.Core.CloudAccounts.Create(existingModel)
			}

			list := new(model.CloudAccountList)
			err := sg.CloudAccounts.List(list)

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

func TestCloudAccountsCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("CloudAccounts Create works correctly", t, func() {

		table := []struct {
			// Input
			model *model.CloudAccount
			// Mocks
			mockValidateAccountError error
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				model: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				mockValidateAccountError: nil,
				err: nil,
			},

			// Invalid provider
			{
				model: &model.CloudAccount{
					Name:        "test",
					Provider:    "nocloud",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				mockValidateAccountError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: Provider: regular expression mismatch"},
			},

			// No name
			{
				model: &model.CloudAccount{
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				mockValidateAccountError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: Name: zero value"},
			},

			// No credentials
			{
				model: &model.CloudAccount{
					Name:     "test",
					Provider: "aws",
				},
				mockValidateAccountError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: Credentials: zero value"},
			},

			// On Provider ValidateAccount error
			{
				model: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				mockValidateAccountError: errors.New("creds aren't working"),
				err: &model.Error{Status: 422, Message: "Validation failed: creds aren't working"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					ValidateAccountFn: func(m *model.CloudAccount) error {
						return item.mockValidateAccountError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			err := sg.CloudAccounts.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestCloudAccountsGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("CloudAccounts Get works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.CloudAccount
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			// For ValidateAccount on Create
			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.CloudAccounts.Create(item.existingModel)

			err := sg.CloudAccounts.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestCloudAccountsUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("CloudAccounts Update works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.CloudAccount
			modelUpdate   *model.CloudAccount
			// Expectations
			err *model.Error
		}{
			// Can't update Name
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				modelUpdate: &model.CloudAccount{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update Provider
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				modelUpdate: &model.CloudAccount{
					Provider: "do",
				},
				err: &model.Error{Status: 422, Message: "Provider cannot be changed"},
			},

			// Can't update Credentials
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				modelUpdate: &model.CloudAccount{
					Credentials: map[string]string{"new": "credz"},
				},
				err: &model.Error{Status: 422, Message: "Credentials cannot be changed"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			// For Create
			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.CloudAccounts.Create(item.existingModel)

			err := sg.CloudAccounts.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestCloudAccountsDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("CloudAccounts Delete works correctly", t, func() {

		table := []struct {
			// Input
			existingModel    *model.CloudAccount
			hasExistingKubes []*model.Kube
			// Expectations
			err *model.Error
		}{
			// Successful example
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				hasExistingKubes: nil,
				err:              nil,
			},

			// Can't delete if there's Kubes (cuz then we couldn't tear down)
			{
				existingModel: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"access_key": "blah", "secret_key": "bleh"},
				},
				hasExistingKubes: []*model.Kube{
					{
						CloudAccountName: "test",
						Name:             "testkube",
						MasterNodeSize:   "t2.micro",
						NodeSizes:        []string{"t2.micro"},
						AWSConfig: &model.AWSKubeConfig{
							Region:           "us-east-1",
							AvailabilityZone: "us-east-1a",
						},
					},
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Cannot delete CloudAccount that has active Kubes"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			// For Create
			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			srv.Core.CloudAccounts.Create(item.existingModel)

			for _, existingKube := range item.hasExistingKubes {
				srv.Core.Kubes.Create(existingKube)
			}

			err := sg.CloudAccounts.Delete(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}
