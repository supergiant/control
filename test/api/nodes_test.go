package api

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"
)

func TestNodesList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Nodes List works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModels     []*model.Node
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
				existingModels: []*model.Node{
					{
						KubeName: "test",
						Name:     "testnode.host",
						Size:     "m4.large",
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
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			for _, existingModel := range item.existingModels {
				srv.Core.Nodes.Create(existingModel)
			}

			list := new(model.NodeList)
			err := sg.Nodes.List(list)

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

func TestNodesCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Nodes Create works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			model              *model.Node
			// Mocks
			mockCreateNodeError error
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
				model: &model.Node{
					KubeName: "test",
					Size:     "m4.large",
				},
				mockCreateNodeError: nil,
				err:                 nil,
			},

			// No Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Node{
					Size: "t2.micro",
				},
				mockCreateNodeError: nil,
				err:                 &model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},

			// Invalid Kube
			{
				parentCloudAccount: nil,
				parentKube:         nil,
				model: &model.Node{
					KubeName: "non-existent",
					Size:     "t2.micro",
				},
				mockCreateNodeError: nil,
				err:                 &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'KubeName' on Node"},
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
				model: &model.Node{
					KubeName: "test",
				},
				mockCreateNodeError: nil,
				err:                 &model.Error{Status: 422, Message: "Validation failed: Size: zero value"},
			},

			// On Provider CreateNode error
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
				model: &model.Node{
					KubeName: "test",
					Size:     "t2.micro",
				},
				mockCreateNodeError: errors.New("error creating Node"),
				err:                 nil,
				statusError:         "error creating Node",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateNodeFn: func(m *model.Node, _ *core.Action) error {
						return item.mockCreateNodeError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			err := sg.Nodes.Create(item.model)

			if item.err == nil {
				// NOTE The Provider part of Create is Async.
				// We can only call this if the non-Async err is nil (meaning the Action started).
				sg.Nodes.Get(item.model.ID, item.model)
				So(item.model.Status.Error, ShouldEqual, item.statusError)

				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestNodesGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Nodes Get works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Node
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
				existingModel: &model.Node{
					KubeName: "test",
					Size:     "t2.micro",
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

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Nodes.Create(item.existingModel)

			err := sg.Nodes.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestNodesUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Nodes Update works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Node
			modelUpdate        *model.Node
			// Expectations
			err *model.Error
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
				existingModel: &model.Node{
					KubeName: "test",
					Size:     "t2.micro",
				},
				modelUpdate: &model.Node{
					KubeName: "new-name",
				},
				err: &model.Error{Status: 422, Message: "KubeName cannot be changed"},
			},

			// Can't update Size
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
				existingModel: &model.Node{
					KubeName: "test",
					Size:     "t2.micro",
				},
				modelUpdate: &model.Node{
					Size: "m4.10xlarge",
				},
				err: &model.Error{Status: 422, Message: "Size cannot be changed"},
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return new(fake_core.Provider)
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Nodes.Create(item.existingModel)

			err := sg.Nodes.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestNodesDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Nodes Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			parentKube         *model.Kube
			existingModel      *model.Node
			// Mocks
			mockProviderDeleteNodeError error
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
				existingModel: &model.Node{
					KubeName:   "test",
					Size:       "t2.micro",
					Name:       "test.host",
					ProviderID: "prov-ID",
				},
				mockProviderDeleteNodeError: nil,
				statusError:                 "",
			},

			// On Provider DeleteNode error
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
				existingModel: &model.Node{
					KubeName:   "test",
					Size:       "t2.micro",
					Name:       "test.host",
					ProviderID: "prov-ID",
				},
				mockProviderDeleteNodeError: errors.New("error deleting Node"),
				statusError:                 "error deleting Node",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteNodeFn: func(_ *model.Node, _ *core.Action) error {
						return item.mockProviderDeleteNodeError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}
			if item.parentKube != nil {
				srv.Core.Kubes.Create(item.parentKube)
			}

			srv.Core.Nodes.Create(item.existingModel)

			err := sg.Nodes.Delete(item.existingModel.ID, item.existingModel)

			So(err, ShouldBeNil)

			// NOTE this is async error, so it is not the error returned from Delete.
			// Should have an update by the time this Get completes
			sg.Nodes.Get(item.existingModel.ID, item.existingModel)

			So(item.existingModel.Status.Error, ShouldEqual, item.statusError)
		}
	})
}
