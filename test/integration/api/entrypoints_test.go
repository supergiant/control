package api

import (
	"errors"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEntrypointsList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Entrypoints List works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModels     []*model.Entrypoint
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModels: []*model.Entrypoint{
					{
						KubeName: "test",
						Name:     "my-entrypoint",
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.NewAPIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			for _, existingModel := range item.existingModels {
				srv.Core.Entrypoints.Create(existingModel)
			}

			list := new(model.EntrypointList)
			err := sg.Entrypoints.List(list)

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

func TestEntrypointsCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Entrypoints Create works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			model              *model.Entrypoint
			// Mocks
			mockCreateEntrypointError error
			// Expectations
			err         *model.Error
			statusError string
		}{
			// A successful example
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				model: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				mockCreateEntrypointError: nil,
				err: nil,
			},

			// No Name
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				model: &model.Entrypoint{
					KubeName: "test",
				},
				mockCreateEntrypointError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: Name: zero value"},
			},

			// No Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Entrypoint{
					Name: "my-entrypoint",
				},
				mockCreateEntrypointError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},

			// Invalid Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Entrypoint{
					KubeName: "non-existent",
					Name:     "my-entrypoint",
				},
				mockCreateEntrypointError: nil,
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'KubeName' on Entrypoint"},
			},

			// On Provider CreateEntrypoint error
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				model: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				mockCreateEntrypointError: errors.New("error creating entrypoint"),
				err:         nil,
				statusError: "error creating entrypoint",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateEntrypointFn: func(m *model.Entrypoint, _ *core.Action) error {
						return item.mockCreateEntrypointError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.NewAPIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			err := sg.Entrypoints.Create(item.model)

			if item.err == nil {
				// NOTE The Provider part of Create is Async.
				// We can only call this if the non-Async err is nil (meaning the Action started).
				sg.Entrypoints.Get(item.model.ID, item.model)
				So(item.model.Status.Error, ShouldEqual, item.statusError)

				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointsGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Entrypoints Get works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Entrypoint
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModel: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
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
			sg := srv.Core.NewAPIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Entrypoints.Create(item.existingModel)

			err := sg.Entrypoints.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointsUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Entrypoints Update works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Entrypoint
			modelUpdate        *model.Entrypoint
			// Expectations
			err *model.Error
		}{
			// Can't update Name
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModel: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				modelUpdate: &model.Entrypoint{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update KubeName
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModel: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				modelUpdate: &model.Entrypoint{
					KubeName: "new-name",
				},
				err: &model.Error{Status: 422, Message: "KubeName cannot be changed"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.NewAPIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Entrypoints.Create(item.existingModel)

			err := sg.Entrypoints.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointsDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Entrypoints Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount          *model.CloudAccount
			parentKube                  *model.Kube
			existingModel               *model.Entrypoint
			existingEntrypointListeners []*model.EntrypointListener
			// Mocks
			mockProviderDeleteEntrypointError error
			// Expectations
			entrypointListenerNamesDeleted []string
			statusError                    string
		}{
			// Successful example
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModel: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				existingEntrypointListeners: []*model.EntrypointListener{
					{
						EntrypointName: "my-entrypoint",
						Name:           "test",
						EntrypointPort: 100,
						NodePort:       101,
					},
				},
				mockProviderDeleteEntrypointError: nil,
				entrypointListenerNamesDeleted:    []string{"test"},
				statusError:                       "",
			},

			// On Provider DeleteEntrypoint error
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				parentKube: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingModel: &model.Entrypoint{
					KubeName: "test",
					Name:     "my-entrypoint",
				},
				existingEntrypointListeners:       nil,
				mockProviderDeleteEntrypointError: errors.New("error deleting Entrypoint"),
				entrypointListenerNamesDeleted:    nil,
				statusError:                       "error deleting Entrypoint",
			},
		}

		for _, item := range table {

			var entrypointListenerNamesDeleted []string

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteEntrypointFn: func(_ *model.Entrypoint) error {
						return item.mockProviderDeleteEntrypointError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.NewAPIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Entrypoints.Create(item.existingModel)

			for _, existingEntrypointListener := range item.existingEntrypointListeners {
				srv.Core.EntrypointListeners.Create(existingEntrypointListener)
			}

			err := sg.Entrypoints.Delete(item.existingModel.ID, item.existingModel)

			So(err, ShouldBeNil)

			// NOTE this is async error, so it is not the error returned from Delete.
			// Should have an update by the time this Get completes
			sg.Entrypoints.Get(item.existingModel.ID, item.existingModel)

			So(item.existingModel.Status.Error, ShouldEqual, item.statusError)

			// NOTE Since EntrypointListeners are deleted directly, we have to detect
			// the names deleted in an inverse fashion.
			for _, existingEntrypointListener := range item.existingEntrypointListeners {
				entrypointListenerNamesDeleted = append(entrypointListenerNamesDeleted, existingEntrypointListener.Name)
			}
			var remainingEntrypointListeners []*model.EntrypointListener
			srv.Core.DB.Find(&remainingEntrypointListeners)
			for _, remainingEntrypointListener := range remainingEntrypointListeners {
				for i, nameDeleted := range entrypointListenerNamesDeleted {
					if remainingEntrypointListener.Name == nameDeleted {
						// Delete
						entrypointListenerNamesDeleted = append(entrypointListenerNamesDeleted[:i], entrypointListenerNamesDeleted[i+1:]...)
					}
				}
			}

			So(entrypointListenerNamesDeleted, ShouldResemble, item.entrypointListenerNamesDeleted)
		}
	})
}
