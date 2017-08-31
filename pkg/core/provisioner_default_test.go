package core_test

import (
	"encoding/json"
	// . "github.com/smartystreets/goconvey/convey"
)

func newRawMessage(str string) *json.RawMessage {
	rawmsg := json.RawMessage([]byte(str))
	return &rawmsg
}

// func TestDefaultProvisionerProvision(t *testing.T) {
// 	Convey("DefaultProvisioner Provision passses correct arguments to KubernetesClient CreateResource(), and passes error through", t, func() {
// 		table := []struct {
// 			// Input
// 			kubeResource *model.KubeResource
//
// 			// Mocks
// 			mockCreateResourceOut   json.RawMessage
// 			mockCreateResourceError error
//
// 			// Assertions
// 			kindPassed      string
// 			namespacePassed string
// 			objInPassed     map[string]interface{}
// 			outSaved        json.RawMessage
// 			errorReturned   error
// 		}{
// 			// A successful example
// 			//------------------------------------------------------------------------
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "test",
// 					Name:      "test",
// 					Kind:      "Service",
// 					Resource: newRawMessage(`{
// 						"spec": {
// 							"ports": [
// 								{
// 									"port": 80
// 								}
// 							]
// 						}
// 					}`),
// 				},
// 				// Mocks
// 				mockCreateResourceOut:   json.RawMessage([]byte(`{"just":"making sure output is captured correctly"}`)),
// 				mockCreateResourceError: nil,
// 				// Assertions
// 				kindPassed:      "Service",
// 				namespacePassed: "test",
// 				objInPassed: map[string]interface{}{
// 					"apiVersion": "v1",
// 					"kind":       "Service",
// 					"metadata": map[string]interface{}{
// 						"namespace": "test",
// 						"name":      "test",
// 					},
// 					"spec": map[string]interface{}{
// 						"ports": []interface{}{
// 							map[string]interface{}{
// 								"port": 80,
// 							},
// 						},
// 					},
// 				},
// 				outSaved:      json.RawMessage([]byte(`{"just":"making sure output is captured correctly"}`)),
// 				errorReturned: nil,
// 			},
// 			// A successful example with metadata provided
// 			//------------------------------------------------------------------------
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "foo",
// 					Name:      "bar",
// 					Kind:      "Service",
// 					Resource: newRawMessage(`{
// 						"apiVersion": "user-provided-version",
// 						"metadata": {
// 							"namespace": "erroneous user input we don't care about",
// 							"labels": {
// 								"cool": "label"
// 							}
// 						},
// 						"spec": {
// 							"ports": [
// 								{
// 									"port": 80
// 								}
// 							]
// 						}
// 					}`),
// 				},
// 				// Mocks
// 				mockCreateResourceOut:   json.RawMessage([]byte(`{"just":"making sure output is captured correctly (again)"}`)),
// 				mockCreateResourceError: nil,
// 				// Assertions
// 				kindPassed:      "Service",
// 				namespacePassed: "foo",
// 				objInPassed: map[string]interface{}{
// 					"apiVersion": "user-provided-version",
// 					"kind":       "Service",
// 					"metadata": map[string]interface{}{
// 						"namespace": "foo",
// 						"name":      "bar",
// 						"labels": map[string]string{
// 							"cool": "label",
// 						},
// 					},
// 					"spec": map[string]interface{}{
// 						"ports": []interface{}{
// 							map[string]interface{}{
// 								"port": 80,
// 							},
// 						},
// 					},
// 				},
// 				outSaved:      json.RawMessage([]byte(`{"just":"making sure output is captured correctly (again)"}`)),
// 				errorReturned: nil,
// 			},
// 			//------------------------------------------------------------------------
// 			// When there's an unexpected error from Kubernetes
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "beep",
// 					Name:      "borp",
// 					Kind:      "Secret",
// 					Resource: newRawMessage(`{
// 						"spec": {}
// 					}`),
// 				},
// 				// Mocks
// 				mockCreateResourceOut:   nil,
// 				mockCreateResourceError: errors.New("dunno what's happening"),
// 				// Assertions
// 				kindPassed:      "Secret",
// 				namespacePassed: "beep",
// 				objInPassed: map[string]interface{}{
// 					"apiVersion": "v1",
// 					"kind":       "Secret",
// 					"metadata": map[string]interface{}{
// 						"namespace": "beep",
// 						"name":      "borp",
// 					},
// 					"spec": map[string]interface{}{},
// 				},
// 				outSaved:      nil,
// 				errorReturned: errors.New("dunno what's happening"),
// 			},
// 		}
//
// 		for _, item := range table {
// 			var kindPassed string
// 			var namespacePassed string
// 			var objInPassed map[string]interface{}
//
// 			var outSaved json.RawMessage
//
// 			c := &core.Core{
// 				K8S: func(_ *model.Kube) kubernetes.ClientInterface {
// 					return &fake_core.KubernetesClient{
// 						CreateResourceFn: func(apiVersion string, kind string, namespace string, objIn interface{}, out interface{}) error {
// 							kindPassed = kind
// 							namespacePassed = namespace
// 							objInPassed = objIn.(map[string]interface{})
//
// 							reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.mockCreateResourceOut))
//
// 							return item.mockCreateResourceError
// 						},
// 					}
// 				},
// 				DB: &fake_core.DB{
// 					SaveFn: func(m model.Model) error {
// 						outSaved = *m.(*model.KubeResource).Resource
// 						return nil
// 					},
// 				},
// 			}
//
// 			provisioner := &core.DefaultProvisioner{c}
// 			err := provisioner.Provision(item.kubeResource)
//
// 			So(err, ShouldResemble, item.errorReturned)
// 			So(kindPassed, ShouldEqual, item.kindPassed)
// 			So(namespacePassed, ShouldEqual, item.namespacePassed)
//
// 			// NOTE I don't know wtf is going on here
// 			// So(objInPassed, ShouldResemble, item.objInPassed)
// 			objInPassedMarshalled, _ := json.Marshal(objInPassed)
// 			itemObjInPassedMarshalled, _ := json.Marshal(item.objInPassed)
//
// 			So(objInPassedMarshalled, ShouldResemble, itemObjInPassedMarshalled)
//
// 			So(outSaved, ShouldResemble, item.outSaved)
// 		}
// 	})
// }

//------------------------------------------------------------------------------

// func TestDefaultProvisionerTeardown(t *testing.T) {
// 	Convey("DefaultProvisioner Teardown passses correct arguments to KubernetesClient DeleteResource(), and passes error through", t, func() {
// 		table := []struct {
// 			// Input
// 			kubeResource *model.KubeResource
//
// 			// Mocks
// 			mockDeleteResourceError error
//
// 			// Assertions
// 			kindPassed      string
// 			namespacePassed string
// 			namePassed      string
// 			errorReturned   error
// 		}{
// 			// A successful example
// 			//------------------------------------------------------------------------
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "fizz",
// 					Name:      "bubb",
// 					Kind:      "Bleeper",
// 				},
// 				// Mocks
// 				mockDeleteResourceError: nil,
// 				// Assertions
// 				kindPassed:      "Bleeper",
// 				namespacePassed: "fizz",
// 				namePassed:      "bubb",
// 				errorReturned:   nil,
// 			},
//
// 			// When there's a 404 error (we don't care)
// 			//------------------------------------------------------------------------
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "test",
// 					Name:      "test",
// 					Kind:      "Kind",
// 				},
// 				// Mocks
// 				mockDeleteResourceError: errors.New("K8S 404 Not Found error"),
// 				// Assertions
// 				kindPassed:      "Kind",
// 				namespacePassed: "test",
// 				namePassed:      "test",
// 				errorReturned:   nil,
// 			},
//
// 			// When there's an unexpected error from Kubernetes
// 			//------------------------------------------------------------------------
// 			{
// 				// Input
// 				kubeResource: &model.KubeResource{
// 					Namespace: "chip",
// 					Name:      "chop",
// 					Kind:      "Krampus",
// 				},
// 				// Mocks
// 				mockDeleteResourceError: errors.New("wots this"),
// 				// Assertions
// 				kindPassed:      "Krampus",
// 				namespacePassed: "chip",
// 				namePassed:      "chop",
// 				errorReturned:   errors.New("wots this"),
// 			},
// 		}
//
// 		for _, item := range table {
// 			var kindPassed string
// 			var namespacePassed string
// 			var namePassed string
//
// 			c := &core.Core{
// 				K8S: func(_ *model.Kube) kubernetes.ClientInterface {
// 					return &fake_core.KubernetesClient{
// 						DeleteResourceFn: func(apiVersion string, kind string, namespace string, name string) error {
// 							kindPassed = kind
// 							namespacePassed = namespace
// 							namePassed = name
//
// 							return item.mockDeleteResourceError
// 						},
// 					}
// 				},
// 			}
//
// 			provisioner := &core.DefaultProvisioner{c}
// 			err := provisioner.Teardown(item.kubeResource)
//
// 			So(err, ShouldResemble, item.errorReturned)
// 			So(kindPassed, ShouldEqual, item.kindPassed)
// 			So(namespacePassed, ShouldEqual, item.namespacePassed)
// 			So(namePassed, ShouldEqual, item.namePassed)
// 		}
// 	})
// }
