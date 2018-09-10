package provisioner

import (
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"testing"
)

func TestNodesFromProfile(t *testing.T) {
	region := "fra1"

	p := &profile.Profile{
		Provider: clouds.DigitalOcean,
		Region:   region,
		MasterProfiles: []map[string]string{
			{
				"image": "ubuntu-16-04-x64",
				"size":  "s-1vcpu-2gb",
			},
		},
		NodesProfiles: []map[string]string{
			{
				"image": "ubuntu-16-04-x64",
				"size":  "s-2vcpu-4gb",
			},
			{
				"image": "ubuntu-16-04-x64",
			},
		},
	}

	masters, nodes := nodesFromProfile(p)

	if len(masters) != len(p.MasterProfiles) {
		t.Errorf("Wrong master node count expected %d actual %d",
			len(p.MasterProfiles), len(masters))
	}

	if len(nodes) != len(p.NodesProfiles) {
		t.Errorf("Wrong master node count expected %d actual %d",
			len(p.NodesProfiles), len(nodes))
	}

	if masters[0].Size != p.MasterProfiles[0]["size"] {
		t.Errorf("Wrong master node size expected %s actual %s",
			masters[0].Size,
			p.MasterProfiles[0]["size"])
	}
}
