package kube

import (
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/kubeconfig"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/sghelm/proxy"
)

func helmProxyFrom(kube *model.Kube) (proxy.Interface, error) {
	if kube == nil {
		return nil, errors.Wrap(sgerrors.ErrNilEntity, "kube model")
	}

	restConf, err := kubeconfig.NewConfigFor(kube)
	if err != nil {
		return nil, err
	}

	coreV1Client, err := corev1.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}

	return proxy.New(coreV1Client, restConf, "")
}
