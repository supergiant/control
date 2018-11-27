package ssh

import (
	"net"
	"time"

	"context"
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

	if config.Host == "" {
		return nil, ErrHostNotSpecified
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
		BannerCallback: func(message string) error {
			logrus.Debug(message)
			return nil
		},
	}, nil
}

func connectionWithBackOff(ctx context.Context, host, port string, config *ssh.ClientConfig, timeout time.Duration, attemptCount int) (*ssh.Client, error) {
	var (
		counter = 0
		c       *ssh.Client
		err     error
	)

	for counter < attemptCount {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)

			if err != nil {
				logrus.Warnf("connect to %s failed, try again in %v seconds, reason: %v",
					fmt.Sprintf("%s:%s", host, port),
					timeout, err)
				time.Sleep(timeout)
				timeout = timeout * 2
			} else {
				return c, err
			}
			counter += 1
		}
	}

	return nil, err
}
