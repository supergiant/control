package ssh

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
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

func NewRunner(config *Config) (*Runner, error) {
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
	if r.Host == "" {
		r.Host = "localhost"
	}
	if r.User == "" {
		r.User = "root"
	}
	if r.Port == 0 {
		r.Port = 22
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port), r.SshClientConfig)
	if err == nil {
		r.client = client
		return nil
	}

	return err
}

// Exec single command on ssh session
func (r *Runner) Exec(cmd string) (err error) {
	if r.client == nil {
		return errors.New("not connected")
	}

	session, err := r.client.NewSession()
	defer session.Close()

	if err != nil {
		return err
	}

	session.Stdout = r.out
	session.Stderr = r.err

	err = session.Run(cmd)

	if err != nil {
		return err
	}

	// We can close session multiple times
	return session.Close()
}
