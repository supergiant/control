package pki

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"github.com/pkg/errors"
)

// CARequest defines a request to generate or use CA if provided to setup PKI for k8s cluster
type CARequest struct {
	DNSDomain string   `json:"dnsDomain" valid:"required"`
	IPs       []string `json:"ips" valid:"required"`
	CA        *PairPEM `json:"ca" valid:"optional"`
}

// Pair defines a certificate and a private key.
type Pair struct {
	Cert *x509.Certificate
	Key  *rsa.PrivateKey
}

// PairPEM defines PEM encoded certificate and private key.
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

func (p *PKI) Marshall() []byte {
	data, _ := json.Marshal(p)
	return data
}

// Encode encodes cert/key with PEM and returns them as a PairPEM.
func Encode(p *Pair) (*PairPEM, error) {
	encoded := new(PairPEM)
	buf := new(bytes.Buffer)

	err := pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: p.Cert.Raw})
	if err != nil {
		return nil, errors.Wrap(err, "encode a certificate with pem")

	}
	encoded.Cert = buf.Bytes()

	buf.Reset()
	err = pem.Encode(buf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p.Key)})
	if err != nil {
		return nil, errors.Wrap(err, "encode a private key with pem")
	}
	encoded.Key = buf.Bytes()

	return encoded, nil
}

// Decode parses a pem encoded cert/key and returns them as a Pair.
func Decode(p *PairPEM) (*Pair, error) {
	pemBytes, rest := pem.Decode(p.Cert)
	if len(rest) > 0 {
		return nil, errors.New("decode pem")
	}
	cert, err := x509.ParseCertificate(pemBytes.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse a raw certificate")
	}

	key, err := x509.ParsePKCS1PrivateKey(p.Key)
	if err != nil {
		return nil, errors.Wrap(err, "parse a raw private key")
	}

	return &Pair{cert, key}, nil
}

// NewPKI creates certificates and key for a kubernetes cluster.
// If no CA cert/key is provided, it creates self-signed ones.
func NewPKI(caPEM *PairPEM) (*PKI, error) {
	if caPEM == nil {
		p, k, err := generateCACert()
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		caPEM = &PairPEM{Cert: p, Key: k}
	}
	ca, err := Decode(caPEM)
	if err != nil {
		return nil, errors.Wrap(err, "decode a CA pair")
	}

	if !ca.Cert.IsCA {
		return nil, ErrInvalidCA
	}

	return &PKI{
		CA: caPEM,
	}, nil
}

//generateCACert will generate a self-signed CA certificate
func generateCACert() ([]byte, []byte, error) {
	crt, key, err := newCertificateAuthority()
	if err != nil {
		return nil, nil, errors.Wrap(err, "generating self signed CA")
	}
	pmCrt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE",
		Bytes: crt.Raw})

	return pmCrt, x509.MarshalPKCS1PrivateKey(key), nil
}
