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
	repository.On("Put", context.Background(), mock.Anything, mock.Anything, nil).Return(nil)

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
	workflows.RegisterWorkFlow("test_master", nil)
	workflows.RegisterWorkFlow("test_node", nil)

	kubeProfile := &profile.KubeProfile{
		MasterProfiles: []profile.NodeProfile{
			{},
		},
		NodesProfiles: []profile.NodeProfile{
			{},
			{},
		},
	}

	tasks, err := provisioner.Provision(context.Background(), kubeProfile, &steps.Config{
		Provider: clouds.DigitalOcean,
	})

	if err != nil {
		t.Errorf("Unexpected error %v while provision", err)
	}

	if len(tasks) != len(kubeProfile.MasterProfiles)+len(kubeProfile.NodesProfiles) {
		t.Errorf("Wrong task count expected %d actual %d",
			len(kubeProfile.MasterProfiles)+
				len(kubeProfile.NodesProfiles),
			len(tasks))
	}
}
