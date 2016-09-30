package api

import (
	"errors"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEntrypointListenersList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("EntrypointListeners List works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			parentEntrypoint   *model.Entrypoint
			existingModels     []*model.EntrypointListener
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModels: []*model.EntrypointListener{
					{
						EntrypointName: "test",
						Name:           "port-test",
						EntrypointPort: 80,
						NodePort:       30303,
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
			if item.parentEntrypoint != nil {
				srv.Core.Entrypoints.Create(item.parentEntrypoint)
			}

			for _, existingModel := range item.existingModels {
				srv.Core.EntrypointListeners.Create(existingModel)
			}

			list := new(model.EntrypointListenerList)
			err := sg.EntrypointListeners.List(list)

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

func TestEntrypointListenersCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("EntrypointListeners Create works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			parentEntrypoint   *model.Entrypoint
			model              *model.EntrypointListener
			// Mocks
			mockCreateEntrypointListenerError error
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				model: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockCreateEntrypointListenerError: nil,
				err: nil,
			},

			// No Entrypoint
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				parentEntrypoint:   nil,
				model: &model.EntrypointListener{
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockCreateEntrypointListenerError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: EntrypointName: zero value"},
			},

			// Invalid Entrypoint
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				parentEntrypoint:   nil,
				model: &model.EntrypointListener{
					EntrypointName: "non-existent",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockCreateEntrypointListenerError: nil,
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'EntrypointName' on EntrypointListener"},
			},

			// On Provider CreateEntrypointListener error
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				model: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockCreateEntrypointListenerError: errors.New("error creating entrypoint"),
				err: &model.Error{Status: 500, Message: "error creating entrypoint"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateEntrypointListenerFn: func(m *model.EntrypointListener) error {
						return item.mockCreateEntrypointListenerError
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
			if item.parentEntrypoint != nil {
				srv.Core.Entrypoints.Create(item.parentEntrypoint)
			}

			err := sg.EntrypointListeners.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("EntrypointListeners Get works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			parentEntrypoint   *model.Entrypoint
			existingModel      *model.EntrypointListener
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
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
			if item.parentEntrypoint != nil {
				srv.Core.Entrypoints.Create(item.parentEntrypoint)
			}

			srv.Core.EntrypointListeners.Create(item.existingModel)

			err := sg.EntrypointListeners.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("EntrypointListeners Update works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			parentEntrypoint   *model.Entrypoint
			existingModel      *model.EntrypointListener
			modelUpdate        *model.EntrypointListener
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update EntrypointName
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					EntrypointName: "new-name",
				},
				err: &model.Error{Status: 422, Message: "EntrypointName cannot be changed"},
			},

			// Can't update EntrypointPort
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					EntrypointPort: 66,
				},
				err: &model.Error{Status: 422, Message: "EntrypointPort cannot be changed"},
			},

			// Can't update EntrypointProtocol
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					EntrypointProtocol: "HTTPS",
				},
				err: &model.Error{Status: 422, Message: "EntrypointProtocol cannot be changed"},
			},

			// Can't update NodePort
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					NodePort: 4444,
				},
				err: &model.Error{Status: 422, Message: "NodePort cannot be changed"},
			},

			// Can't update NodeProtocol
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				modelUpdate: &model.EntrypointListener{
					NodeProtocol: "UDP",
				},
				err: &model.Error{Status: 422, Message: "NodeProtocol cannot be changed"},
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
			if item.parentEntrypoint != nil {
				srv.Core.Entrypoints.Create(item.parentEntrypoint)
			}

			srv.Core.EntrypointListeners.Create(item.existingModel)

			err := sg.EntrypointListeners.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestEntrypointListenersDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("EntrypointListeners Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			parentEntrypoint   *model.Entrypoint
			existingModel      *model.EntrypointListener
			// Mocks
			mockDeleteEntrypointListenerError error
			// Expectations
			statusError string
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockDeleteEntrypointListenerError: nil,
				statusError:                       "",
			},

			// On Provider DeleteEntrypointListener error
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
				parentEntrypoint: &model.Entrypoint{
					KubeName: "test",
					Name:     "test",
				},
				existingModel: &model.EntrypointListener{
					EntrypointName: "test",
					Name:           "port-test",
					EntrypointPort: 80,
					NodePort:       30303,
				},
				mockDeleteEntrypointListenerError: errors.New("error deleting EntrypointListener"),
				statusError:                       "error deleting EntrypointListener",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteEntrypointListenerFn: func(_ *model.EntrypointListener) error {
						return item.mockDeleteEntrypointListenerError
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
			if item.parentEntrypoint != nil {
				srv.Core.Entrypoints.Create(item.parentEntrypoint)
			}

			srv.Core.EntrypointListeners.Create(item.existingModel)

			sg.EntrypointListeners.Delete(item.existingModel.ID, item.existingModel)

			// NOTE this is async error, so it is not the error returned from Delete.
			// Should have an update by the time this Get completes
			sg.EntrypointListeners.Get(item.existingModel.ID, item.existingModel)

			So(item.existingModel.Status.Error, ShouldEqual, item.statusError)
		}
	})
}
