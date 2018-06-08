package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/ayufan/gitlab-ci-multi-runner/helpers"

	"github.com/supergiant/supergiant/pkg/runner/command"
)

type Config struct {
	Host    string
	Port    int
	User    string
	Timeout int

	KeyFilePath string

	SshClientConfig *ssh.ClientConfig
}

type Runner struct {
	*Config
	out io.Writer
	err io.Writer

	client *ssh.Client
}

// NewRunner creates ssh runner object. It requires two io.Writer
// to send output of ssh session and config for ssh client.
func NewRunner(out, err io.Writer, config *Config) (*Runner, error) {
	if sshConfig, err := getSshConfig(config); err != nil {
		config.SshClientConfig = sshConfig
	}

	return &Runner{
		config,
		os.Stdout,
		os.Stderr,
		nil,
	}, nil
}

// Connect to server with ssh
func (r *Runner) Connect() error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port), r.SshClientConfig)
	if err == nil {
		r.client = client
		return nil
	}

	return err
}

//TODO(stgleb): Add  more context like env variables?
// Exec single command on ssh session
func (r *Runner) Run(c command.Command) (err error) {
	if r.client == nil {
		return errors.New("not connected")
	}

	cmd := helpers.ShellEscape(c.FullCommand())
	session, err := r.client.NewSession()
	defer session.Close()

	if err != nil {
		return err
	}

	session.Stdout = r.out
	session.Stderr = r.err

	err = session.Start(cmd)

	if err != nil {
		return err
	}

	waitCh := make(chan error)
	go func() {
		waitCh <- session.Wait()
	}()

	select {
	case <-c.Ctx.Done():

		if c.Ctx.Err() == context.Canceled {
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
