package ssh

import (
	"time"

	"net"

	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	ErrUserNotSpecified = errors.New("user not specified")
	ErrHostNotSpecified = errors.New("host not specified")
)

func getSshConfig(config Config) (*ssh.ClientConfig, error) {
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
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			logrus.Debugf("hostname %s,addr %s key %s", hostname, remote.String(), string(key.Type()))
			return nil
		},
	}, nil
}

func connectionWithBackOff(host, port string, config *ssh.ClientConfig, timeout time.Duration, attemptCount int) (*ssh.Client, error) {
	var (
		counter = 0
		c       *ssh.Client
		err     error
	)

	for counter < attemptCount {
		c, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)

		if err != nil {
			time.Sleep(timeout)
			timeout = timeout * 2
		}
		counter += 1
	}

	return c, err
}
