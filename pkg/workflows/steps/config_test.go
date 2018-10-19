package steps

import (
	"encoding/json"
	"testing"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/node"
)

func TestMarshalConfig(t *testing.T) {
	nodes := []*node.Node{{ID: "1"}, {ID: "2"}}
	masterMap := make(map[string]*node.Node)

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
	discoveryUrl := "https://etcd.io"
	expectedMasterCount := 3
	expectedNodeCount := 5

	p := profile.Profile{
		MasterProfiles: make([]profile.NodeProfile, expectedMasterCount),
		NodesProfiles:  make([]profile.NodeProfile, expectedNodeCount),
	}

	cfg := NewConfig(clusterName, discoveryUrl, cloudAccountName, p)

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
	n := &node.Node{
		Role: node.RoleMaster,
	}

	cfg := &Config{
		Masters: Map{
			internal: make(map[string]*node.Node),
		},
	}

	cfg.AddMaster(n)

	if len(cfg.Masters.internal) != 1 {
		t.Errorf("Wrong masters count expected %d actual %d",
			1, len(cfg.Masters.internal))
	}
}

func TestAddNode(t *testing.T) {
	n := &node.Node{
		Role: node.RoleNode,
	}

	cfg := &Config{
		Nodes: Map{
			internal: make(map[string]*node.Node),
		},
	}

	cfg.AddNode(n)

	if len(cfg.Nodes.internal) != 1 {
		t.Errorf("Wrong node count expected %d actual %d",
			1, len(cfg.Nodes.internal))
	}
}

func TestConfigGetMaster(t *testing.T) {
	n := &node.Node{
		Name:  "master-1",
		State: node.StateActive,
		Role:  node.RoleMaster,
	}

	testCases := []struct {
		cfg          *Config
		expectedNode *node.Node
	}{
		{
			cfg: &Config{
				Masters: Map{
					internal: map[string]*node.Node{
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
	n1, n2 := &node.Node{
		Name:  "node-1",
		State: node.StateError,
		Role:  node.RoleNode,
	}, &node.Node{
		Name:  "node-2",
		State: node.StateActive,
		Role:  node.RoleNode,
	}

	testCases := []struct {
		cfg          *Config
		expectedNode *node.Node
	}{
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*node.Node{
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
					internal: map[string]*node.Node{
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
					internal: map[string]*node.Node{
						"node-1": {},
					},
				},
			},
			expectedCount: 1,
		},
		{
			cfg: &Config{
				Masters: Map{
					internal: map[string]*node.Node{},
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
					internal: map[string]*node.Node{
						"node-1": {},
					},
				},
			},
			expectedCount: 1,
		},
		{
			cfg: &Config{
				Nodes: Map{
					internal: map[string]*node.Node{},
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
