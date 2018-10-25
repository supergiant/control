package kube

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/sghelm/proxy"
)

func helmProxyFrom(kube *model.Kube) (proxy.Interface, error) {
	if kube == nil {
		return nil, errors.Wrap(sgerrors.ErrNilEntity, "kube model")
	}

	restConf, err := NewConfigFor(kube)
	if err != nil {
		return nil, err
	}

	coreV1Client, err := corev1.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}

	return proxy.New(coreV1Client, restConf, "")
}
