package task

import (
	"context"
	"os"

	"github.com/supergiant/supergiant/pkg/runner/command"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

func RunRemoteScript(ctx context.Context, script, user, host, cert string, timeoutSec int) error {
	key := []byte(cert)

	runner, err := ssh.NewRunner(os.Stdout, os.Stderr, &ssh.Config{
		User:    user,
		Host:    host,
		Timeout: timeoutSec,
		Port:    22,
		Key:     key,
	})

	if err != nil {
		return err
	}
	cmd := command.NewCommand(ctx, os.Stdout, os.Stderr, script)
	return runner.Run(cmd)
}
