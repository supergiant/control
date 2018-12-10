package provisioner

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/testutils"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
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
	m.data[k.ID] = k
	return m.createErr
}

func (m *mockKubeService) Get(ctx context.Context, kname string) (*model.Kube, error) {
	return m.data[kname], m.getError
}

func TestNewProvisioner(t *testing.T) {
	storage := &testutils.MockStorage{}
	service := &mockKubeService{}
	interval := time.Second * 1

	p := NewProvisioner(storage, service, interval)

	if p.repository != storage {
		t.Errorf("Wrong repository expected %v actual %v",
			storage, p.repository)
	}

	if p.kubeService != service {
		t.Errorf("Wrong kube service value expected %v actual %v",
			service, p.kubeService)
	}

	if p.cancelMap == nil {
		t.Errorf("Cancel map must not be nil")
	}

	if p.provisionMap == nil {
		t.Errorf("Provision map must not be nil")
	}
}

func TestProvisionCluster(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything,
		mock.Anything).Return(nil)

	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	svc := &mockKubeService{
		data: make(map[string]*model.Kube),
	}

	provisioner := TaskProvisioner{
		svc,
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				ProvisionMaster: "test_master",
				ProvisionNode:   "test_node",
				PreProvision:    "",
			},
		},
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_master", []steps.Step{})
	workflows.RegisterWorkFlow("test_node", []steps.Step{})
	workflows.RegisterWorkFlow(workflows.Cluster, []steps.Step{})

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

	if len(provisioner.cancelMap) != 1 {
		t.Errorf("Unexpected size of cancel map expected %d actual %d",
			1, len(provisioner.cancelMap))
	}

	if k := svc.data[cfg.ClusterID]; k == nil {
		t.Errorf("Kube %s not found", k.ID)

		if len(k.Tasks) != len(p.MasterProfiles)+len(p.NodesProfiles)+1 {
			t.Errorf("Wrong count of tasks in kube expected %d actual %d",
				len(p.MasterProfiles)+len(p.NodesProfiles)+1, len(k.Tasks))
		}
	}
}

func TestProvisionNodes(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	repository.On("Get", mock.Anything, mock.Anything,
		mock.Anything).Return()
	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	kubeID := "1234"

	k := &model.Kube{
		ID:       kubeID,
		Provider: clouds.DigitalOcean,
		Masters: map[string]*node.Node{
			"1": {
				ID:        "1",
				PrivateIp: "10.0.0.1",
				PublicIp:  "10.20.30.40",
				State:     node.StateActive,
				Region:    "fra1",
				Size:      "s-2vcpu-4gb",
			},
		},
		BootstrapPublicKey:  []byte(""),
		BootstrapPrivateKey: []byte(""),
		SshPublicKey:        []byte(""),
		CloudSpec:           make(map[string]string),
	}

	provisioner := TaskProvisioner{
		&mockKubeService{
			data: map[string]*model.Kube{
				k.ID: k,
			},
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
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_node", []steps.Step{})

	nodeProfile := profile.NodeProfile{
		"size":  "s-2vcpu-4gb",
		"image": "ubuntu-18-04-x64",
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
	config.ClusterID = k.ID

	_, err := provisioner.ProvisionNodes(context.Background(),
		[]profile.NodeProfile{nodeProfile}, k, config)

	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}

	if len(provisioner.cancelMap) != 1 {
		t.Errorf("Unexpected size of cancel map expected %d actual %d",
			1, len(provisioner.cancelMap))
	}

}

func TestRestartProvisionCluster(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything,
		mock.Anything).Return(nil)
	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return([]byte(`{"id": "task_id", "type": "preprovision"}`),
		nil)

	bc := &bufferCloser{
		bytes.Buffer{},
		nil,
	}

	svc := &mockKubeService{
		data: map[string]*model.Kube{
			"kubeID": {
				ID: "kubeID",
			},
		},
	}

	provisioner := TaskProvisioner{
		svc,
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				ProvisionMaster: workflows.AWSMaster,
				ProvisionNode:   workflows.AWSNode,
				PreProvision:    workflows.AWSPreProvision,
			},
		},
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.AWSMaster, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.AWSNode, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.Cluster, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.AWSPreProvision, []steps.Step{})

	p := &profile.Profile{
		Provider: clouds.AWS,
		MasterProfiles: []profile.NodeProfile{
			{},
		},
		NodesProfiles: []profile.NodeProfile{
			{},
			{},
		},
	}

	taskMap := map[string][]string{
		workflows.PreProvisionTask: {
			"task_id",
		},
	}
	cfg := steps.NewConfig("kube_name",
		"", "", *p)
	cfg.ClusterID = "kubeID"

	err := provisioner.
		RestartClusterProvisioning(context.Background(),
			p, cfg, taskMap)

	time.Sleep(time.Second * 1)
	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}

	if _, ok := provisioner.cancelMap["kubeID"]; !ok {
		t.Errorf("cancal map for kube must not be empty")
	}
}

func TestDeserializeTasks(t *testing.T) {
	repository := &testutils.MockStorage{}

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "1234", "type": "preprovision"}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(
			`{"id": "4567", "type": "master"}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "9876", "type": "node"}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "abcd", "type": "cluster"}`),
		nil)

	provisioner := TaskProvisioner{
		repository: repository,
	}

	taskIdMap := map[string][]string{
		workflows.PreProvisionTask: {"1234"},
		workflows.MasterTask:       {"4567"},
		workflows.NodeTask:         {"9876"},
		workflows.ClusterTask:      {"abcd"},
	}

	taskMap, err := provisioner.deserializeClusterTasks(context.Background(),
		taskIdMap)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if len(taskMap) != len(taskIdMap) {
		t.Errorf("Wrong task count expected %d actual %d",
			len(taskIdMap), len(taskMap))
	}
}

func TestDeserializeTasksError(t *testing.T) {
	repository := &testutils.MockStorage{}

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		nil,
		sgerrors.ErrNotFound)

	provisioner := TaskProvisioner{
		repository: repository,
	}

	taskIdMap := map[string][]string{
		workflows.PreProvisionTask: {"1234"},
	}

	taskMap, err := provisioner.deserializeClusterTasks(context.Background(),
		taskIdMap)

	if errors.Cause(err) != sgerrors.ErrNotFound {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if taskMap != nil {
		t.Error("Task map must be nil")
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
				ID:      "1234",
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
				ID:      "1234",
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
				testCase.kube.ID: testCase.kube,
			},
		}

		p := &TaskProvisioner{
			kubeService: svc,
		}
		cfg := steps.NewConfig(
			"test",
			"",
			"test",
			profile.Profile{})
		cfg.ClusterID = testCase.kube.ID
		logrus.Println(testCase.kube.ID)

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

func TestTaskProvisioner_Cancel(t *testing.T) {
	clusterID := "1234"
	called := false
	f := func() {
		called = true
	}

	tp := &TaskProvisioner{
		cancelMap: map[string]func(){
			clusterID: f,
		},
	}

	tp.Cancel(clusterID)

	if !called {
		t.Errorf("Cancel function was not called")
	}
}

func TestTaskProvisioner_CancelNotFound(t *testing.T) {
	clusterID := "1234"
	called := false
	f := func() {
		called = true
	}

	tp := &TaskProvisioner{
		cancelMap: map[string]func(){
			clusterID: f,
		},
	}

	err := tp.Cancel("not_found")

	if called {
		t.Errorf("Cancel function must not been called")
	}

	if err != sgerrors.ErrNotFound {
		t.Errorf("Unexpected error value %v", err)
	}
}

func TestBuildInitialCluster(t *testing.T) {
	service := &mockKubeService{
		data: make(map[string]*model.Kube),
	}
	clusterID := "clusterID"
	tp := &TaskProvisioner{
		kubeService: service,
	}

	taskIds := map[string][]string{
		workflows.MasterTask: {"1234", "5678", "abcd"},
	}

	tp.buildInitialCluster(context.Background(), &profile.Profile{}, nil, nil, &steps.Config{
		ClusterID: clusterID,
	}, taskIds)

	if k := service.data[clusterID]; k == nil {
		t.Errorf("Cluster %s not found", clusterID)
	} else {
		if len(k.Tasks) != len(taskIds) {
			t.Errorf("Wrong number of tasks in cluster "+
				"expected %d actual %d", len(taskIds), len(k.Tasks))
		}
	}
}
