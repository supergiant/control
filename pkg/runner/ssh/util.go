package ssh

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	ErrUserNotSpecified = errors.New("user not specified")
	ErrHostNotSpecified = errors.New("host not specified")
)

func getSshConfig(config *Config) (*ssh.ClientConfig, error) {
	if config.User == "" {
		return nil, ErrUserNotSpecified
	}
	key, err := ssh.ParsePrivateKey(config.Key)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: time.Duration(config.Timeout) * time.Second,
	}, nil
}
