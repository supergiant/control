package provisioner

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestTaskProvisioner(t *testing.T) {
	repository := &testutils.MockStorage{}
	repository.On("Put", context.Background(), mock.Anything, mock.Anything, mock.Anything).Return(nil)

	provisioner := TaskProvisioner{
		repository,
		func(string) (io.WriteCloser, error) {
			return os.Stdout, nil
		},
		map[clouds.Name][]string{
			clouds.DigitalOcean: {"test_master", "test_node"},
		},
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_master", []steps.Step{})
	workflows.RegisterWorkFlow("test_node", []steps.Step{})

	p := &profile.Profile{
		Provider: clouds.DigitalOcean,
		MasterProfiles: []map[string]string{
			{
				"size":  "s-1vcpu-2gb",
				"image": "ubuntu-18-04-x64",
			},
		},
		NodesProfiles: []map[string]string{
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
	taskMap, err := provisioner.Provision(context.Background(), p, cfg)

	if err != nil {
		t.Errorf("Unexpected error %v while provision", err)
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
