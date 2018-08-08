package pki

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateSelfSignedCAKey(t *testing.T) {
	pki, err := NewPKI(nil, "master", []net.IP{
		net.ParseIP("127.0.0.1"),
	})
	require.NoError(t, err)

	require.NotNil(t, pki.CA)
	require.NotNil(t, pki.CA.Key)
	require.True(t, len(pki.CA.Cert) > 0)
	require.True(t, len(pki.CA.Key) > 0)

	require.NotNil(t, pki.APIServer)
	require.True(t, len(pki.APIServer.Cert) > 0)
	require.True(t, len(pki.APIServer.Key) > 0)

	require.NotNil(t, pki.Kubelet)
	require.NotNil(t, len(pki.Kubelet.Cert))
	require.NotNil(t, len(pki.Kubelet.Key))
}
