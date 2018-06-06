package ssh

import (
	"io/ioutil"
	"os/user"
	"time"

	"golang.org/x/crypto/ssh"
)

func getKey(keyFilePath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func getSshConfig(config *Config) (*ssh.ClientConfig, error) {
	key, err := getKey(config.KeyFilePath)

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
