package task

import (
	"context"
	"os"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

func RunRemoteScript(ctx context.Context, script, user, host, cert string, timeoutSec int) error {
	key := []byte(cert)

	cfg := &ssh.Config{
		User:    user,
		Host:    host,
		Timeout: timeoutSec,
		Port:    "22",
		Key:     key,
	}

	run, err := ssh.NewRunner(cfg)
	if err != nil {
		return err
	}

	cmd := runner.NewCommand(ctx, script, os.Stdout, os.Stderr)
	return run.Run(cmd)
}
