package api

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_core"
)

var kubeResourceID int64 = 24

//------------------------------------------------------------------------------

func TestKubeResourcesCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	var lastStartCalledOn *int64

	realCreateFn := srv.Core.KubeResources.Create

	srv.Core.KubeResources = &fake_core.KubeResources{
		CreateFn: realCreateFn,
		StartFn: func(id *int64, m *model.KubeResource) core.ActionInterface {
			return &core.Action{
				Status: &model.ActionStatus{
					Description: "starting",
				},
				Core:       srv.Core,
				Model:      m,
				ResourceID: m.GetUUID(),
				Fn: func(_ *core.Action) error {
					lastStartCalledOn = id
					return nil
				},
			}
		},
	}

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("KubeResources Create works correctly", t, func() {

		table := []struct {
			// Input
			kubeResource *model.KubeResource
			// Output
			startCalled bool
			err         *model.Error
		}{
			// A successful example
			{
				&model.KubeResource{
					KubeName:  kube.Name,
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Resource: newRawMessage(`{
						"spec": {
							"containers": [
								{
									"name": "jenkins",
									"image": "jenkins"
								}
							]
						}
					}`),
				},
				true,
				nil,
			},
			// No KubeName provided
			{
				&model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Resource: newRawMessage(`{
						"spec": {
							"containers": [
								{
									"name": "jenkins",
									"image": "jenkins"
								}
							]
						}
					}`),
				},
				false,
				&model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},
			// Kube does not exist
			{
				&model.KubeResource{
					KubeName:  "crab",
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Resource: newRawMessage(`{
						"spec": {
							"containers": [
								{
									"name": "jenkins",
									"image": "jenkins"
								}
							]
						}
					}`),
				},
				false,
				&model.Error{Status: 422, Message: "Parent does not exist, foreign key 'KubeName' on KubeResource"},
			},
		}

		for _, item := range table {
			lastStartCalledOn = nil

			err := sg.KubeResources.Create(item.kubeResource)
			startCalled := lastStartCalledOn != nil && *lastStartCalledOn == *item.kubeResource.ID

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			So(startCalled, ShouldEqual, item.startCalled)
		}
	})
}

//------------------------------------------------------------------------------

func TestKubeResourcesStart(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("KubeResources Start uses correct Provisioner, and ensures existence of correct Namespace", t, func() {
		table := []struct {
			// Input
			kubeResource *model.KubeResource
			// Mocks
			mockStartTimeout bool
			// Expectations
			namespaceEnsured  string
			provisionerCalled string
			errorReturned     error
		}{
			// Pod
			{
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{ID: &kubeResourceID},
					KubeName:  kube.Name,
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Resource:  newRawMessage(`{}`),
				},
				mockStartTimeout: false,
				namespaceEnsured: "test",
				// provisionerCalled: "pod",
				errorReturned: nil,
			},
			// Service
			{
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{ID: &kubeResourceID},
					KubeName:  kube.Name,
					Namespace: "beep",
					Name:      "test",
					Kind:      "Service",
					Resource:  newRawMessage(`{}`),
				},
				mockStartTimeout: false,
				namespaceEnsured: "beep",
				// provisionerCalled: "service",
				errorReturned: nil,
			},
			// Anything else
			{
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{ID: &kubeResourceID},
					KubeName:  kube.Name,
					Namespace: "test",
					Name:      "test",
					Kind:      "LiterallyAnythingElse",
					Resource:  newRawMessage(`{}`),
				},
				mockStartTimeout: false,
				namespaceEnsured: "test",
				// provisionerCalled: "default",
				errorReturned: nil,
			},
			// Reports error on StartTimeout
			{
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{ID: &kubeResourceID},
					KubeName:  kube.Name,
					Namespace: "foo",
					Name:      "test",
					Kind:      "Service",
					Resource:  newRawMessage(`{}`),
				},
				mockStartTimeout: true,
				namespaceEnsured: "foo",
				// provisionerCalled: "service",
				errorReturned: errors.New("Service 'test' in Namespace 'foo' to start: context deadline exceeded"),
			},
		}

		for _, item := range table {

			var namespaceEnsured string
			// var provisionerCalled string

			srv.Core.KubeResourceStartTimeout = time.Nanosecond

			srv.Core.K8S = func(_ *model.Kube) kubernetes.ClientInterface {
				return &fake_core.KubernetesClient{
					EnsureNamespaceFn: func(name string) error {
						namespaceEnsured = name
						return nil
					},
				}
			}

			isRunningFn := func(_ *model.KubeResource) (bool, error) {
				return !item.mockStartTimeout, nil
			}

			srv.Core.DefaultProvisioner = &fake_core.Provisioner{
				ProvisionFn: func(_ *model.KubeResource) error {
					// provisionerCalled = "default"
					return nil
				},
				IsRunningFn: isRunningFn,
			}
			// srv.Core.PodProvisioner = &fake_core.Provisioner{
			// 	ProvisionFn: func(_ *model.KubeResource) error {
			// 		// provisionerCalled = "pod"
			// 		return nil
			// 	},
			// 	IsRunningFn: isRunningFn,
			// }
			// srv.Core.ServiceProvisioner = &fake_core.Provisioner{
			// 	ProvisionFn: func(_ *model.KubeResource) error {
			// 		// provisionerCalled = "service"
			// 		return nil
			// 	},
			// 	IsRunningFn: isRunningFn,
			// }

			sg.KubeResources.Create(item.kubeResource)
			sg.KubeResources.Start(item.kubeResource.ID, item.kubeResource)

			if item.errorReturned != nil {
				// Reload to get new status
				time.Sleep(1 * time.Second)
				sg.KubeResources.Get(item.kubeResource.ID, item.kubeResource)
				So(item.kubeResource.Status.Error, ShouldResemble, item.errorReturned.Error())
			}
			So(namespaceEnsured, ShouldEqual, item.namespaceEnsured)
			// So(provisionerCalled, ShouldEqual, item.provisionerCalled)

			// Cleanup
			srv.Core.DB.Delete(item.kubeResource)
		}
	})
}
