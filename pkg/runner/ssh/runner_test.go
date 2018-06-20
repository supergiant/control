package ssh

import (
	"context"
	"io/ioutil"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/supergiant/supergiant/pkg/runner"
)

func TestRunner_Run(t *testing.T) {
	testCases := []struct {
		runner      runner.Runner
		command     *runner.Command
		expectedErr error
	}{
		{
			runner: &Runner{
				client: &ssh.Client{},
			},
			command:     runner.NewCommand(context.Background(), "", ioutil.Discard, ioutil.Discard),
			expectedErr: ErrEmptyScript,
		},
		{
			runner:      &Runner{},
			command:     runner.NewCommand(context.Background(), "", ioutil.Discard, ioutil.Discard),
			expectedErr: ErrNotConnected,
		},
	}

	for _, testCase := range testCases {
		err := testCase.runner.Run(testCase.command)

		if err != testCase.expectedErr {
			t.Errorf("Wrong error expected %v actual %v", testCase.expectedErr, err)
		}
	}
}
