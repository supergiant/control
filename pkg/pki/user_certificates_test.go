package pki

import "testing"

func TestNewAdminPair(t *testing.T) {
	cert, key, _ := newCertificateAuthority()

	pair := &Pair{
		Cert: cert,
		Key:  key,
	}

	caPEMPair, _ := Encode(pair)

	pairPem, err := NewAdminPair(caPEMPair)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if pairPem == nil {
		t.Errorf("pair pem must not be nil")
	}
}

func TestNewUserPair(t *testing.T) {
	cert, key, _ := newCertificateAuthority()

	pair := &Pair{
		Cert: cert,
		Key:  key,
	}

	pemPair, _ := Encode(pair)

	pairPem, err := NewUserPair("kubernetes-admin",
		[]string{MastersGroup}, pemPair)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if pairPem == nil {
		t.Errorf("pair pem must not be nil")
	}
}
