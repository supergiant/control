package steps

import (
	"encoding/json"
	"testing"

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
