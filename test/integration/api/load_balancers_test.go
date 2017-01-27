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

var loadBalancerID int64 = 31

func TestLoadBalancersList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("LoadBalancers List works correctly", t, func() {

		table := []struct {
			// Input
			existingModels []*model.LoadBalancer
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModels: []*model.LoadBalancer{
					&model.LoadBalancer{
						KubeName:  kube.Name,
						Name:      "test",
						Namespace: "default",
					},
				},
				err: nil,
			},
		}

		for _, item := range table {

			for _, existingModel := range item.existingModels {
				srv.Core.LoadBalancers.Create(existingModel)
			}

			list := new(model.LoadBalancerList)
			err := sg.LoadBalancers.List(list)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			So(len(list.Items), ShouldEqual, len(item.existingModels))

			// NOTE we have to clean up LoadBalancers manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.LoadBalancer{})
		}
	})
}

//------------------------------------------------------------------------------

func TestLoadBalancersCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("LoadBalancers Create works correctly", t, func() {

		table := []struct {
			// Input
			model *model.LoadBalancer
			// Mocks
			mockProviderCreateLoadBalancerError error
			// Expectations
			err      *model.Error
			asyncErr string // because Create has a sync and async phase
		}{
			// A successful example
			{
				model: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				err:      nil,
				asyncErr: "",
			},

			// On Provider CreateLoadBalancer error
			{
				model: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				mockProviderCreateLoadBalancerError: errors.New("butts"),
				err:      nil,
				asyncErr: "butts",
			},

			// No KubeName
			{
				model: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					Name:      "test",
					Namespace: "default",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: KubeName: zero value"},
			},

			// No Name
			{
				model: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Namespace: "default",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Name: zero value"},
			},

			// No Namespace
			{
				model: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
				},
				err: &model.Error{Status: 422, Message: "Validation failed: Namespace: zero value"},
			},
		}

		for _, item := range table {

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					CreateLoadBalancerFn: func(_ *model.LoadBalancer, _ *core.Action) error {
						return item.mockProviderCreateLoadBalancerError
					},
				}
			}

			err := sg.LoadBalancers.Create(item.model)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// Reload to get new status
			// We use fresh to get rid of status (which was set by Create, and is nil at Get)
			freshModel := new(model.LoadBalancer)
			sg.LoadBalancers.Get(item.model.ID, freshModel)

			if item.asyncErr == "" {
				So(freshModel.Status, ShouldBeNil)
			} else {
				So(freshModel.Status.Error, ShouldEqual, item.asyncErr)
			}

			// NOTE we have to clean up LoadBalancers manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.LoadBalancer{})
		}
	})
}

//------------------------------------------------------------------------------

func TestLoadBalancersGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("LoadBalancers Get works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.LoadBalancer
			// Expectations
			err *model.Error
		}{
			// A successful example
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				err: nil,
			},
		}

		for _, item := range table {

			srv.Core.LoadBalancers.Create(item.existingModel)

			err := sg.LoadBalancers.Get(item.existingModel.ID, item.existingModel)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// NOTE we have to clean up LoadBalancers manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.LoadBalancer{})
		}
	})
}

//------------------------------------------------------------------------------

func TestLoadBalancersUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("LoadBalancers Update works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.LoadBalancer
			modelUpdate   *model.LoadBalancer
			// Mocks
			mockProviderUpdateLoadBalancerError error
			// Expectations
			err      *model.Error
			asyncErr string // because Update has a sync and async phase
		}{
			// Successful example
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				modelUpdate: &model.LoadBalancer{
					Ports: map[int]int{80: 8080},
				},
				err:      nil,
				asyncErr: "",
			},

			// On Provider UpdateLoadBalancer error
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				modelUpdate: &model.LoadBalancer{
					Ports: map[int]int{80: 8080},
				},
				mockProviderUpdateLoadBalancerError: errors.New("something unexpected"),
				err:      nil,
				asyncErr: "something unexpected",
			},

			// Can't update KubeName
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				modelUpdate: &model.LoadBalancer{
					KubeName: "butt",
				},
				err: &model.Error{Status: 422, Message: "KubeName cannot be changed"},
			},

			// Can't update Name
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				modelUpdate: &model.LoadBalancer{
					Name: "butt",
				},
				err: &model.Error{Status: 422, Message: "Name cannot be changed"},
			},

			// Can't update Namespace
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				modelUpdate: &model.LoadBalancer{
					Namespace: "butt",
				},
				err: &model.Error{Status: 422, Message: "Namespace cannot be changed"},
			},
		}

		for _, item := range table {

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					UpdateLoadBalancerFn: func(_ *model.LoadBalancer, _ *core.Action) error {
						return item.mockProviderUpdateLoadBalancerError
					},
				}
			}

			srv.Core.LoadBalancers.Create(item.existingModel)

			err := sg.LoadBalancers.Update(item.existingModel.ID, item.modelUpdate)

			if item.err == nil {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldResemble, item.err)
			}

			// Reload to get new status
			// We use fresh to get rid of status (which was set by Update, and is nil at Get)
			freshModel := new(model.LoadBalancer)
			sg.LoadBalancers.Get(item.existingModel.ID, freshModel)

			if item.asyncErr == "" {
				So(freshModel.Status, ShouldBeNil)
			} else {
				So(freshModel.Status.Error, ShouldEqual, item.asyncErr)
			}

			// NOTE we have to clean up LoadBalancers manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.LoadBalancer{})
		}
	})
}

//------------------------------------------------------------------------------

func TestLoadBalancersDelete(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	requestor := createAdmin(srv.Core)
	sg := srv.Core.APIClient("token", requestor.APIToken)

	srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
		return new(fake_core.Provider)
	}
	kube := createKube(sg)

	Convey("LoadBalancers Delete works correctly", t, func() {

		table := []struct {
			// Input
			existingModel *model.LoadBalancer
			// Mocks
			mockProviderDeleteLoadBalancerError error
			// Expectations
			asyncErr string
		}{
			// Successful example
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				asyncErr: "",
			},

			// On Provider DeleteLoadBalancer error
			{
				existingModel: &model.LoadBalancer{
					BaseModel: model.BaseModel{ID: &loadBalancerID},
					KubeName:  kube.Name,
					Name:      "test",
					Namespace: "default",
				},
				mockProviderDeleteLoadBalancerError: errors.New("something unexpected"),
				asyncErr: "something unexpected",
			},
		}

		for _, item := range table {

			srv.Core.AWSProvider = func(_ map[string]string) core.Provider {
				return &fake_core.Provider{
					DeleteLoadBalancerFn: func(_ *model.LoadBalancer, _ *core.Action) error {
						return item.mockProviderDeleteLoadBalancerError
					},
				}
			}

			srv.Core.LoadBalancers.Create(item.existingModel)

			// NOTE no need to test this error here
			_ = sg.LoadBalancers.Delete(item.existingModel.ID, item.existingModel)

			time.Sleep(10 * time.Millisecond)

			// Reload to get new status
			// We use fresh to get rid of status (which was set by Delete, and is nil at Get)
			freshModel := new(model.LoadBalancer)
			getErr := sg.LoadBalancers.Get(item.existingModel.ID, freshModel)

			if item.asyncErr == "" {
				So(getErr, ShouldNotBeNil)
			} else {
				So(freshModel.Status.Error, ShouldEqual, item.asyncErr)
			}

			// NOTE we have to clean up LoadBalancers manually since we do not wipe DB each time
			srv.Core.DB.Delete(&model.LoadBalancer{})
		}
	})
}
