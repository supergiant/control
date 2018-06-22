package ssh

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"

	"github.com/supergiant/supergiant/pkg/runner"
)

const (
	DefaultPort = "22"
)

// Config is a set of params needed to create valid ssh.ClientConfig
type Config struct {
	Host    string
	Port    string
	User    string
	Timeout int
	Key     []byte
}

// Runner is implementation of runner interface for ssh
type Runner struct {
	host    string
	port    string
	sshConf *ssh.ClientConfig
}

// NewRunner creates ssh runner object. It requires two io.Writer
// to send output of ssh session and config for ssh client.
// TODO: Does it safe to pass Config as a pointer?
func NewRunner(config *Config) (*Runner, error) {
	if strings.TrimSpace(config.Host) == "" {
		return nil, ErrHostNotSpecified
	}
	sshConfig, err := getSshConfig(config)
	if err != nil {
		return nil, err
	}

	r := &Runner{host: config.Host, sshConf: sshConfig}
	if config.Port == "" {
		r.port = DefaultPort
	}

	return r, nil
}

//TODO(stgleb): Add  more context like env variables?
// Run executes a single command on ssh session.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
func (r *Runner) Run(cmd *runner.Command) (err error) {
	if cmd == nil || strings.TrimSpace(cmd.Script) == "" {
		return nil
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", r.host, r.port), r.sshConf)
	if err != nil {
		return errors.Wrap(err, "ssh: dial")
	}
	session, err := c.NewSession()
	if err != nil {
		return errors.Wrap(err, "ssh: new session")
	}
	defer session.Close()

	session.Stdout = cmd.Out
	session.Stderr = cmd.Err

	waitCh := make(chan error)
	go func() {
		waitCh <- session.Run(cmd.Script)
	}()

	select {
	case <-cmd.Ctx.Done():
		if cmd.Ctx.Err() == context.Canceled {
			session.Signal(ssh.SIGKILL)
			session.Close()
		}
		return <-waitCh
	case err := <-waitCh:
		return err
	}

	// We can close session multiple times
	return session.Close()
}
