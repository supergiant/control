package steps

import (
	"encoding/json"
	"testing"

	"github.com/supergiant/supergiant/pkg/node"
)

func TestMarshalConfig(t *testing.T) {
	nodes := []*node.Node{{Id: "1"}, {Id: "2"}}
	masterMap := make(map[string]*node.Node)

	for _, n := range nodes {
		masterMap[n.Id] = n
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
		_, ok := cfg2.Masters.internal[n.Id]

		if !ok {
			t.Errorf("Node id %s not found in master map %v", n.Id, cfg2.Masters.internal)
			return
		}
	}
}
