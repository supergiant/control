package digitalocean

import (
	"testing"
	"github.com/supergiant/supergiant/pkg/runner/command"
	"github.com/supergiant/supergiant/pkg/runner"
)

type fakeRunner struct{
	commands []string
}

func (f *fakeRunner) Run(command *command.Command) error {
	f.commands = append(f.commands, command.FullCommand())
	return nil
}

func TestJob_ProvisionNode(t *testing.T) {
	testCases := []struct{
		r runner.Runner
		config string
		kubeletService string
		kubeletScript  string
		proxyScript    string
	}{
		{

		},
		{

		},
	}

	for _, testCase := range testCases {
		
	}
}
