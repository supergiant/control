package kube

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
)

func NewConfigFor(k *model.Kube) (*rest.Config, error) {
	kubeConf, err := kubeConfigFrom(k)
	if err != nil {
		return nil, errors.Wrap(err, "build kubeconfig")
	}

	restConf, err := clientcmd.NewNonInteractiveClientConfig(
		kubeConf,
		kubeConf.CurrentContext,
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()

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

// kubeConfigFrom returns a kube config for provided cluster.
func kubeConfigFrom(k *model.Kube) (clientcmddapi.Config, error) {
	if len(k.Masters) == 0 {
		// TODO: use another base error, not ErrNotFound
		return clientcmddapi.Config{}, errors.Wrap(sgerrors.ErrNotFound, "master nodes")
	}
	m := util.GetRandomNode(k.Masters)

	var apiAddr string
	if k.APIPort != "" {
		apiAddr = fmt.Sprintf("https://%s:%s", m.PublicIp, k.APIPort)
	} else {
		// TODO: apiPort has been hardcoded in provisioner, use it if no apiPort has been provided
		apiAddr = fmt.Sprintf("http://%s:8080", m.PublicIp)
	}

	return clientcmddapi.Config{
		AuthInfos: map[string]*clientcmddapi.AuthInfo{
			k.Auth.Username: {
				Token:                 k.Auth.Token,
				ClientCertificateData: []byte(k.Auth.Cert),
				ClientKeyData:         []byte(k.Auth.Key),
			},
		},
		Clusters: map[string]*clientcmddapi.Cluster{
			k.Auth.Username: {
				Server:                   apiAddr,
				CertificateAuthorityData: []byte(k.Auth.CA),
			},
		},
		Contexts: map[string]*clientcmddapi.Context{
			k.Auth.Username: {
				AuthInfo: k.Auth.Username,
				Cluster:  k.Auth.Username,
			},
		},
		CurrentContext: k.Auth.Username,
	}, nil
}

func setGroupDefaults(config *rest.Config, gv schema.GroupVersion) {
	config.GroupVersion = &gv
	if len(gv.Group) == 0 {
		config.APIPath = "/api"
	} else {
		config.APIPath = "/apis"
	}
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if len(config.UserAgent) == 0 {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
}
