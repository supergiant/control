package pki

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"

	certutil "k8s.io/client-go/util/cert"
)

const (
	// APIServerCertCommonName defines API's server certificate common name (CN)
	APIServerCertCommonName = "kube-apiserver"

	// APIServerKubeletClientCertCommonName defines kubelet client certificate common name (CN)
	APIServerKubeletClientCertCommonName = "kube-apiserver-kubelet-client"

	// MastersGroup defines the well-known group for the apiservers. This group is also superuser by default
	// (i.e. bound to the cluster-admin ClusterRole)
	MastersGroup = "system:masters"
)

// NewAPIServerCertAndKey generate certificate for apiserver, signed by the given CA.
func NewAPIServerCertAndKey(caCert *x509.Certificate, caKey *rsa.PrivateKey, dnsDomain string, ips []net.IP) ([]byte, []byte, error) {
	config := certutil.Config{
		CommonName: APIServerCertCommonName,
		AltNames:   getAPIServerAltNames(dnsDomain, ips),
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	return newEncodedCertAndKey(caCert, caKey, config)
}

// NewAPIServerKubeletClientCertAndKey generate certificate for the apiservers to connect to the kubelets securely, signed by the given CA.
func NewAPIServerKubeletClientCertAndKey(caCert *x509.Certificate, caKey *rsa.PrivateKey) ([]byte, []byte, error) {
	config := certutil.Config{
		CommonName:   APIServerKubeletClientCertCommonName,
		Organization: []string{MastersGroup},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	return newEncodedCertAndKey(caCert, caKey, config)
}

func getAPIServerAltNames(dnsDomain string, ips []net.IP) certutil.AltNames {
	return certutil.AltNames{
		DNSNames: []string{
			"kubernetes",
			"kubernetes.default",
			"kubernetes.default.svc",
			fmt.Sprintf("kubernetes.default.svc.%s", dnsDomain),
		},
		IPs: ips,
	}
}
