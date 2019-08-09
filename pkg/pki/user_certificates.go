package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math"
	"math/big"
	"time"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
)

const (
	MastersGroup = "system:masters"

	duration365d = time.Hour * 24 * 365
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

	key, err := newPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "create private key")
	}

	cfg := certutil.Config{
		CommonName:   userName,
		Organization: userGroups,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	cert, err := newSignedCert(cfg, key, ca.Cert, ca.Key)
	if err != nil {
		return nil, errors.Wrap(err, "sign certificate")
	}

	return Encode(&Pair{
		Cert: cert,
		Key:  key,
	})
}

// newSignedCert creates a signed certificate using the given CA certificate and key
func newSignedCert(cfg certutil.Config, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}
	if len(cfg.Usages) == 0 {
		return nil, errors.New("must specify at least one ExtKeyUsage")
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     time.Now().Add(duration365d).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  cfg.Usages,
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}
