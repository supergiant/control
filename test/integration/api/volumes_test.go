package api

import (
	"errors"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVolumesList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Volumes List works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModels     []*model.Volume
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
				existingModels: []*model.Volume{
					{
						KubeName: "test",
						Name:     "my-volume",
						Type:     "gp2",
						Size:     10,
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
				srv.Core.Volumes.Create(existingModel)
			}

			list := new(model.VolumeList)
			err := sg.Volumes.List(list)

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

func TestVolumesCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Volumes Create works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			model              *model.Volume
			// Mocks
			mockCreateVolumeError error
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
				model: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				mockCreateVolumeError: nil,
				err: nil,
			},

			// No Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Volume{
					Name: "my-volume",
					Type: "gp2",
					Size: 10,
				},
				mockCreateVolumeError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},

			// Invalid Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Volume{
					KubeName: "non-existent",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				mockCreateVolumeError: nil,
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'KubeName' on Volume"},
			},

			// No Size
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
				model: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
				},
				mockCreateVolumeError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: Size: zero value"},
			},

			// On Provider CreateVolume error
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
				model: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				mockCreateVolumeError: errors.New("error creating Volume"),
				err: &model.Error{Status: 500, Message: "error creating Volume"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateVolumeFn: func(m *model.Volume, _ *core.Action) error {
						return item.mockCreateVolumeError
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

			err := sg.Volumes.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestVolumesGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Volumes Get works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Volume
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
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

			srv.Core.Volumes.Create(item.existingModel)

			err := sg.Volumes.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestVolumesUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Volumes Update works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Volume
			modelUpdate        *model.Volume
			// Expectations
			resizeTriggered bool
			err             *model.Error
		}{
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				modelUpdate: &model.Volume{
					KubeName: "new-name",
				},
				err: &model.Error{Status: 422, Message: "KubeName cannot be changed"},
			},

			// Can update Size, and it triggers Resize
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				modelUpdate: &model.Volume{
					Size: 20,
				},
				resizeTriggered: true,
				err:             nil,
			},

			// Can't update Type
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				modelUpdate: &model.Volume{
					Type: "new-type",
				},
				err: &model.Error{Status: 422, Message: "Type cannot be changed"},
			},

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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				modelUpdate: &model.Volume{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},
		}

		for _, item := range table {

			var resizeTriggered bool

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					ResizeVolumeFn: func(_ *model.Volume, _ *core.Action) error {
						resizeTriggered = true
						return nil
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

			srv.Core.Volumes.Create(item.existingModel)

			err := sg.Volumes.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			So(resizeTriggered, ShouldEqual, item.resizeTriggered)
		}
	})
}

//------------------------------------------------------------------------------

func TestVolumesDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Volumes Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Volume
			// Mocks
			mockProviderDeleteVolumeError error
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				mockProviderDeleteVolumeError: nil,
				statusError:                   "",
			},

			// On Provider DeleteVolume error
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
				existingModel: &model.Volume{
					KubeName: "test",
					Name:     "my-volume",
					Type:     "gp2",
					Size:     10,
				},
				mockProviderDeleteVolumeError: errors.New("error deleting Volume"),
				statusError:                   "error deleting Volume",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteVolumeFn: func(_ *model.Volume) error {
						return item.mockProviderDeleteVolumeError
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

			srv.Core.Volumes.Create(item.existingModel)

			err := sg.Volumes.Delete(item.existingModel.ID, item.existingModel)

			So(err, ShouldBeNil)

			// NOTE this is async error, so it is not the error returned from Delete.
			// Should have an update by the time this Get completes
			sg.Volumes.Get(item.existingModel.ID, item.existingModel)

			So(item.existingModel.Status.Error, ShouldEqual, item.statusError)
		}
	})
}
