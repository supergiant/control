package ssh

import (
	"os/user"
	"time"

	"golang.org/x/crypto/ssh"
)

func getKey(rawKey []byte) (ssh.Signer, error) {
	key, err := ssh.ParsePrivateKey(rawKey)

	if err != nil {
		return nil, err
	}
	return key, nil
}

func getSshConfig(config *Config) (*ssh.ClientConfig, error) {
	key, err := getKey(config.Key)

	if err != nil {
		return nil, err
	}

	if config.User == "" {
		u, _ := user.Current()
		config.User = u.Name
	}

	return &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: time.Duration(config.Timeout) * time.Second,
	}, nil
}
