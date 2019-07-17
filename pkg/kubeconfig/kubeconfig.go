package kubeconfig

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewConfigFor(k *model.Kube) (*rest.Config, error) {
	kubeConf, err := AdminKubeConfig(k)
	if err != nil {
		return nil, errors.Wrap(err, "build kubeconfig")
	}

	restConf, err := clientcmd.NewNonInteractiveClientConfig(
		kubeConf,
		kubeConf.CurrentContext,
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()

	restConf.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if len(restConf.UserAgent) == 0 {
		restConf.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return restConf, errors.Wrap(err, "build rest config")
}

func RestClientForGroupVersion(k *model.Kube, gv schema.GroupVersion) (rest.Interface, error) {
	cfg, err := NewConfigFor(k)
	if err != nil {
		return nil, err
	}

	setGroupDefaults(cfg, gv)
	return rest.RESTClientFor(cfg)
}

func discoveryClient(k *model.Kube) (*discovery.DiscoveryClient, error) {
	cfg, err := NewConfigFor(k)
	if err != nil {
		return nil, err
	}
	return discovery.NewDiscoveryClientForConfig(cfg)
}

func CoreV1Client(k *model.Kube) (corev1client.CoreV1Interface, error) {
	cfg, err := NewConfigFor(k)
	if err != nil {
		return nil, err
	}
	return corev1client.NewForConfig(cfg)
}

// adminKubeConfig returns a cluster-admin kubeconfig for provided cluster.
func AdminKubeConfig(k *model.Kube) (clientcmddapi.Config, error) {
	// TODO: this should be an address of the master load balancer
	if k == nil || (k.ExternalDNSName == "" && len(k.Masters) == 0) {
		// TODO: use another base error, not ErrNotFound
		return clientcmddapi.Config{}, errors.Wrap(sgerrors.ErrNotFound, "master nodes")
	}

	var apiAddr string
	if k.ExternalDNSName != "" {
		if strings.HasPrefix(k.ExternalDNSName, "https") {
			apiAddr = fmt.Sprintf("%s:%d", k.ExternalDNSName, k.APIServerPort)
		} else {
			apiAddr = fmt.Sprintf("https://%s:%d", k.ExternalDNSName, k.APIServerPort)
		}
	} else {
		// Use public IP in case if DNS name is absent
		m := util.GetRandomNode(k.Masters)
		apiAddr = fmt.Sprintf("https://%s:%d", m.PublicIp, k.APIServerPort)
	}

	// TODO: add validation
	return clientcmddapi.Config{
		AuthInfos: map[string]*clientcmddapi.AuthInfo{
			adminContext(k.Name): {
				ClientCertificateData: []byte(k.Auth.AdminCert),
				ClientKeyData:         []byte(k.Auth.AdminKey),
			},
		},
		Clusters: map[string]*clientcmddapi.Cluster{
			k.Name: {
				Server:                   apiAddr,
				CertificateAuthorityData: []byte(k.Auth.CACert),
			},
		},
		Contexts: map[string]*clientcmddapi.Context{
			adminContext(k.Name): {
				AuthInfo: adminContext(k.Name),
				Cluster:  k.Name,
			},
		},
		CurrentContext: adminContext(k.Name),
	}, nil
}

func setGroupDefaults(config *rest.Config, gv schema.GroupVersion) {
	config.GroupVersion = &gv
	if len(gv.Group) == 0 {
		config.APIPath = "/api"
	} else {
		config.APIPath = "/apis"
	}
}

func adminContext(clusterName string) string {
	return "admin@" + clusterName
}
