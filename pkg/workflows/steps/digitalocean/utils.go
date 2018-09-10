package digitalocean

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/digitalocean/godo"
	"strings"
	"encoding/base64"
	"github.com/pkg/errors"
)

// Returns private ip
func getPrivateIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "private" {
			return network.IPAddress
		}
	}

	return ""
}

// Returns public ip
func getPublicIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "public" {
			return network.IPAddress
		}
	}

	return ""
}

func generateKeyPair(size int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		return "", "", err
	}

	privateKeyPem := encodePrivateKeyToPEM(privateKey)
	publicKey, err := generatePublicKey(privateKey)

	if err != nil {
		return "", "", err
	}

	return string(privateKeyPem), string(publicKey), nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func generatePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)

	privateKey.Public()
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

func fingerprint(key string) (string, error) {
	parts := strings.Fields(string(key))

	k, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.Wrap(err, "fingerprint decode string")
	}

	fp := md5.Sum(k)
	buffer := &bytes.Buffer{}

	for i, b := range fp {
		fmt.Fprintf(buffer, "%02x", b)

		if i < len(fp)-1 {
			fmt.Fprint(buffer, ":")
		}
	}

	return buffer.String(), nil
}
