package api

import (
	"errors"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKubesList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kubes List works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			existingModels     []*model.Kube
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
				existingModels: []*model.Kube{
					{
						CloudAccountName: "test",
						Name:             "test",
						MasterNodeSize:   "t2.micro",
						NodeSizes:        []string{"t2.micro"},
						AWSConfig: &model.AWSKubeConfig{
							Region:           "us-east-1",
							AvailabilityZone: "us-east-1a",
						},
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

			for _, existingModel := range item.existingModels {
				srv.Core.Kubes.Create(existingModel)
			}

			list := new(model.KubeList)
			err := sg.Kubes.List(list)

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

func TestKubesCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kubes Create works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			model              *model.Kube
			// Mocks
			mockProviderCreateKubeError error
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
				model: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: nil,
				err: nil,
			},

			// No CloudAccount
			{
				parentCloudAccount: nil,
				model: &model.Kube{
					Name:           "test",
					MasterNodeSize: "t2.micro",
					NodeSizes:      []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: nil,
				err: &model.Error{Status: 422, Message: "Validation failed: CloudAccountName: zero value"},
			},

			// Invalid CloudAccount
			{
				parentCloudAccount: nil,
				model: &model.Kube{
					CloudAccountName: "non-existent",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: nil,
				err: &model.Error{Status: 422, Message: "Parent does not exist, foreign key 'CloudAccountName' on Kube"},
			},

			// On Provider CreateKube error
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				model: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: errors.New("error creating kube"),
				err:         nil,
				statusError: "error creating kube",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateKubeFn: func(m *model.Kube, _ *core.Action) error {
						return item.mockProviderCreateKubeError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}

			err := sg.Kubes.Create(item.model)

			if item.err == nil {
				// NOTE The Provider part of Create is Async.
				// We can only call this if the non-Async err is nil (meaning the Action started).
				sg.Kubes.Get(item.model.ID, item.model)
				So(item.model.Status.Error, ShouldEqual, item.statusError)

				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------
//
// TODO this repeats some of the Create test above. The Create test should
// really check validations, and then assert Provision was called.
//
//
//
func TestKubesProvision(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kubes Provision works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			model              *model.Kube
			// Mocks
			mockProviderCreateKubeError error
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
				model: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					Username:         "username",
					Password:         "password",
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: nil,
				err: nil,
			},

			// On Provider CreateKube error
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				model: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					Username:         "username",
					Password:         "password",
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				mockProviderCreateKubeError: errors.New("error creating kube"),
				err:         nil,
				statusError: "error creating kube",
			},
		}

		for _, item := range table {

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateKubeFn: func(m *model.Kube, _ *core.Action) error {
						return item.mockProviderCreateKubeError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}

			// Create model directly in DB first
			srv.Core.DB.Create(item.model)

			err := sg.Kubes.Provision(item.model.ID, item.model)

			if item.err == nil {
				// NOTE The Provider part of Create is Async.
				// We can only call this if the non-Async err is nil (meaning the Action started).
				sg.Kubes.Get(item.model.ID, item.model)
				So(item.model.Status.Error, ShouldEqual, item.statusError)

				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestKubesGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kubes Get works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			existingModel      *model.Kube
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
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
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

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}

			srv.Core.Kubes.Create(item.existingModel)

			err := sg.Kubes.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestKubesUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kubes Update works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount *model.CloudAccount
			existingModel      *model.Kube
			modelUpdate        *model.Kube
			// Expectations
			err *model.Error
		}{
			// Can update NodeSizes
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					NodeSizes: []string{"t2.small"},
				},
				err: nil,
			},

			// Can't update Name
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					Name: "new-name",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update CloudAccountName
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					CloudAccountName: "new-name",
				},
				err: &model.Error{Status: 422, Message: "CloudAccountName cannot be changed"},
			},

			// Can't update MasterNodeSize
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					MasterNodeSize: "newsize",
				},
				err: &model.Error{Status: 422, Message: "MasterNodeSize cannot be changed"},
			},

			// Can't update Username
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					Username: "Username",
				},
				err: &model.Error{Status: 422, Message: "Username cannot be changed"},
			},

			// Can't update Password
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					Password: "Password",
				},
				err: &model.Error{Status: 422, Message: "Password cannot be changed"},
			},

			// Can't update HeapsterVersion
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					HeapsterVersion: "v2.0.0",
				},
				err: &model.Error{Status: 422, Message: "HeapsterVersion cannot be changed"},
			},

			// Can't update HeapsterMetricResolution
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					HeapsterMetricResolution: "10s",
				},
				err: &model.Error{Status: 422, Message: "HeapsterMetricResolution cannot be changed"},
			},

			// Can't update AWSConfig
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					AWSConfig: &model.AWSKubeConfig{
						Region: "us-east-7",
					},
				},
				err: &model.Error{Status: 422, Message: "AWSConfig cannot be changed"},
			},

			// Can't update DigitalOceanConfig
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				modelUpdate: &model.Kube{
					DigitalOceanConfig: &model.DOKubeConfig{
						Region: "nyc1",
					},
				},
				err: &model.Error{Status: 422, Message: "DigitalOceanConfig cannot be changed"},
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

			srv.Core.Kubes.Create(item.existingModel)

			err := sg.Kubes.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}
		}
	})
}

//------------------------------------------------------------------------------

func TestKubeDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Kube Delete works correctly", t, func() {

		table := []struct {
			// Input
			parentCloudAccount    *model.CloudAccount
			existingModel         *model.Kube
			existingNodes         []*model.Node
			existingKubeResources []*model.KubeResource
			// Mocks
			mockProviderDeleteKubeError error
			// Expectations
			entrypointListenerNamesDeleted []string
			entrypointNamesDeleted         []string
			nodeNamesDeleted               []string
			kubeResourceNamesDeleted       []string
			volumeNamesDeleted             []string
			statusError                    string
		}{
			// Successful example
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingNodes: []*model.Node{
					{
						KubeName: "test",
						Name:     "test",
						Size:     "t2.micro",
					},
				},
				existingKubeResources: []*model.KubeResource{
					{
						KubeName:  "test",
						Namespace: "test",
						Kind:      "Thingy",
						Name:      "test",
						Resource:  newRawMessage(`{}`),
					},
				},
				mockProviderDeleteKubeError:    nil,
				entrypointNamesDeleted:         []string{"test"},
				entrypointListenerNamesDeleted: []string{"test"},
				nodeNamesDeleted:               []string{"test"},
				kubeResourceNamesDeleted:       []string{"test"},
				volumeNamesDeleted:             []string{"test"},
				statusError:                    "",
			},

			// On Provider DeleteKube error
			{
				parentCloudAccount: &model.CloudAccount{
					Name:        "test",
					Provider:    "aws",
					Credentials: map[string]string{"test": "test"},
				},
				existingModel: &model.Kube{
					CloudAccountName: "test",
					Name:             "test",
					MasterNodeSize:   "t2.micro",
					NodeSizes:        []string{"t2.micro"},
					AWSConfig: &model.AWSKubeConfig{
						Region:           "us-east-1",
						AvailabilityZone: "us-east-1a",
					},
				},
				existingNodes:                  nil,
				existingKubeResources:          nil,
				mockProviderDeleteKubeError:    errors.New("error deleting Kube"),
				entrypointNamesDeleted:         nil,
				entrypointListenerNamesDeleted: nil,
				nodeNamesDeleted:               nil,
				kubeResourceNamesDeleted:       nil,
				volumeNamesDeleted:             nil,
				statusError:                    "error deleting Kube",
			},
		}

		for _, item := range table {

			var nodeNamesDeleted []string
			var kubeResourceNamesDeleted []string

			wipeAndInitialize(srv.Core)

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteKubeFn: func(_ *model.Kube, _ *core.Action) error {
						return item.mockProviderDeleteKubeError
					},
				}
			}

			requestor := createAdmin(srv.Core)
			sg := srv.Core.APIClient("token", requestor.APIToken)

			if item.parentCloudAccount != nil {
				srv.Core.CloudAccounts.Create(item.parentCloudAccount)
			}

			srv.Core.Kubes.Create(item.existingModel)

			for _, existingNode := range item.existingNodes {
				srv.Core.Nodes.Create(existingNode)
			}
			for _, existingKubeResource := range item.existingKubeResources {
				srv.Core.KubeResources.Create(existingKubeResource)
			}

			err := sg.Kubes.Delete(item.existingModel.ID, item.existingModel)

			// TODO (cuz Async delete is too slow for following DB lookups)
			time.Sleep(time.Second)

			So(err, ShouldBeNil)

			// NOTE this is async error, so it is not the error returned from Delete.
			// Should have an update by the time this Get completes
			sg.Kubes.Get(item.existingModel.ID, item.existingModel)

			So(item.existingModel.Status.Error, ShouldEqual, item.statusError)

			// Find Nodes deleted
			for _, existingNode := range item.existingNodes {
				nodeNamesDeleted = append(nodeNamesDeleted, existingNode.Name)
			}
			var remainingNodes []*model.Node
			srv.Core.DB.Find(&remainingNodes)
			for _, remainingNode := range remainingNodes {
				for i, nameDeleted := range nodeNamesDeleted {
					if remainingNode.Name == nameDeleted {
						// Delete
						nodeNamesDeleted = append(nodeNamesDeleted[:i], nodeNamesDeleted[i+1:]...)
					}
				}
			}
			So(nodeNamesDeleted, ShouldResemble, item.nodeNamesDeleted)

			// Find KubeResourcesDeleted
			for _, existingKubeResource := range item.existingKubeResources {
				kubeResourceNamesDeleted = append(kubeResourceNamesDeleted, existingKubeResource.Name)
			}
			var remainingKubeResources []*model.KubeResource
			srv.Core.DB.Find(&remainingKubeResources)
			for _, remainingKubeResource := range remainingKubeResources {
				for i, nameDeleted := range kubeResourceNamesDeleted {
					if remainingKubeResource.Name == nameDeleted {
						// Delete
						kubeResourceNamesDeleted = append(kubeResourceNamesDeleted[:i], kubeResourceNamesDeleted[i+1:]...)
					}
				}
			}
			So(kubeResourceNamesDeleted, ShouldResemble, item.kubeResourceNamesDeleted)
		}
	})
}
