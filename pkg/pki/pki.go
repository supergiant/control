package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/pkg/errors"
)

const (
	// RSAPrivateKeyBlockType is a possible value for pem.Block.Type.
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
	// PrivateKeyBlockType is a possible value for pem.Block.Type.
	PublicKeyBlockType = "PUBLIC KEY"
	// CertificateBlockType is a possible value for pem.Block.Type.
	CertificateBlockType = "CERTIFICATE"
)

// CARequest defines a request to generate or use CA if provided to setup PKI for k8s cluster
type CARequest struct {
	DNSDomain string   `json:"dnsDomain" valid:"required"`
	IPs       []string `json:"ips" valid:"required"`
	CA        []byte   `json:"ca" valid:"optional"`
}

// Pair defines a certificate and a private key.
type Pair struct {
	Cert *x509.Certificate
	Key  *rsa.PrivateKey
}

// PairPEM defines PEM encoded certificate and private key.
// TODO: user cert pair in the kube model or get rid of it.
type PairPEM struct {
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

// PKI defines a set of certificates/keys for a kubernetes cluster.
type PKI struct {
	ID string   `json:"id"`
	CA *PairPEM `json:"ca"`
	//KubeName is a sg specific name of a k8s cluster
	KubeName string `json:"kubeName"`
}

// Encode encodes cert/key with PEM and returns them as a PairPEM.
func Encode(p *Pair) (*PairPEM, error) {
	if p == nil || p.Cert == nil || p.Key == nil {
		return nil, ErrEmptyPair
	}
	return &PairPEM{
		//Cert: certutil.EncodeCertPEM(p.Cert),
		Cert: encodeCertPEM(p.Cert),
		Key:  encodePrivateKeyPEM(p.Key),
	}, nil
}

// Decode parses a pem encoded cert/key and returns them as a Pair.
func Decode(p *PairPEM) (*Pair, error) {
	if p == nil || p.Cert == nil || p.Key == nil {
		return nil, ErrEmptyPair
	}

	pemBytes, rest := pem.Decode(p.Cert)
	if len(rest) > 0 {
		return nil, errors.New("decode pem")
	}
	cert, err := x509.ParseCertificate(pemBytes.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse a raw certificate")
	}

	pemBytes, rest = pem.Decode(p.Key)

	if len(rest) > 0 {
		return nil, errors.New("decode pem")
	}

	key, err := x509.ParsePKCS1PrivateKey(pemBytes.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse a raw private key")
	}

	return &Pair{cert, key}, nil
}

// NewCAPair creates certificates and key for a kubernetes cluster.
// If no CA cert/key is provided, it creates self-signed ones.
func NewCAPair(parentBytes []byte) (*PairPEM, error) {
	var caPem *PairPEM

	if parentBytes == nil || len(parentBytes) == 0 {
		p, k, err := generateCACert()
		if err != nil {
			return nil, err
		}
		caPem = &PairPEM{Cert: p, Key: k}
	} else {
		pemBlock, rest := pem.Decode(parentBytes)
		if len(rest) > 0 {
			return nil, errors.New("error decode parent cert")
		}

		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "parse parent cert bytes")
		}

		certBytes, keyBytes, err := generateCertFromParent(cert)
		if err != nil {
			return nil, errors.Wrap(err, "create cert from parent")
		}

		caPem = &PairPEM{Cert: certBytes, Key: keyBytes}
	}

	ca, err := Decode(caPem)
	if err != nil {
		return nil, errors.Wrap(err, "decode a CA pair")
	}

	// Check that cert generates is CA cert
	if !ca.Cert.IsCA {
		return nil, ErrInvalidCA
	}

	return caPem, nil
}

// EncodePublicKeyPEM returns PEM-encoded public data
func EncodePublicKeyPEM(key *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return []byte{}, err
	}
	block := pem.Block{
		Type:  PublicKeyBlockType,
		Bytes: der,
	}
	return pem.EncodeToMemory(&block), nil
}

// generateCACert will generate a self-signed CA certificate
func generateCACert() ([]byte, []byte, error) {
	crt, key, err := newCertificateAuthority()
	if err != nil {
		return nil, nil, errors.Wrap(err, "generating self signed CA")
	}
	pmCrt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE",
		Bytes: crt.Raw})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return pmCrt, keyBytes, nil
}

func generateCertFromParent(parent *x509.Certificate) ([]byte, []byte, error) {
	// Generate a key.
	key, err := newPrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "generate private key")
	}
	// Fill out the template.
	template := x509.Certificate{
		SerialNumber:          new(big.Int).SetInt64(0),
		Subject:               pkix.Name{Organization: []string{"Qbox Inc"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if parent.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	if parent == nil {
		parent = &template
	}
	// Generate the certificate.
	// TODO: there is no ca key, is it valid?
	cert, err := x509.CreateCertificate(rand.Reader, &template, parent, &key.PublicKey, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create certificate from parent")
	}
	// Marshal the key.
	b := x509.MarshalPKCS1PrivateKey(key)

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}),
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}),
		nil
}

// newPrivateKey creates a RSA private key.
func newPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// encodePrivateKeyPEM returns PEM-encoded private key data
func encodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

// EncodeCertPEM returns PEM-endcoded certificate data
func encodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}