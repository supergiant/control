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

// newCertAndKey creates new certificate and key by passing the certificate authority certificate and key
func newCertAndKey(caCert *x509.Certificate, caKey *rsa.PrivateKey, config certutil.Config) (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := certutil.NewPrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "create private key")
	}

	cert, err := certutil.NewSignedCert(config, key, caCert, caKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to sign certificate")
	}

	return cert, key, nil
}

func newEncodedCertAndKey(caCert *x509.Certificate, caKey *rsa.PrivateKey, config certutil.Config) ([]byte, []byte, error) {
	crt, key, err := newCertAndKey(caCert, caKey, config)
	if err != nil {
		return nil, nil, err
	}

	encoded, err := Encode(&Pair{crt, key})
	if err != nil {
		return nil, nil, err
	}

	return encoded.Cert, encoded.Key, nil
}
