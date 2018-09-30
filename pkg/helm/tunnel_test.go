package helm

import (
	"k8s.io/client-go/rest"
	"testing"
)

func TestNewTunnel(t *testing.T) {
	var (
		client    rest.Interface = nil
		config                   = &rest.Config{}
		namespace                = "namespace"
		pod                      = "podname"
		remote                   = 0
	)

	tunnel := newTunnel(client, config, namespace, pod, remote)

	if tunnel.client != client {
		t.Errorf("expected client %v actual %v", client, tunnel.client)
	}

	if tunnel.config != config {
		t.Errorf("expected config %v actual %v", config, tunnel.config)
	}

	if tunnel.Namespace != namespace {
		t.Errorf("expected namespace %s actual %s", namespace,
			tunnel.Namespace)
	}

	if tunnel.PodName != pod {
		t.Errorf("expected pod %s actual %s", pod, tunnel.PodName)
	}

	if tunnel.Remote != remote {
		t.Errorf("expected remote %d actual %d", remote, tunnel.Remote)
	}
}

func TestTunnelClose(t *testing.T) {
	tunnel := newTunnel(nil, nil, "", "", 1)
	tunnel.close()
}

func TestTunnelClosePanic(t *testing.T) {
	tunnel := newTunnel(nil, nil, "", "", 1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic")
		}
	}()
	tunnel.close()
	tunnel.close()
}
