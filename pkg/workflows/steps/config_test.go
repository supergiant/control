package steps

import (
	"encoding/json"
	"testing"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
)

func TestMarshalConfig(t *testing.T) {
	nodes := []*model.Machine{{ID: "1"}, {ID: "2"}}
	masterMap := make(map[string]*model.Machine)

	for _, n := range nodes {
		masterMap[n.ID] = n
	}

	cfg := &Config{
		Masters: Map{
			internal: masterMap,
		},
	}

	data, err := json.Marshal(cfg)

	if err != nil {
		t.Errorf("Marshall json %v", err)
	}

	cfg2 := &Config{}

	if err := json.Unmarshal(data, cfg2); err != nil {
		t.Errorf("Unmarshall json %v", err)
	}

	for _, n := range nodes {
		_, ok := cfg2.Masters.internal[n.ID]

		if !ok {
			t.Errorf("Node id %s not found in master map %v", n.ID, cfg2.Masters.internal)
			return
		}
	}
}

func TestNewConfig(t *testing.T) {
	clusterName := "testCluster"
	cloudAccountName := "cloudAccountName"
	expectedMasterCount := 3
	expectedNodeCount := 5

	p := profile.Profile{
		MasterProfiles: make([]profile.NodeProfile, expectedMasterCount),
		NodesProfiles:  make([]profile.NodeProfile, expectedNodeCount),
	}

	cfg, err := NewConfig(clusterName, cloudAccountName, p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if cfg.ClusterName != clusterName {
		t.Errorf("Wrong cluster name expected %s actual %s", clusterName, cfg.ClusterName)
	}

	if cfg.CloudAccountName != cloudAccountName {
		t.Errorf("Wrong cloud account name expected %s actual %s",
			cloudAccountName, cfg.CloudAccountName)
	}

	if cap(cfg.nodeChan) != expectedNodeCount+expectedMasterCount {
		t.Errorf("Wrong node chan capacity expected %d actual %d",
			expectedNodeCount+expectedMasterCount, len(cfg.Nodes.internal)+len(cfg.Masters.internal))
	}
}

func TestAddMaster(t *testing.T) {
	n := &model.Machine{
		Role: model.RoleMaster,
	}

	cfg := &Config{
		Masters: Map{
			internal: make(map[string]*model.Machine),
		},
	}

	cfg.AddMaster(n)

	if len(cfg.Masters.internal) != 1 {
		t.Errorf("Wrong masters count expected %d actual %d",
			1, len(cfg.Masters.internal))
	}
}

func TestAddNode(t *testing.T) {
	n := &model.Machine{
		Role: model.RoleNode,
	}

	cfg := &Config{
		Nodes: Map{
			internal: make(map[string]*model.Machine),
		},
	}

	cfg.AddNode(n)

	if len(cfg.Nodes.internal) != 1 {
		t.Errorf("Wrong node count expected %d actual %d",
			1, len(cfg.Nodes.internal))
	}
}

func TestConfigGetMaster(t *testing.T) {
	n := &model.Machine{
		Name:  "master-1",
		State: model.MachineStateActive,
		Role:  model.RoleMaster,
	}

	testCases := []struct {
		cfg          *Config
		expectedNode *model.Machine
	}{
		{
			cfg: &Config{
				Masters: Map{
					internal: map[string]*model.Machine{
						n.Name: n,
					},
				},
			},
			expectedNode: n,
		},
		{
			cfg:          &Config{},
			expectedNode: nil,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.cfg.GetMaster()

		if actual != testCase.expectedNode {
			t.Errorf("Wrong master node expected %v actual %v", testCase.expectedNode, actual)
		}
	}
}

func TestConfigGetNode(t *testing.T) {
	n1, n2 := &model.Machine{
		Name:  "node-1",
		State: model.MachineStateError,
		Role:  model.RoleNode,
	}, &model.Machine{
		Name:  "node-2",
		State: model.MachineStateActive,
		Role:  model.RoleNode,
	}

	testCases := []struct {
		cfg          *Config
		expectedNode *model.Machine
	}{
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*model.Machine{
						n1.Name: n1,
						n2.Name: n2,
					},
				},
			},
			expectedNode: n2,
		},
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*model.Machine{
						n1.Name: n1,
					},
				},
			},
			expectedNode: nil,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.cfg.GetNode()

		if actual != testCase.expectedNode {
			t.Errorf("Wrong node expected %v actual %v",
				testCase.expectedNode, actual)
		}
	}
}

func TestConfigGetMasters(t *testing.T) {
	testCases := []struct {
		cfg           *Config
		expectedCount int
	}{
		{
			cfg: &Config{
				Masters: Map{
					internal: map[string]*model.Machine{
						"node-1": {},
					},
				},
			},
			expectedCount: 1,
		},
		{
			cfg: &Config{
				Masters: Map{
					internal: map[string]*model.Machine{},
				},
			},
			expectedCount: 0,
		},
	}

	for _, testCase := range testCases {
		masters := testCase.cfg.GetMasters()

		if len(masters) != testCase.expectedCount {
			t.Errorf("Wrong amount of masters expected %d actual %d",
				testCase.expectedCount, len(masters))
		}
	}
}

func TestConfigGetNodes(t *testing.T) {
	testCases := []struct {
		cfg           *Config
		expectedCount int
	}{
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*model.Machine{
						"node-1": {},
					},
				},
			},
			expectedCount: 1,
		},
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*model.Machine{},
				},
			},
			expectedCount: 0,
		},
	}

	for _, testCase := range testCases {
		nodes := testCase.cfg.GetNodes()

		if len(nodes) != testCase.expectedCount {
			t.Errorf("Wrong amount of nodes expected %d actual %d",
				testCase.expectedCount, len(nodes))
		}
	}
}

func TestToCloudProviderOpt(t *testing.T) {
	for _, tc := range []struct {
		in  clouds.Name
		out string
	}{
		{clouds.AWS, "aws"},
		{clouds.GCE, "gce"},
		{clouds.DigitalOcean, ""},
	} {
		if toCloudProviderOpt(tc.in) != tc.out {
			t.Logf("toCloudProvider(%s) = %s", tc.in, tc.out)
		}
	}
}

func TestNewConfigFromKube(t *testing.T) {
	expectedMasterCount := 3
	expectedNodeCount := 5

	p := profile.Profile{
		MasterProfiles: make([]profile.NodeProfile, expectedMasterCount),
		NodesProfiles:  make([]profile.NodeProfile, expectedNodeCount),
	}

	k := &model.Kube{
		ID:          "ClusteID",
		Name:        "ClusterName",
		AccountName: "CloudAccount",
		CloudSpec: map[string]string{
			clouds.AwsImageID: "ImageID",
			clouds.AwsVpcID:   "VpcID",
		},
	}

	cfg, err := NewConfigFromKube(&p, k)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if cfg.ClusterName != k.Name {
		t.Errorf("Wrong cluster name expected %s actual %s", k.Name, cfg.ClusterName)
	}

	if cfg.CloudAccountName != k.AccountName {
		t.Errorf("Wrong cloud account name expected %s actual %s",
			k.AccountName, cfg.CloudAccountName)
	}

	if cfg.AWSConfig.VPCID != k.CloudSpec[clouds.AwsVpcID] {
		t.Errorf("Wrong VPCID value expected %s actual %s",
			k.CloudSpec[clouds.AwsVpcID], cfg.AWSConfig.VPCID)
	}

	if cfg.AWSConfig.ImageID != k.CloudSpec[clouds.AwsImageID] {
		t.Errorf("Wrong AWS Image ID value expected %s actual %s",
			k.CloudSpec[clouds.AwsImageID], cfg.AWSConfig.ImageID)
	}

	if cap(cfg.nodeChan) != expectedNodeCount+expectedMasterCount {
		t.Errorf("Wrong node chan capacity expected %d actual %d",
			expectedNodeCount+expectedMasterCount, len(cfg.Nodes.internal)+len(cfg.Masters.internal))
	}
}
