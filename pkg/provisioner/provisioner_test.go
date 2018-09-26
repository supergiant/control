package provisioner

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

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

type mockKubeService struct {
	getError  error
	createErr error
	data      map[string]*model.Kube
}

func (m *mockKubeService) Create(ctx context.Context, k *model.Kube) error {
	m.data[k.Name] = k
	return m.createErr
}

func (m *mockKubeService) Get(ctx context.Context, kname string) (*model.Kube, error) {
	return m.data[kname], m.getError
}

func TestProvisionCluster(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	provisioner := TaskProvisioner{
		&mockKubeService{
			data: make(map[string]*model.Kube),
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
		&mockKubeService{
			data: make(map[string]*model.Kube),
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
				State:     node.StateActive,
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

func TestMonitorCluster(t *testing.T) {
	testCases := []struct {
		nodes                []node.Node
		states               []model.KubeState
		kube                 *model.Kube
		expectedNodeCount    int
		expectedClusterState model.KubeState
	}{
		{
			nodes: []node.Node{
				{
					Name:  "test",
					Role:  node.RoleMaster,
					State: node.StatePlanned,
				},
				{
					Name:  "test",
					Role:  node.RoleMaster,
					State: node.StateBuilding,
				},
				{
					Name:  "test",
					Role:  node.RoleMaster,
					State: node.StateProvisioning,
				},
				{
					Name:  "test",
					Role:  node.RoleMaster,
					State: node.StateActive,
				},
			},
			states: []model.KubeState{
				model.StateProvisioning,
				model.StateOperational,
			},
			kube: &model.Kube{
				Name:    "test",
				Masters: make(map[string]*node.Node),
				Nodes:   make(map[string]*node.Node),
			},
			expectedNodeCount:    1,
			expectedClusterState: model.StateOperational,
		},
		{
			nodes: []node.Node{
				{
					Name:  "test1",
					Role:  node.RoleMaster,
					State: node.StateProvisioning,
				},
				{
					Name:  "test2",
					Role:  node.RoleMaster,
					State: node.StateError,
				},
				{
					Name:  "test1",
					Role:  node.RoleMaster,
					State: node.StateProvisioning,
				},
				{
					Name:  "test2",
					Role:  node.RoleMaster,
					State: node.StateActive,
				},
			},
			states: []model.KubeState{
				model.StateProvisioning,
				model.StateFailed,
			},
			kube: &model.Kube{
				Name:    "test",
				Masters: make(map[string]*node.Node),
				Nodes:   make(map[string]*node.Node),
			},
			expectedNodeCount:    2,
			expectedClusterState: model.StateFailed,
		},
		{
			kube: &model.Kube{
				Name:  "test",
				State: model.StateProvisioning,
			},
			expectedNodeCount:    0,
			expectedClusterState: model.StateProvisioning,
		},
	}

	for _, testCase := range testCases {
		svc := &mockKubeService{
			data: map[string]*model.Kube{
				testCase.kube.Name: testCase.kube,
			},
		}

		p := &TaskProvisioner{
			kubeService: svc,
		}
		cfg := steps.NewConfig("test",
			"",
			"test",
			profile.Profile{})

		ctx, cancel := context.WithCancel(context.Background())
		go p.monitorClusterState(ctx, cfg)

		for _, n := range testCase.nodes {
			cfg.NodeChan() <- n
		}

		for _, state := range testCase.states {
			cfg.KubeStateChan() <- state
		}

		time.Sleep(time.Millisecond * 1)
		cancel()

		if len(testCase.kube.Masters)+len(testCase.kube.Nodes) != testCase.expectedNodeCount {
			t.Errorf("Wrong node count in the end of provisioning")
		}

		if testCase.kube.State != testCase.expectedClusterState {
			t.Errorf("Wrong cluster state in the end of provisioning")
		}
	}
}
