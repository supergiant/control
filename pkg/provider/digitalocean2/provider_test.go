package digitalocean2

import (
	"testing"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/test/fake_core"
	"github.com/Sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/digitalocean/godo"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/supergiant/pkg/provision"
	"github.com/stretchr/testify/require"
)

type DOClientMock struct {
	mock.Mock
}

var _ DigitalOceanClient = (*DOClientMock)(nil)

func (m *DOClientMock) NewDroplet(k *model.Kube, s string, c context.Context) (*godo.Droplet, error) {
	args := m.Called(k, s, c)
	dr := args.Get(0).(godo.Droplet)
	return &dr, args.Error(1)
}

func (m *DOClientMock) DeleteDroplet(ID DropletID, c context.Context) (error) {
	args := m.Called(ID, c)
	return args.Error(0)
}

type ProvisionMock struct {
	mock.Mock
}

var _ provision.Interface = (*ProvisionMock)(nil)

func (m *ProvisionMock) CreateMaster(kube *model.Kube, ips []string, ctx context.Context) error {
	args := m.Called(kube, ips, ctx)
	return args.Error(0)
}

func TestProvider_CreateKube(t *testing.T) {
	c := &core.Core{
		DB:  new(fake_core.DB),
		Log: logrus.New(),
	}
	clientMock := new(DOClientMock)
	provMock := new(ProvisionMock)

	p := Provider{
		DOClient:  clientMock,
		Core:      c,
		Provision: provMock,
	}

	clientMock.On("NewDroplet", mock.Anything, mock.Anything, mock.Anything).Return(godo.Droplet{
		Networks: &godo.Networks{
			V4: []godo.NetworkV4{
				{
					Gateway:   "10.1.1.1",
					IPAddress: "127.0.0.1",
					Netmask:   "255.255.225.0",
					Type:      "public",
				}},
		},
	}, nil)
	provMock.On("CreateMaster",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(nil)

	m := &model.Kube{
	}

	action := &core.Action{
		Core: c,
	}

	err := p.CreateKube(m, action)

	clientMock.AssertExpectations(t)
	provMock.AssertExpectations(t)
	require.NoError(t, err, "error while calling provider.createKube")
}
