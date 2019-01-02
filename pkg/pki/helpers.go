package pki

import (
	"crypto/rsa"
	"crypto/x509"
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
