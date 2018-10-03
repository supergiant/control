package pki

import (
	"testing"

	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
)

func TestGenerateSelfSignedCAKey(t *testing.T) {
	pki, err := NewPKI(nil)
	require.NoError(t, err)

	require.NotNil(t, pki.CA)
	require.NotNil(t, pki.CA.Key)
	require.True(t, len(pki.CA.Cert) > 0)
	require.True(t, len(pki.CA.Key) > 0)
}

func TestPKI_Marshall(t *testing.T) {
	id := "1"
	kubeName := "kubeName"
	ca := &PairPEM{
		Key:  []byte(`key`),
		Cert: []byte(`pem`),
	}

	pki := &PKI{
		ID:       id,
		CA:       ca,
		KubeName: kubeName,
	}

	data := pki.Marshall()

	if len(data) == 0 {
		t.Errorf("empty array of bytes")
	}

	if !bytes.Contains(data, []byte(kubeName)) {
		t.Errorf("kubeName %s not found in %s", kubeName, string(data))
	}
}

func TestUnmarshal(t *testing.T) {
	ca := &PairPEM{
		Key:  []byte(`key`),
		Cert: []byte(`pem`),
	}

	pki := &PKI{
		ID:       "id",
		CA:       ca,
		KubeName: "name",
	}

	data := pki.Marshall()
	pki2 := &PKI{}

	err := json.Unmarshal(data, pki2)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if pki.KubeName != pki2.KubeName {
		t.Errorf("wrong kube name expected %s actual %s",
			pki.KubeName, pki2.KubeName)
	}

	if !bytes.Equal(pki.CA.Cert, pki2.CA.Cert) {
		t.Errorf("Wrong CA cert expected %s actual %s",
			string(pki.CA.Cert), string(pki.CA.Cert))
	}
}
