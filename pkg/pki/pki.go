package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"math/big"
	"time"

	"github.com/supergiant/supergiant/pkg/sgerrors"
	certutil "k8s.io/client-go/util/cert"
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
	} else {
		ca, err := Decode(caPEM)
		if err != nil {
			return nil, errors.Wrap(err, "decode parent pem")
		}

		cert, key, err := generateCertFromParent(ca.Cert)

		if err != nil {
			return nil, errors.Wrap(err, "create cert from parent")
		}
		caPEM = &PairPEM{Cert: cert, Key: key}
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
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return pmCrt, keyBytes, nil
}

func generateCertFromParent(parent *x509.Certificate) ([]byte, []byte, error) {
	if parent == nil {
		return nil, nil, sgerrors.ErrInvalidCredentials
	}

	// Generate a key.
	key, err := certutil.NewPrivateKey()
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
