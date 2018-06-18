package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/supergiant/supergiant/pkg/runner/command"
	"golang.org/x/crypto/ssh"
)

// Config is a set of params needed to create valid ssh.ClientConfig
type Config struct {
	Host    string
	Port    int
	User    string
	Timeout int

	Key []byte

	SshClientConfig *ssh.ClientConfig
}

// Runner is implementation of runner interface for ssh
type Runner struct {
	*Config
	out io.Writer
	err io.Writer

	client *ssh.Client
}

// NewRunner creates ssh runner object. It requires two io.Writer
// to send output of ssh session and config for ssh client.
func NewRunner(outStream, errStream io.Writer, config *Config) (*Runner, error) {
	if sshConfig, err := getSshConfig(config); err != nil {
		config.SshClientConfig = sshConfig
	}

	r := &Runner{
		config,
		outStream,
		errStream,
		nil,
	}

	err := r.connect()

	if err != nil {
		return nil, err
	}

	return r, nil
}

// Connect to server with ssh
func (r *Runner) connect() error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port), r.SshClientConfig)
	if err == nil {
		r.client = client
		return nil
	}

	return err
}

//TODO(stgleb): Add  more context like env variables?
// Run executes a single command on ssh session.
func (r *Runner) Run(c *command.Command) (err error) {
	if r.client == nil {
		return errors.New("not connected")
	}

	cmd := strings.TrimSpace(c.FullCommand())
	if cmd == "" {
		return nil
	}

	session, err := r.client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdout = io.MultiWriter(r.out, c.Out)
	session.Stderr = io.MultiWriter(r.out, c.Out)

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
