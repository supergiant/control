package core_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/test/fake_core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/pkg/provider/digitalocean"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPodProvisionerProvision(t *testing.T) {
	Convey("PodProvisioner Provision works correctly", t, func() {

		var kubeResourceID int64 = 16 // just need a number we can grab a pointer to

		table := []struct {
			// Input
			kubeResource *model.KubeResource

			// Mocks
			mockExistingVolumes          []*model.Volume
			mockDefaultProvisionError    error
			mockVolumeFetchError         error
			mockVolumeCreateError        error
			mockVolumeProviderIDAssigner func(*model.Volume)
			mockPodStartTimeout          bool

			// Expectations
			volumesCreated                       []*model.Volume
			definitionPassedToDefaultProvisioner map[string]interface{}
			errorReturned                        error
		}{
			// A basic, successful example
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "aws",
						},
					},
					KubeName:  "mykube",
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "my-lil-volume",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "gp2",
										"size": 10
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes:       nil,
				mockDefaultProvisionError: nil,
				mockVolumeFetchError:      nil,
				mockVolumeCreateError:     nil,
				mockVolumeProviderIDAssigner: func(volume *model.Volume) {
					volume.ProviderID = "zzzzzz"
				},
				// Expectations
				volumesCreated: []*model.Volume{
					{
						KubeName:       "mykube",
						KubeResourceID: &kubeResourceID,
						Name:           "my-lil-volume",
						Type:           "gp2",
						Size:           10,
						ProviderID:     "zzzzzz",
					},
				},
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"volumes": []interface{}{
							map[string]interface{}{
								"name": "my-lil-volume",
								"awsElasticBlockStore": map[string]interface{}{
									"volumeID": "zzzzzz",
									"fsType":   "ext4",
								},
							},
						},
					},
				},
				errorReturned: nil,
			},

			// A successful example where some Volumes have already been created
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "aws",
						},
					},
					KubeName:  "mykube",
					Namespace: "fooz",
					Name:      "bars",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "vol0",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "io1",
										"size": 120
									}
								},
								{
									"name": "vol1",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "io1",
										"size": 15
									}
								},
								{
									"name": "vol2",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "gp2",
										"size": 33
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes: []*model.Volume{
					{
						KubeName:       "mykube",
						KubeResourceID: &kubeResourceID,
						Name:           "vol1",
						Type:           "io1",
						Size:           15,
						ProviderID:     "vol1-ID",
					},
				},
				mockDefaultProvisionError: nil,
				mockVolumeFetchError:      nil,
				mockVolumeCreateError:     nil,
				mockVolumeProviderIDAssigner: func(volume *model.Volume) {
					volume.ProviderID = volume.Name + "-ID"
				},
				// Expectations
				volumesCreated: []*model.Volume{
					{
						KubeName:       "mykube",
						KubeResourceID: &kubeResourceID,
						Name:           "vol0",
						Type:           "io1",
						Size:           120,
						ProviderID:     "vol0-ID",
					},
					{
						KubeName:       "mykube",
						KubeResourceID: &kubeResourceID,
						Name:           "vol2",
						Type:           "gp2",
						Size:           33,
						ProviderID:     "vol2-ID",
					},
				},
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"volumes": []interface{}{
							map[string]interface{}{
								"name": "vol0",
								"awsElasticBlockStore": map[string]interface{}{
									"volumeID": "vol0-ID",
									"fsType":   "ext4",
								},
							},
							map[string]interface{}{
								"name": "vol1",
								"awsElasticBlockStore": map[string]interface{}{
									"volumeID": "vol1-ID",
									"fsType":   "ext4",
								},
							},
							map[string]interface{}{
								"name": "vol2",
								"awsElasticBlockStore": map[string]interface{}{
									"volumeID": "vol2-ID",
									"fsType":   "ext4",
								},
							},
						},
					},
				},
				errorReturned: nil,
			},

			// A successful example where there are other non-SG volumes defined
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "aws",
						},
					},
					KubeName:  "mykube",
					Namespace: "fooz",
					Name:      "bars",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "vol0",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "io1",
										"size": 120
									}
								},
								{
									"name": "non-sg-vol",
									"emptyDir": {}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes: []*model.Volume{
					{
						KubeName:       "mykube",
						KubeResourceID: &kubeResourceID,
						Name:           "vol0",
						Type:           "io1",
						Size:           120,
						ProviderID:     "vol0-ID",
					},
				},
				mockDefaultProvisionError:    nil,
				mockVolumeFetchError:         nil,
				mockVolumeCreateError:        nil,
				mockVolumeProviderIDAssigner: nil,
				// Expectations
				volumesCreated: nil,
				definitionPassedToDefaultProvisioner: map[string]interface{}{
					"spec": map[string]interface{}{
						"volumes": []interface{}{
							map[string]interface{}{
								"name": "vol0",
								"awsElasticBlockStore": map[string]interface{}{
									"volumeID": "vol0-ID",
									"fsType":   "ext4",
								},
							},
							map[string]interface{}{
								"name":     "non-sg-vol",
								"emptyDir": map[string]interface{}{},
							},
						},
					},
				},
				errorReturned: nil,
			},

			// When there's an error fetching existing Volumes
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "digitalocean",
						},
					},
					KubeName:  "nothing",
					Namespace: "really",
					Name:      "matters",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "vol0",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "io1",
										"size": 120
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes:          nil,
				mockDefaultProvisionError:    nil,
				mockVolumeFetchError:         errors.New("well what's this here"),
				mockVolumeCreateError:        nil,
				mockVolumeProviderIDAssigner: nil,
				// Expectations
				volumesCreated:                       nil,
				definitionPassedToDefaultProvisioner: nil,
				errorReturned:                        errors.New("well what's this here"),
			},

			// When there's an error creating Volumes (due to user-input or unknown)
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "digitalocean",
						},
					},
					KubeName:  "test",
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "vol0",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "io1",
										"size": 120
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes:          nil,
				mockDefaultProvisionError:    nil,
				mockVolumeFetchError:         nil,
				mockVolumeCreateError:        errors.New("bad"),
				mockVolumeProviderIDAssigner: nil,
				// Expectations
				volumesCreated:                       nil,
				definitionPassedToDefaultProvisioner: nil,
				errorReturned:                        errors.New("bad"),
			},

			// When there's an error with the DefaultProvisioner (most likely unknown)
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					BaseModel: model.BaseModel{
						ID: &kubeResourceID,
					},
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Provider: "digitalocean",
						},
					},
					KubeName:  "test",
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
					Template: newRawMessage(`{
						"spec": {
							"volumes": [
								{
									"name": "vol",
									"SUPERGIANT_EXTERNAL_VOLUME": {
										"type": "gp2",
										"size": 10
									}
								}
							]
						}
					}`),
				},
				// Mocks
				mockExistingVolumes:       nil,
				mockDefaultProvisionError: errors.New("really not good this time"),
				mockVolumeFetchError:      nil,
				mockVolumeCreateError:     nil,
				mockVolumeProviderIDAssigner: func(volume *model.Volume) {
					volume.ProviderID = volume.Name + "-id"
				},
				// Expectations
				volumesCreated: []*model.Volume{
					{
						KubeName:       "test",
						KubeResourceID: &kubeResourceID,
						Name:           "vol",
						Type:           "gp2",
						Size:           10,
						ProviderID:     "vol-id",
					},
				},
				definitionPassedToDefaultProvisioner: nil,
				errorReturned:                        errors.New("really not good this time"),
			},
		}

		for _, item := range table {
			var volumesCreated []*model.Volume
			var definitionPassedToDefaultProvisioner map[string]interface{}

			c := &core.Core{
				// We call K8S client to wait for Pod start
				K8S: func(_ *model.Kube) kubernetes.ClientInterface {
					return &fake_core.KubernetesClient{
						GetResourceFn: func(kind string, namespace string, name string, out *json.RawMessage) error {
							status := "True"
							if item.mockPodStartTimeout {
								status = "False"
							}
							*out = json.RawMessage([]byte(fmt.Sprintf(`{"status": {"conditions": [{"type": "Ready", "status": "%s"}]}}`, status)))
							return nil
						},
					}
				},

				DB: &fake_core.DB{
					// Used to fetch Volumes
					FindFn: func(out interface{}, _ ...interface{}) error {
						if err := item.mockVolumeFetchError; err != nil {
							return err
						}
						// All this insane line does is set out = item.mockExistingVolumes
						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.mockExistingVolumes))
						return nil
					},
				},

				Volumes: &fake_core.Volumes{
					CreateFn: func(volume *model.Volume) error {
						if err := item.mockVolumeCreateError; err != nil {
							return err
						}
						volumesCreated = append(volumesCreated, volume)
						if idFn := item.mockVolumeProviderIDAssigner; idFn != nil {
							idFn(volume)
						}
						return nil
					},
				},

				DefaultProvisioner: &fake_core.Provisioner{
					ProvisionFn: func(kubeResource *model.KubeResource) error {
						if err := item.mockDefaultProvisionError; err != nil {
							return err
						}
						json.Unmarshal(*kubeResource.Definition, &definitionPassedToDefaultProvisioner)
						return nil
					},
				},
			}

			// Needed to load up the appropriate Provider
			c.CloudAccounts = &core.CloudAccounts{core.Collection{Core: c}}
			// We can use the real provider here, just need it for Volume def
			c.AWSProvider = func(creds map[string]string) core.Provider {
				return &aws.Provider{Core: c, Credentials: creds}
			}
			c.DOProvider = func(_ map[string]string) core.Provider {
				return &digitalocean.Provider{Core: c}
			}

			provisioner := &core.PodProvisioner{c}
			err := provisioner.Provision(item.kubeResource)

			So(err, ShouldResemble, item.errorReturned)
			So(volumesCreated, ShouldResemble, item.volumesCreated)
			So(definitionPassedToDefaultProvisioner, ShouldResemble, item.definitionPassedToDefaultProvisioner)
		}
	})
}

//------------------------------------------------------------------------------

func TestPodProvisionerTeardown(t *testing.T) {
	Convey("PodProvisioner Teardown works correctly", t, func() {
		table := []struct {
			// Input
			kubeResource *model.KubeResource

			// Mocks
			mockOwnedVolumes         []*model.Volume
			mockVolumeFetchError     error
			mockVolumeDeleteError    error
			mockDefaultTeardownError error

			// Assertions
			volumeNamesDeleted    []string
			defaultTeardownCalled bool
			errorReturned         error
		}{
			// A successful example
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
				},
				// Mocks
				mockOwnedVolumes: []*model.Volume{
					{
						Name: "one-volume",
					},
					{
						Name: "two-volume",
					},
				},
				mockVolumeFetchError:     nil,
				mockVolumeDeleteError:    nil,
				mockDefaultTeardownError: nil,
				// Assertions
				volumeNamesDeleted:    []string{"one-volume", "two-volume"},
				defaultTeardownCalled: true,
				errorReturned:         nil,
			},

			// When there's an error fetching Volumes
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
				},
				// Mocks
				mockOwnedVolumes: []*model.Volume{
					{
						Name: "this-volume",
					},
				},
				mockVolumeFetchError:     errors.New("crap"),
				mockVolumeDeleteError:    nil,
				mockDefaultTeardownError: nil,
				// Assertions
				volumeNamesDeleted:    nil,
				defaultTeardownCalled: false,
				errorReturned:         errors.New("crap"),
			},

			// When there's an error deleting Volumes
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
				},
				// Mocks
				mockOwnedVolumes: []*model.Volume{
					{
						Name: "this-volume",
					},
				},
				mockVolumeFetchError:     nil,
				mockVolumeDeleteError:    errors.New("whoops"),
				mockDefaultTeardownError: nil,
				// Assertions
				volumeNamesDeleted:    nil,
				defaultTeardownCalled: true,
				errorReturned:         errors.New("whoops"),
			},

			// When there's an error from default Teardown
			//------------------------------------------------------------------------
			{
				// Input
				kubeResource: &model.KubeResource{
					Namespace: "test",
					Name:      "test",
					Kind:      "Pod",
				},
				// Mocks
				mockOwnedVolumes: []*model.Volume{
					{
						Name: "this-volume",
					},
				},
				mockVolumeFetchError:     nil,
				mockVolumeDeleteError:    nil,
				mockDefaultTeardownError: errors.New("splat"),
				// Assertions
				volumeNamesDeleted:    nil,
				defaultTeardownCalled: true,
				errorReturned:         errors.New("splat"),
			},
		}

		for _, item := range table {

			var volumeNamesDeleted []string
			var defaultTeardownCalled bool

			c := &core.Core{
				DB: &fake_core.DB{
					// Used to fetch Volumes
					FindFn: func(out interface{}, _ ...interface{}) error {
						if err := item.mockVolumeFetchError; err != nil {
							return err
						}
						// All this insane line does is set out = item.mockExistingVolumes
						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(item.mockOwnedVolumes))
						return nil
					},
				},

				DefaultProvisioner: &fake_core.Provisioner{
					TeardownFn: func(kubeResource *model.KubeResource) error {
						if kubeResource.Name == item.kubeResource.Name { // just for sanity
							defaultTeardownCalled = true
						}
						return item.mockDefaultTeardownError
					},
				},
			}
			c.Volumes = &fake_core.Volumes{
				DeleteFn: func(_ *int64, volume *model.Volume) core.ActionInterface {
					return &fake_core.Action{
						NowFn: func() error {
							if err := item.mockVolumeDeleteError; err != nil {
								return err
							}
							volumeNamesDeleted = append(volumeNamesDeleted, volume.Name)
							return nil
						},
					}
				},
			}

			provisioner := &core.PodProvisioner{c}
			err := provisioner.Teardown(item.kubeResource)

			So(err, ShouldResemble, item.errorReturned)

			// So(volumeNamesDeleted, ShouldResemble, item.volumeNamesDeleted)
			for _, volumeNameDeleted := range item.volumeNamesDeleted {
				So(volumeNamesDeleted, ShouldContain, volumeNameDeleted)
			}

			So(defaultTeardownCalled, ShouldEqual, item.defaultTeardownCalled)
		}
	})
}
