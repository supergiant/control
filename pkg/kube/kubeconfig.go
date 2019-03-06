package kube

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
)

func NewConfigFor(k *model.Kube) (*rest.Config, error) {
	kubeConf, err := adminKubeConfig(k)
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

func restClientForGroupVersion(k *model.Kube, gv schema.GroupVersion) (rest.Interface, error) {
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

func corev1Client(k *model.Kube) (corev1client.CoreV1Interface, error) {
	cfg, err := NewConfigFor(k)
	if err != nil {
		return nil, err
	}
	return corev1client.NewForConfig(cfg)
}

// adminKubeConfig returns a cluster-admin kubeconfig for provided cluster.
func adminKubeConfig(k *model.Kube) (clientcmddapi.Config, error) {
	// TODO: this should be an address of the master load balancer
	if k == nil || len(k.Masters) == 0 {
		// TODO: use another base error, not ErrNotFound
		return clientcmddapi.Config{}, errors.Wrap(sgerrors.ErrNotFound, "master nodes")
	}
	m := util.GetRandomNode(k.Masters)

	var apiAddr string
	if k.APIPort != "" {
		apiAddr = fmt.Sprintf("https://%s:%s", m.PublicIp, k.APIPort)
	} else {
		// TODO: apiPort has been hardcoded in provisioner, use 443 by default
		apiAddr = fmt.Sprintf("https://%s", m.PublicIp)
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
				Server: apiAddr,
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
