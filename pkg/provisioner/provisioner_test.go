package provisioner

import (
	"context"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/testutils"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockKubeService struct {
	getError  error
	createErr error
	lock      sync.RWMutex
	data      map[string]model.Kube
}

func (m *mockKubeService) Create(ctx context.Context, k *model.Kube) error {
	if k == nil {
		return nil
	}
	m.lock.Lock()
	m.data[k.ID] = *k
	m.lock.Unlock()
	return m.createErr
}

func (m *mockKubeService) Get(ctx context.Context, kname string) (*model.Kube, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	k := m.data[kname]
	return &k, m.getError
}

type mockStep struct {
}

func (m *mockStep) Run(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (m *mockStep) Name() string {
	return ""
}

func (m *mockStep) Description() string {
	return ""
}

func (m *mockStep) Depends() []string {
	return nil
}

func (m *mockStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func TestNewProvisioner(t *testing.T) {
	storage := &testutils.MockStorage{}
	service := &mockKubeService{}
	interval := time.Second * 1

	p := NewProvisioner(storage, service, interval, "")

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
}

func TestProvisionCluster(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything,
		mock.Anything).Return(nil)

	bc := &bufferCloser{
		ioutil.Discard,
		nil,
	}

	svc := &mockKubeService{
		data: make(map[string]model.Kube),
	}

	provisioner := TaskProvisioner{
		svc,
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.ProvisionMaster, []steps.Step{
		&mockStep{},
	})
	workflows.RegisterWorkFlow(workflows.ProvisionNode, []steps.Step{
		&mockStep{},
	})
	workflows.RegisterWorkFlow(workflows.PostProvision, []steps.Step{
		&mockStep{},
	})
	workflows.RegisterWorkFlow(workflows.DigitalOceanInfra, []steps.Step{
		&mockStep{},
	})

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

	cfg, err := steps.NewConfig("", "", *p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	_, err = provisioner.ProvisionCluster(ctx, p, cfg)
	time.Sleep(time.Millisecond * 10)

	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}

	// Cancel context to shut down cluster state monitoring
	cancel()
	if k, err := svc.Get(context.Background(), cfg.Kube.ID); err == nil && k == nil {
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
		ioutil.Discard,
		nil,
	}

	kubeID := "1234"

	k := &model.Kube{
		ID:       kubeID,
		Provider: clouds.DigitalOcean,
		Masters: map[string]*model.Machine{
			"1": {
				ID:        "1",
				PrivateIp: "10.0.0.1",
				PublicIp:  "10.20.30.40",
				State:     model.MachineStateActive,
				Region:    "fra1",
				Size:      "s-2vcpu-4gb",
			},
		},
		CloudSpec: make(map[string]string),
	}

	provisioner := TaskProvisioner{
		&mockKubeService{
			data: map[string]model.Kube{
				k.ID: *k,
			},
		},
		repository,
		func(string) (io.WriteCloser, error) {
			return bc, nil
		},
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.ProvisionNode, []steps.Step{
		&mockStep{},
	})

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

	config, err := steps.NewConfig(k.Name, k.AccountName, kubeProfile)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	config.Kube.ID = k.ID

	_, err = provisioner.ProvisionNodes(context.Background(),
		[]profile.NodeProfile{nodeProfile}, k, config)

	time.Sleep(time.Millisecond * 10)
	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}

	if len(provisioner.cancelMap) != 1 {
		t.Errorf("Unexpected size of cancel map expected %d actual %d",
			1, len(provisioner.cancelMap))
	}

}

func TestRestartProvisionClusterSuccess(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything,
		mock.Anything).Return(nil)
	repository.On("Get", mock.Anything,
		mock.Anything,
		mock.Anything).Return([]byte(`{"id": "task_id", 
		"type": "PreProvision", "stepsStatuses":[{"status": "error"}], "config": {}}`),
		nil).Once()

	repository.On("Get", mock.Anything,
		mock.Anything,
		mock.Anything).Return([]byte(`{"id": "task_id", 
		"type": "ProvisionMaster", "stepsStatuses":[{"status": "error"}] "config": {}}`),
		nil).Once()

	bc := &bufferCloser{
		ioutil.Discard,
		nil,
	}

	svc := &mockKubeService{
		data: map[string]model.Kube{
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
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.DigitalOceanInfra, []steps.Step{
		&mockStep{},
	})

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
	cfg, err := steps.NewConfig("kube_name", "", *p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	cfg.Kube.ID = "kubeID"

	err = provisioner.
		RestartClusterProvisioning(context.Background(),
			p, cfg, taskMap)

	time.Sleep(time.Millisecond * 10)

	if err != nil {
		t.Errorf("Unexpected error %v while provisionCluster", err)
	}
}

func TestRestartProvisionClusterError(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", mock.Anything,
		mock.Anything, mock.Anything,
		mock.Anything).Return(nil)
	repository.On("Get", mock.Anything,
		mock.Anything,
		mock.Anything).Return([]byte(`}`),
		nil)

	bc := &bufferCloser{
		ioutil.Discard,
		nil,
	}

	svc := &mockKubeService{
		data: map[string]model.Kube{
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
		NewRateLimiter(time.Nanosecond * 1),
		make(map[string]func()),
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.ProvisionMaster, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.ProvisionNode, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.PostProvision, []steps.Step{})
	workflows.RegisterWorkFlow(workflows.DigitalOceanInfra, []steps.Step{
		&mockStep{},
	})

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
	cfg, err := steps.NewConfig("kube_name",
		"", *p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	cfg.Kube.ID = "kubeID"

	err = provisioner.
		RestartClusterProvisioning(context.Background(),
			p, cfg, taskMap)

	time.Sleep(time.Millisecond * 10)
	if err == nil {
		t.Errorf("Error must not be nil")
	}
}

func TestDeserializeTasks(t *testing.T) {
	repository := &testutils.MockStorage{}

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "1234", "type": "preprovision", "config": {"provider": "aws"}}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(
			`{"id": "4567", "type": "master", "config": {"provider": "aws"}"}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "9876", "type": "node", "config": {"provider": "aws"}}`),
		nil)

	repository.On("Get", mock.Anything,
		mock.Anything, mock.Anything).Return(
		[]byte(`{"id": "abcd", "type": "cluster", "config": {"provider": "aws"}}`),
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

	taskMap, err := provisioner.deserializeClusterTasks(context.Background(), &steps.Config{
		Provider: clouds.AWS,
	},
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

	taskMap, err := provisioner.deserializeClusterTasks(context.Background(), &steps.Config{},
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
		name                 string
		nodes                []model.Machine
		states               []model.KubeState
		kube                 model.Kube
		expectedMastersNum   int
		expectedClusterState model.KubeState
	}{
		{
			name: "state_operational",
			nodes: []model.Machine{
				{
					Name:  "test",
					Role:  model.RoleMaster,
					State: model.MachineStatePlanned,
				},
				{
					Name:  "test",
					Role:  model.RoleMaster,
					State: model.MachineStateBuilding,
				},
				{
					Name:  "test",
					Role:  model.RoleMaster,
					State: model.MachineStateProvisioning,
				},
				{
					Name:  "test",
					Role:  model.RoleMaster,
					State: model.MachineStateActive,
				},
			},
			states: []model.KubeState{
				model.StateProvisioning,
				model.StateOperational,
			},
			kube: model.Kube{
				ID:      "1234",
				Name:    "test",
				Masters: make(map[string]*model.Machine),
				Nodes:   make(map[string]*model.Machine),
			},
			expectedMastersNum:   1,
			expectedClusterState: model.StateOperational,
		},
		{
			name: "state_failed",
			nodes: []model.Machine{
				{
					Name:  "test1",
					Role:  model.RoleMaster,
					State: model.MachineStateProvisioning,
				},
				{
					Name:  "test2",
					Role:  model.RoleMaster,
					State: model.MachineStateError,
				},
				{
					Name:  "test1",
					Role:  model.RoleMaster,
					State: model.MachineStateProvisioning,
				},
				{
					Name:  "test2",
					Role:  model.RoleMaster,
					State: model.MachineStateActive,
				},
			},
			states: []model.KubeState{
				model.StateProvisioning,
				model.StateFailed,
			},
			kube: model.Kube{
				ID:      "1234",
				Name:    "test",
				Masters: make(map[string]*model.Machine),
				Nodes:   make(map[string]*model.Machine),
			},
			expectedMastersNum:   2,
			expectedClusterState: model.StateFailed,
		},
		{
			name: "state_provisioning",
			kube: model.Kube{
				Name:  "test",
				State: model.StateProvisioning,
			},
			expectedClusterState: model.StateProvisioning,
		},
	}

	for _, testCase := range testCases {
		svc := &mockKubeService{
			data: map[string]model.Kube{
				testCase.kube.ID: testCase.kube,
			},
		}

		p := &TaskProvisioner{
			kubeService: svc,
		}
		cfg, err := steps.NewConfig(
			"test",
			"test",
			profile.Profile{})

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		cfg.Kube.ID = testCase.kube.ID
		logrus.Println(testCase.kube.ID)

		ctx, cancel := context.WithCancel(context.Background())
		go p.monitorClusterState(ctx, cfg.Kube.ID, cfg.NodeChan(),
			cfg.KubeStateChan(), cfg.ConfigChan())

		for _, n := range testCase.nodes {
			cfg.NodeChan() <- n
		}

		for _, state := range testCase.states {
			cfg.KubeStateChan() <- state
		}

		time.Sleep(time.Millisecond * 1)
		cancel()

		k, _ := svc.Get(context.Background(), testCase.kube.ID)

		if len(k.Masters) != testCase.expectedMastersNum {
			t.Fatalf("TC#%s: nodes number: expected=%d got=%d", testCase.name, testCase.expectedMastersNum, len(k.Masters))
		}

		if k.State != testCase.expectedClusterState {
			t.Fatalf("TC#%s: cluster state: expected=%s got=%s",
				testCase.name, testCase.expectedClusterState, testCase.kube.State)
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
		data: make(map[string]model.Kube),
	}
	clusterID := "clusterID"
	tp := &TaskProvisioner{
		kubeService: service,
	}

	taskIds := map[string][]string{
		workflows.MasterTask: {"1234", "5678", "abcd"},
	}

	tp.buildInitialCluster(context.Background(), &profile.Profile{}, nil, nil, &steps.Config{
		Kube: model.Kube{
			ID: clusterID,
		},
	}, taskIds)

	if k, _ := service.Get(context.Background(), clusterID); k == nil {
		t.Errorf("Cluster %s not found", clusterID)
	} else {
		if len(k.Tasks) != len(taskIds) {
			t.Errorf("Wrong number of tasks in cluster "+
				"expected %d actual %d", len(taskIds), len(k.Tasks))
		}
	}
}
