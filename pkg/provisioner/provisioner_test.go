package provisioner

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/clouds"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type bufferCloser struct {
	bytes.Buffer
	err error
}

func (b *bufferCloser) Close() error {
	return b.err
}

type mockKubeCreater struct {
	data map[string]string
}

func (m *mockKubeCreater) Create(ctx context.Context, k *model.Kube) error {
	return nil
}

func TestProvisionCluster(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	provisioner := TaskProvisioner{
		&mockKubeCreater{
			data: make(map[string]string),
		},
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				ProvisionMaster: "test_master",
				ProvisionNode:   "test_node",
			},
		},
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_master", []steps.Step{})
	workflows.RegisterWorkFlow("test_node", []steps.Step{})

	p := &profile.Profile{
		Provider: clouds.DigitalOcean,
		MasterProfiles: []profile.NodeProfile{
			{
				"size":  "s-1vcpu-2gb",
				"image": "ubuntu-18-04-x64",
			},
		},
		NodesProfiles: []profile.NodeProfile{
			{
				"size":  "s-2vcpu-4gb",
				"image": "ubuntu-18-04-x64",
			},
			{
				"size":  "s-2vcpu-4gb",
				"image": "ubuntu-18-04-x64",
			},
		},
	}

	cfg := steps.NewConfig("", "", "", *p)
	taskMap, err := provisioner.ProvisionCluster(context.Background(), p, cfg)

	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}

	if len(taskMap) != 3 {
		t.Errorf("Expected task map len 3 actul %d", len(taskMap))
	}

	if len(taskMap["master"])+len(taskMap["node"]) != len(p.MasterProfiles)+len(p.NodesProfiles) {
		t.Errorf("Wrong task count expected %d actual %d",
			len(p.MasterProfiles)+
				len(p.NodesProfiles),
			len(taskMap))
	}
}

func TestProvisionNodes(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", context.Background(),
		mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	provisioner := TaskProvisioner{
		&mockKubeCreater{
			data: make(map[string]string),
		},
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				ProvisionMaster: "test_master",
				ProvisionNode:   "test_node"},
		},
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_node", []steps.Step{})

	nodeProfile := profile.NodeProfile{
		"size":  "s-2vcpu-4gb",
		"image": "ubuntu-18-04-x64",
	}

	k := &model.Kube{
		Masters: map[string]*node.Node{
			"1": {
				Id:        "1",
				PrivateIp: "10.0.0.1",
				PublicIp:  "10.20.30.40",
				Active:    true,
				Region:    "fra1",
				Size:      "s-2vcpu-4gb",
			},
		},
	}

	kubeProfile := profile.Profile{
		Provider:        clouds.DigitalOcean,
		Region:          k.Region,
		Arch:            k.Arch,
		OperatingSystem: k.OperatingSystem,
		UbuntuVersion:   k.OperatingSystemVersion,
		DockerVersion:   k.DockerVersion,
		K8SVersion:      k.K8SVersion,
		HelmVersion:     k.HelmVersion,

		NetworkType:    k.Networking.Type,
		CIDR:           k.Networking.CIDR,
		FlannelVersion: k.Networking.Version,

		NodesProfiles: []profile.NodeProfile{
			nodeProfile,
		},

		RBACEnabled: k.RBACEnabled,
	}

	config := steps.NewConfig(k.Name, "", k.AccountName, kubeProfile)
	_, err := provisioner.ProvisionNodes(context.Background(), []profile.NodeProfile{nodeProfile}, k, config)

	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}
}
