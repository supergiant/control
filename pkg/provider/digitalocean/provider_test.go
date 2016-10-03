package digitalocean_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/digitalocean"
	"github.com/supergiant/supergiant/test/fake_core"
	"github.com/supergiant/supergiant/test/fake_digitalocean_provider"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDigitalOceanProviderValidateAccount(t *testing.T) {
	Convey("DigitalOcean Provider ValidateAccount works correctly", t, func() {
		table := []struct {
			// Input
			cloudAccount *model.CloudAccount
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				cloudAccount: &model.CloudAccount{},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							ListFn: func(_ *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
								return nil, nil, nil
							},
						},
					}
				},
			}

			err := provider.ValidateAccount(item.cloudAccount)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderCreateKube(t *testing.T) {
	Convey("DigitalOcean Provider CreateKube works correctly", t, func() {
		table := []struct {
			// Input
			kube *model.Kube
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube: &model.Kube{
					NodeSizes:          []string{"2gb"},
					DigitalOceanConfig: &model.DOKubeConfig{},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),

				K8S: func(*model.Kube) kubernetes.ClientInterface {
					return &fake_core.KubernetesClient{
						ListNodesFn: func(query string) ([]*kubernetes.Node, error) {
							return []*kubernetes.Node{
								{
									Metadata: kubernetes.Metadata{
										Name: "created-node",
									},
								},
							}, nil
						},
					}
				},

				Nodes: new(fake_core.Nodes),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							// Create
							CreateFn: func(_ *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
								return &godo.Droplet{
									ID: 1,
								}, nil, nil
							},
							// Get
							GetFn: func(int) (*godo.Droplet, *godo.Response, error) {
								return &godo.Droplet{
									ID: 1,
									Networks: &godo.Networks{
										V4: []godo.NetworkV4{
											{
												Type:      "public",
												IPAddress: "99.99.99.99",
											},
										},
									},
								}, nil, nil
							},
						},
						Tags: &fake_digitalocean_provider.Tags{},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateKube(item.kube, action)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderDeleteKube(t *testing.T) {
	Convey("DigitalOcean Provider DeleteKube works correctly", t, func() {
		table := []struct {
			// Input
			kube *model.Kube
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"2gb"},
					DigitalOceanConfig: &model.DOKubeConfig{
						MasterID: 16678899,
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							// Delete
							DeleteFn: func(_ int) (*godo.Response, error) {
								return nil, nil
							},
						},
					}
				},
			}

			err := provider.DeleteKube(item.kube)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderCreateNode(t *testing.T) {
	Convey("DigitalOcean Provider CreateNode works correctly", t, func() {
		table := []struct {
			// Input
			node *model.Node
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				node: &model.Node{
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Credentials: map[string]string{"token": "my-special-token"},
						},
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							// Create
							CreateFn: func(_ *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
								return &godo.Droplet{
									ID: 1,
								}, nil, nil
							},
							// Get
							GetFn: func(int) (*godo.Droplet, *godo.Response, error) {
								return &godo.Droplet{
									ID: 1,
									Networks: &godo.Networks{
										V4: []godo.NetworkV4{
											{
												Type:      "public",
												IPAddress: "99.99.99.99",
											},
										},
									},
								}, nil, nil
							},
						},
						Tags: &fake_digitalocean_provider.Tags{},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateNode(item.node, action)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderDeleteNode(t *testing.T) {
	Convey("DigitalOcean Provider DeleteNode works correctly", t, func() {
		table := []struct {
			// Input
			node *model.Node
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				node: &model.Node{
					ProviderID: "1",
					Kube: &model.Kube{
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							// Delete
							DeleteFn: func(_ int) (*godo.Response, error) {
								return nil, nil
							},
						},
					}
				},
			}

			err := provider.DeleteNode(item.node)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderCreateVolume(t *testing.T) {
	Convey("DigitalOcean Provider CreateVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Credentials: map[string]string{"token": "my-special-token"},
						},
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Storage: &fake_digitalocean_provider.Storage{
							CreateVolumeFn: func(*godo.VolumeCreateRequest) (*godo.Volume, *godo.Response, error) {
								return &godo.Volume{
									ID: "volumeID",
								}, nil, nil
							},
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateVolume(item.volume, action)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderKubernetesVolumeDefinition(t *testing.T) {
	Convey("DigitalOcean Provider KubernetesVolumeDefinition works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			kubeVol *kubernetes.Volume
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					Name:       "test",
					ProviderID: "provider-ID",
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Credentials: map[string]string{"token": "my-special-token"},
						},
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
				// Expectations
				kubeVol: &kubernetes.Volume{
					Name: "test",
					FlexVolume: &kubernetes.FlexVolume{
						Driver: "supergiant.io/digitalocean",
						FSType: "ext4",
						Options: map[string]string{
							"volumeID": "provider-ID",
							"name":     "test",
						},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
			}

			kubeVol := provider.KubernetesVolumeDefinition(item.volume)

			So(kubeVol, ShouldResemble, item.kubeVol)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderResizeVolume(t *testing.T) {
	Convey("DigitalOcean Provider ResizeVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					ProviderID: "volumeID",
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Credentials: map[string]string{"token": "my-special-token"},
						},
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						StorageActions: &fake_digitalocean_provider.StorageActions{
							ResizeFn: func(volumeID string, sizeGigabytes int, regionSlug string) (*godo.Action, *godo.Response, error) {
								return nil, nil, nil
							},
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.ResizeVolume(item.volume, action)

			So(err, ShouldEqual, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestDigitalOceanProviderDeleteVolume(t *testing.T) {
	Convey("DigitalOcean Provider DeleteVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					ProviderID: "volumeID",
					Kube: &model.Kube{
						CloudAccount: &model.CloudAccount{
							Credentials: map[string]string{"token": "my-special-token"},
						},
						DigitalOceanConfig: &model.DOKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Storage: &fake_digitalocean_provider.Storage{
							DeleteVolumeFn: func(_ string) (*godo.Response, error) {
								return nil, nil
							},
						},
					}
				},
			}

			err := provider.DeleteVolume(item.volume)

			So(err, ShouldEqual, item.err)
		}
	})
}
