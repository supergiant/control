package pki

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
)

// newCertificateAuthority creates new certificate and private key for the certificate authority
func newCertificateAuthority() (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := certutil.NewPrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "create private key")
	}

	config := certutil.Config{
		CommonName: "kubernetes",
	}
	cert, err := certutil.NewSelfSignedCACert(config, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create self-signed certificate")
	}

	return cert, key, nil
}

// NewCertAndKey creates signed certificate and key for the provided CA.
func NewCertAndKey(caCert *x509.Certificate, caKey *rsa.PrivateKey, config *certutil.Config) (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := certutil.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("create private key: %s", err)
	}

	cert, err := certutil.NewSignedCert(*config, key, caCert, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("create signed cert: %s", err)
	}

	return cert, key, nil
}
