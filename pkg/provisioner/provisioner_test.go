package provisioner

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
)

type mockTokenGetter struct {
	getToken func(int) (string, error)
}

func (t *mockTokenGetter) GetToken(num int) (string, error) {
	return t.getToken(num)
}

func TestTaskProvisioner(t *testing.T) {
	tokenGetter := &mockTokenGetter{
		func(num int) (string, error) {
			return fmt.Sprintf("%s", util.RandomString(4)), nil
		},
	}

	repository := &testutils.MockStorage{}
	repository.On("Put", context.Background(), mock.Anything, mock.Anything, nil).Return(nil)

	provisioner := TaskProvisioner{
		repository,
		tokenGetter,
	}

	workflows.Init()
	workflows.RegisterWorkFlow("test_master", nil)
	workflows.RegisterWorkFlow("test_node", nil)

	// Mock digital ocean workflows with nil ones
	provisionMap[clouds.DigitalOcean] = []string{"test_master", "test_node"}

	kubeProfile := &profile.KubeProfile{
		Provider: clouds.DigitalOcean,
		MasterProfiles: []profile.NodeProfile{
			{},
		},
		NodesProfiles: []profile.NodeProfile{
			{},
			{},
		},
	}

	tasks, err := provisioner.Provision(context.Background(), kubeProfile, nil)

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
