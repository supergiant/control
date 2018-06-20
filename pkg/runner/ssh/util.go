package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

var (
	ErrUserNotSpecified = errors.New("user not specified")
	ErrNotConnected     = errors.New("runner not connected")
	ErrEmptyScript      = errors.New("script is empty")
)

func getSshConfig(config *Config) (*ssh.ClientConfig, error) {
	key, err := ssh.ParsePrivateKey(config.Key)

	if err != nil {
		return nil, err
	}

	if config.User == "" {
		return nil, ErrUserNotSpecified
	}

	return &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: time.Duration(config.Timeout) * time.Second,
	}, nil
}
