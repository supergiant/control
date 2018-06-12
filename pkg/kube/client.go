package kube

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"
)

func restClientForGroupVersion(k *Kube, gv schema.GroupVersion) (rest.Interface, error) {
	cfg, err := buildConfig(k.APIAddr, k.Auth)
	if err != nil {
		return nil, err
	}

	setGroupDefaults(cfg, gv)
	return rest.RESTClientFor(cfg)
}

func discoveryClient(k *Kube) (*discovery.DiscoveryClient, error) {
	cfg, err := buildConfig(k.APIAddr, k.Auth)
	if err != nil {
		return nil, err
	}

	return discovery.NewDiscoveryClientForConfig(cfg)
}

// buildKubeConfig returns a kube config for provided options.
func buildKubeConfig(addr string, auth Auth) clientcmddapi.Config {
	return clientcmddapi.Config{
		AuthInfos: map[string]*clientcmddapi.AuthInfo{
			auth.Username: {
				Token: auth.Token,
				ClientCertificateData: []byte(auth.Cert),
				ClientKeyData:         []byte(auth.Key),
			},
		},
		Clusters: map[string]*clientcmddapi.Cluster{
			auth.Username: {
				Server: addr,
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

func buildConfig(addr string, auth Auth) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveClientConfig(
		buildKubeConfig(addr, auth),
		"",
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()
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
