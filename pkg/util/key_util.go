package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const keySize = 4096

// Create bootstrap key pair and save to config ssh section
func BootstrapKeys(config *steps.Config) error {
	private, public, err := generateKeyPair(keySize)

	if err != nil {
		return err
	}

	config.Kube.SSHConfig.BootstrapPrivateKey = private
	config.Kube.SSHConfig.BootstrapPublicKey = public

	return nil
}

func generateKeyPair(size int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		return "", "", err
	}

	privateKeyPem := encodePrivateKeyToPEM(privateKey)
	publicKey, err := generatePublicKey(&privateKey.PublicKey)

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

func generatePublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)

	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}
