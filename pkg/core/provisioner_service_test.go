package core_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/core/fake"
	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceProvisionerProvision(t *testing.T) {
	Convey("ServiceProvisioner Provision works correctly", t, func() {

		var kubeResourceID int64 = 7 // just need a number we can grab a pointer to

		table := []struct {
			// Input
			kubeResource *model.KubeResource

			// Mocks
			mockExistingEntrypointListeners   []*model.EntrypointListener
			mockDefaultProvisionError         error
			mockEntrypointListenerFetchError  error
			mockEntrypointListenerCreateError error
			mockEntrypointListenerDeleteError error
			mockNodePortAssigner              func(portDef map[string]interface{})

			// Expectations
			definitionPassedToDefaultProvisioner map[string]interface{}
			entrypointListenersCreated           []*model.EntrypointListener
			entrypointListenerNamesDeleted       []string
			errorReturned                        error
		}{
			// A basic, successful example (starting from scratch)
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					KubeName:  "kube2",
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "website-http",
									"port": 8080,
									"protocol": "TCP",
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "my-entrypoint",
										"entrypoint_port": 80,
										"entrypoint_protocol": "HTTP"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners:   nil,
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner: func(portDef map[string]interface{}) {
					portDef["nodePort"] = portDef["port"].(float64) + 9
				},
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name":     "website-http",
								"port":     8080,
								"protocol": "TCP",
							},
						},
					},
				},
				entrypointListenersCreated: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "my-entrypoint",
						Name:               "website-http",
						EntrypointPort:     80,
						EntrypointProtocol: "HTTP",
						NodePort:           8089,
						NodeProtocol:       "TCP",
					},
				},
				entrypointListenerNamesDeleted: nil,
				errorReturned:                  nil,
			},

			// When there are existing EntrypointListeners, and no changes
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "website-http",
									"port": 8080,
									"protocol": "TCP",
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "my-entrypoint",
										"entrypoint_port": 80,
										"entrypoint_protocol": "HTTP"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "my-entrypoint",
						Name:               "website-http",
						EntrypointPort:     80,
						EntrypointProtocol: "HTTP",
						NodePort:           36277,
						NodeProtocol:       "TCP",
					},
				},
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name":     "website-http",
								"port":     8080,
								"nodePort": 36277,
								"protocol": "TCP",
							},
						},
					},
				},
				entrypointListenersCreated:     nil,
				entrypointListenerNamesDeleted: nil,
				errorReturned:                  nil,
			},

			// When there are other non-SG ports defined
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "website-http",
									"port": 8080,
									"protocol": "TCP",
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "my-entrypoint",
										"entrypoint_port": 80,
										"entrypoint_protocol": "HTTP"
									}
								},
								{
									"name": "other",
									"port": 2277
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "my-entrypoint",
						Name:               "website-http",
						EntrypointPort:     80,
						EntrypointProtocol: "HTTP",
						NodePort:           36277,
						NodeProtocol:       "TCP",
					},
				},
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name":     "website-http",
								"port":     8080,
								"nodePort": 36277,
								"protocol": "TCP",
							},
							map[string]interface{}{
								"name": "other",
								"port": 2277,
							},
						},
					},
				},
				entrypointListenersCreated:     nil,
				entrypointListenerNamesDeleted: nil,
				errorReturned:                  nil,
			},

			// When there are existing EntrypointListeners that are no longer defined
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "cool-new-port",
									"port": 8080,
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "my-entrypoint",
										"entrypoint_port": 80
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "my-entrypoint",
						Name:               "website-http",
						EntrypointPort:     80,
						EntrypointProtocol: "HTTP",
						NodePort:           36277,
					},
				},
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner: func(portDef map[string]interface{}) {
					portDef["nodePort"] = portDef["port"].(float64) + 100
				},
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name": "cool-new-port",
								"port": 8080,
							},
						},
					},
				},
				entrypointListenersCreated: []*model.EntrypointListener{
					{
						KubeResourceID: &kubeResourceID,
						EntrypointName: "my-entrypoint",
						Name:           "cool-new-port",
						EntrypointPort: 80,
						NodePort:       8180,
					},
				},
				entrypointListenerNamesDeleted: []string{"website-http"},
				errorReturned:                  nil,
			},

			// When there are existing EntrypointListeners, and changes to apply
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "my-port",
									"port": 7777,
									"protocol": "UDP",
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "new-entrypoint",
										"entrypoint_port": 443,
										"entrypoint_protocol": "HTTPS"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "old-entrypoint",
						Name:               "my-port",
						EntrypointPort:     443,
						EntrypointProtocol: "HTTPS",
						NodePort:           33333,
						NodeProtocol:       "UDP",
					},
				},
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name":     "my-port",
								"port":     7777,
								"protocol": "UDP",
								"nodePort": 33333,
							},
						},
					},
				},
				entrypointListenersCreated: []*model.EntrypointListener{
					{
						KubeResourceID:     &kubeResourceID,
						EntrypointName:     "new-entrypoint",
						Name:               "my-port",
						EntrypointPort:     443,
						EntrypointProtocol: "HTTPS",
						NodePort:           33333,
						NodeProtocol:       "UDP",
					},
				},
				entrypointListenerNamesDeleted: []string{"my-port"},
				errorReturned:                  nil,
			},

			// When Default Provision returns an error
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "porty",
									"port": 7777,
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "new-entrypoint",
										"entrypoint_port": 443,
										"entrypoint_protocol": "HTTPS"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners:   nil,
				mockDefaultProvisionError:         errors.New("some error"),
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name": "porty",
								"port": 7777,
							},
						},
					},
				},
				entrypointListenersCreated:     nil,
				entrypointListenerNamesDeleted: nil,
				errorReturned:                  errors.New("some error"),
			},

			// When there's an error fetching existing EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "porty",
									"port": 7777,
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "new-entrypoint",
										"entrypoint_port": 443,
										"entrypoint_protocol": "HTTPS"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners:   nil,
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  errors.New("some error"),
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: nil,
				entrypointListenersCreated:           nil,
				entrypointListenerNamesDeleted:       nil,
				errorReturned:                        errors.New("some error"),
			},

			// When there's an error creating EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": [
								{
									"name": "porty",
									"port": 7777,
									"SUPERGIANT_ENTRYPOINT_LISTENER": {
										"entrypoint_name": "new-entrypoint",
										"entrypoint_port": 443,
										"entrypoint_protocol": "HTTPS"
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners:   nil,
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: errors.New("some error"),
				mockEntrypointListenerDeleteError: nil,
				mockNodePortAssigner: func(portDef map[string]interface{}) {
					portDef["nodePort"] = portDef["port"].(float64) + 1
				},
				// Expectations
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"ports": []interface{}{
							map[string]interface{}{
								"name": "porty",
								"port": 7777,
							},
						},
					},
				},
				entrypointListenersCreated:     nil,
				entrypointListenerNamesDeleted: nil,
				errorReturned:                  errors.New("some error"),
			},

			// When there's an error deleting EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
					Template: newRawMessage(`{
						"spec": {
							"ports": []
						}
					}`),
				},
				// Mocks
				mockExistingEntrypointListeners: []*model.EntrypointListener{
					{
						KubeResourceID: &kubeResourceID,
						EntrypointName: "entrypoint",
						Name:           "port",
					},
				},
				mockDefaultProvisionError:         nil,
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerCreateError: nil,
				mockEntrypointListenerDeleteError: errors.New("some error"),
				mockNodePortAssigner:              nil,
				// Expectations
				definitionPassedToDefaultProvisioner: nil,
				entrypointListenersCreated:           nil,
				entrypointListenerNamesDeleted:       nil,
				errorReturned:                        errors.New("some error"),
			},
		}

		for _, item := range table {
			var definitionPassedToDefaultProvisioner map[string]interface{}
			var entrypointListenersCreated []*model.EntrypointListener
			var entrypointListenerNamesDeleted []string

			c := &core.Core{
				DB: &fake.DB{
					// Used to fetch EntrypointListeners
					FindFn: func(out interface{}, _ ...interface{}) error {
						if err := item.mockEntrypointListenerFetchError; err != nil {
							return err
						}
						// All this insane line does is set out = item.mockExistingEntrypointListeners
						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.mockExistingEntrypointListeners))
						return nil
					},
				},

				EntrypointListeners: &fake.EntrypointListeners{
					CreateFn: func(entrypointListener *model.EntrypointListener) error {
						if err := item.mockEntrypointListenerCreateError; err != nil {
							return err
						}
						entrypointListenersCreated = append(entrypointListenersCreated, entrypointListener)
						return nil
					},

					// Delete is used during Provision if there are orphaned Listeners due
					// to Template change
					DeleteFn: func(_ *int64, entrypointListener *model.EntrypointListener) core.ActionInterface {
						return &fake.Action{
							NowFn: func() error {
								if err := item.mockEntrypointListenerDeleteError; err != nil {
									return err
								}
								entrypointListenerNamesDeleted = append(entrypointListenerNamesDeleted, entrypointListener.Name)
								return nil
							},
						}
					},
				},

				DefaultProvisioner: &fake.Provisioner{
					ProvisionFn: func(kubeResource *model.KubeResource) error {

						json.Unmarshal(*kubeResource.Definition, &definitionPassedToDefaultProvisioner)

						if err := item.mockDefaultProvisionError; err != nil {
							return err
						}

						// We copy Definition to Artifact here to emulate how the Artifact
						// is unmarshalled to by kubernetes CreateResource
						artifact := make(json.RawMessage, len(*kubeResource.Definition))
						copy(artifact, *kubeResource.Definition)
						kubeResource.Artifact = &artifact

						var artifactMap map[string]interface{}
						json.Unmarshal(*kubeResource.Artifact, &artifactMap)

						// We do this here, because inside this function is where we know
						// the NodePort should be assigned.
						spec := artifactMap["spec"].(map[string]interface{})
						ports := spec["ports"].([]interface{})
						var newPorts []interface{}
						for _, p := range ports {
							port := p.(map[string]interface{})
							if port["nodePort"] == nil && item.mockNodePortAssigner != nil {
								item.mockNodePortAssigner(port)
							}
							newPorts = append(newPorts, port)
						}
						spec["ports"] = newPorts

						marshalledArtifact, _ := json.Marshal(artifactMap)
						rawArtifact := json.RawMessage(marshalledArtifact)
						kubeResource.Artifact = &rawArtifact

						return nil
					},
				},
			}

			provisioner := &core.ServiceProvisioner{c}
			err := provisioner.Provision(item.kubeResource)

			So(err, ShouldResemble, item.errorReturned)

			// NOTE I don't know wtf is going on here
			// So(definitionPassedToDefaultProvisioner, ShouldResemble, item.definitionPassedToDefaultProvisioner)
			definitionPassedToDefaultProvisionerMarshalled, _ := json.Marshal(definitionPassedToDefaultProvisioner)
			itemDefinitionPassedToDefaultProvisionerMarshalled, _ := json.Marshal(item.definitionPassedToDefaultProvisioner)

			So(definitionPassedToDefaultProvisionerMarshalled, ShouldResemble, itemDefinitionPassedToDefaultProvisionerMarshalled)

			So(entrypointListenersCreated, ShouldResemble, item.entrypointListenersCreated)
			So(entrypointListenerNamesDeleted, ShouldResemble, item.entrypointListenerNamesDeleted)
		}
	})
}

//------------------------------------------------------------------------------

func TestServiceProvisionerTeardown(t *testing.T) {
	Convey("ServiceProvisioner Teardown works correctly", t, func() {
		table := []struct {
			// Input
			kubeResource *model.KubeResource

			// Mocks
			mockOwnedEntrypointListeners      []*model.EntrypointListener
			mockEntrypointListenerFetchError  error
			mockEntrypointListenerDeleteError error
			mockDefaultTeardownError          error

			// Assertions
			entrypointListenerNamesDeleted []string
			defaultTeardownCalled          bool
			errorReturned                  error
		}{
			// A successful example
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
				},
				// Mocks
				mockOwnedEntrypointListeners: []*model.EntrypointListener{
					{
						Name: "one-port",
					},
					{
						Name: "two-port",
					},
				},
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerDeleteError: nil,
				mockDefaultTeardownError:          nil,
				// Assertions
				entrypointListenerNamesDeleted: []string{"one-port", "two-port"},
				defaultTeardownCalled:          true,
				errorReturned:                  nil,
			},

			// When there's an error fetching EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
				},
				// Mocks
				mockOwnedEntrypointListeners: []*model.EntrypointListener{
					{
						Name: "this-port",
					},
				},
				mockEntrypointListenerFetchError:  errors.New("crap"),
				mockEntrypointListenerDeleteError: nil,
				mockDefaultTeardownError:          nil,
				// Assertions
				entrypointListenerNamesDeleted: nil,
				defaultTeardownCalled:          false,
				errorReturned:                  errors.New("crap"),
			},

			// When there's an error deleting EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
				},
				// Mocks
				mockOwnedEntrypointListeners: []*model.EntrypointListener{
					{
						Name: "this-port",
					},
				},
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerDeleteError: errors.New("whoops"),
				mockDefaultTeardownError:          nil,
				// Assertions
				entrypointListenerNamesDeleted: nil,
				defaultTeardownCalled:          false,
				errorReturned:                  errors.New("whoops"),
			},

			// When there's an error deleting EntrypointListeners
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Service",
				},
				// Mocks
				mockOwnedEntrypointListeners: []*model.EntrypointListener{
					{
						Name: "this-port",
					},
				},
				mockEntrypointListenerFetchError:  nil,
				mockEntrypointListenerDeleteError: nil,
				mockDefaultTeardownError:          errors.New("splat"),
				// Assertions
				entrypointListenerNamesDeleted: []string{"this-port"},
				defaultTeardownCalled:          true,
				errorReturned:                  errors.New("splat"),
			},
		}

		for _, item := range table {

			var entrypointListenerNamesDeleted []string
			var defaultTeardownCalled bool

			c := &core.Core{
				DB: &fake.DB{
					// Used to fetch EntrypointListeners
					FindFn: func(out interface{}, _ ...interface{}) error {
						if err := item.mockEntrypointListenerFetchError; err != nil {
							return err
						}
						// All this insane line does is set out = item.mockExistingEntrypointListeners
						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.mockOwnedEntrypointListeners))
						return nil
					},
				},

				DefaultProvisioner: &fake.Provisioner{
					TeardownFn: func(kubeResource *model.KubeResource) error {
						if kubeResource.Name == item.kubeResource.Name { // just for sanity
							defaultTeardownCalled = true
						}
						return item.mockDefaultTeardownError
					},
				},
			}
			c.EntrypointListeners = &fake.EntrypointListeners{
				DeleteFn: func(_ *int64, entrypointListener *model.EntrypointListener) core.ActionInterface {
					return &fake.Action{
						NowFn: func() error {
							if err := item.mockEntrypointListenerDeleteError; err != nil {
								return err
							}
							entrypointListenerNamesDeleted = append(entrypointListenerNamesDeleted, entrypointListener.Name)
							return nil
						},
					}
				},
			}

			provisioner := &core.ServiceProvisioner{c}
			err := provisioner.Teardown(item.kubeResource)

			So(err, ShouldResemble, item.errorReturned)
			// So(entrypointListenerNamesDeleted, ShouldResemble, item.entrypointListenerNamesDeleted)
			for _, entrypointListenerNameDeleted := range item.entrypointListenerNamesDeleted {
				So(entrypointListenerNamesDeleted, ShouldContain, entrypointListenerNameDeleted)
			}
			So(defaultTeardownCalled, ShouldEqual, item.defaultTeardownCalled)
		}
	})
}
