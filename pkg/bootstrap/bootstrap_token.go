package bootstrap

import (
	"bufio"
	"crypto/rand"
	"fmt"
)

const (
	bootstrapTokenIDLen     = 6
	bootstrapTokenSecretLen = 16
	bootstrapTokenChars     = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func GenerateBootstrapToken() (string, error) {
	id, err := randBytes(bootstrapTokenIDLen)
	if err != nil {
		return "", err
	}

	secret, err := randBytes(bootstrapTokenSecretLen)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", id, secret), nil
}

func randBytes(length int) (string, error) {
	const maxByteValue = 252

	var (
		b     byte
		err   error
		token = make([]byte, length)
	)

	reader := bufio.NewReaderSize(rand.Reader, length*2)
	for i := range token {
		for {
			if b, err = reader.ReadByte(); err != nil {
				return "", err
			}
			if b < maxByteValue {
				break
			}
		}

		token[i] = bootstrapTokenChars[int(b)%len(bootstrapTokenChars)]
	}

	return string(token), nil
}
