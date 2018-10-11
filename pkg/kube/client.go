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
	if len(k.Masters) == 0 {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "master nodes")
	}
	m := util.GetRandomNode(k.Masters)

	// TODO: there is no Scheme/MasterPort config parameters
	apiAdrr := fmt.Sprintf("%s", m.PublicIp)

	return clientcmd.NewNonInteractiveClientConfig(
		buildKubeConfig(apiAdrr, k.Auth),
		"",
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()
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

// buildKubeConfig returns a kube config for provided options.
func buildKubeConfig(endpoint string, auth model.Auth) clientcmddapi.Config {
	return clientcmddapi.Config{
		AuthInfos: map[string]*clientcmddapi.AuthInfo{
			auth.Username: {
				Token:                 auth.Token,
				ClientCertificateData: []byte(auth.Cert),
				ClientKeyData:         []byte(auth.Key),
			},
		},
		Clusters: map[string]*clientcmddapi.Cluster{
			auth.Username: {
				Server:                   endpoint,
				CertificateAuthorityData: []byte(auth.CA),
			},
		},
		Contexts: map[string]*clientcmddapi.Context{
			auth.Username: {
				AuthInfo: auth.Username,
				Cluster:  auth.Username,
			},
		},
		CurrentContext: auth.Username,
	}
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
