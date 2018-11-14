package pki

import (
	"crypto/x509"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
)

const (
	MastersGroup = "system:masters"
)

// NewAdminPair creates certificates for the kubernetes admin user.
func NewAdminPair(ca *PairPEM) (*PairPEM, error) {
	return NewUserPair("kubernetes-admin", []string{MastersGroup}, ca)
}

// NewUserPair creates certificates for a kubernetes user.
func NewUserPair(userName string, userGroups []string, caEncoded *PairPEM) (*PairPEM, error) {
	ca, err := Decode(caEncoded)
	if err != nil {
		return nil, errors.Wrap(err, "decode ca cert/key")
	}

	key, err := certutil.NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "create private key")
	}

	cfg := certutil.Config{
		CommonName:   userName,
		Organization: userGroups,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	cert, err := certutil.NewSignedCert(cfg, key, ca.Cert, ca.Key)
	if err != nil {
		return nil, errors.Wrap(err, "sign certificate")
	}

	return Encode(&Pair{
		Cert: cert,
		Key:  key,
	})
}
