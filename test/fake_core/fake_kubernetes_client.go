package fake_core

import (
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/kubernetes"
)

type KubernetesClient struct {
	EnsureNamespaceFn                func(name string) error
	GetResourceFn                    func(kind string, namespace string, name string, out *json.RawMessage) error
	CreateResourceFn                 func(kind string, namespace string, objIn map[string]interface{}, out *json.RawMessage) error
	DeleteResourceFn                 func(kind string, namespace string, name string) error
	ListNamespacesFn                 func(query string) ([]*kubernetes.Namespace, error)
	ListEventsFn                     func(query string) ([]*kubernetes.Event, error)
	ListNodesFn                      func(query string) ([]*kubernetes.Node, error)
	ListPodsFn                       func(query string) ([]*kubernetes.Pod, error)
	ListNodeHeapsterStatsFn          func() ([]*kubernetes.HeapsterStats, error)
	ListPodHeapsterCPUUsageMetricsFn func(namespace string, name string) ([]*kubernetes.HeapsterMetric, error)
	ListPodHeapsterRAMUsageMetricsFn func(namespace string, name string) ([]*kubernetes.HeapsterMetric, error)
}

func (k *KubernetesClient) EnsureNamespace(name string) error {
	if k.EnsureNamespaceFn == nil {
		return nil
	}
	return k.EnsureNamespaceFn(name)
}

func (k *KubernetesClient) GetResource(kind string, namespace string, name string, out *json.RawMessage) error {
	if k.GetResourceFn == nil {
		return nil
	}
	return k.GetResourceFn(kind, namespace, name, out)
}

func (k *KubernetesClient) CreateResource(kind string, namespace string, objIn map[string]interface{}, out *json.RawMessage) error {
	if k.CreateResourceFn == nil {
		return nil
	}
	return k.CreateResourceFn(kind, namespace, objIn, out)
}

func (k *KubernetesClient) DeleteResource(kind string, namespace string, name string) error {
	if k.DeleteResourceFn == nil {
		return nil
	}
	return k.DeleteResourceFn(kind, namespace, name)
}

func (k *KubernetesClient) ListNamespaces(query string) ([]*kubernetes.Namespace, error) {
	if k.ListNamespacesFn == nil {
		return nil, nil
	}
	return k.ListNamespacesFn(query)
}

func (k *KubernetesClient) ListEvents(query string) ([]*kubernetes.Event, error) {
	if k.ListEventsFn == nil {
		return nil, nil
	}
	return k.ListEventsFn(query)
}

func (k *KubernetesClient) ListNodes(query string) ([]*kubernetes.Node, error) {
	if k.ListNodesFn == nil {
		return nil, nil
	}
	return k.ListNodesFn(query)
}

func (k *KubernetesClient) ListPods(query string) ([]*kubernetes.Pod, error) {
	if k.ListPodsFn == nil {
		return nil, nil
	}
	return k.ListPodsFn(query)
}

func (k *KubernetesClient) ListNodeHeapsterStats() ([]*kubernetes.HeapsterStats, error) {
	if k.ListNodeHeapsterStatsFn == nil {
		return nil, nil
	}
	return k.ListNodeHeapsterStatsFn()
}

func (k *KubernetesClient) ListPodHeapsterCPUUsageMetrics(namespace string, name string) ([]*kubernetes.HeapsterMetric, error) {
	if k.ListPodHeapsterCPUUsageMetricsFn == nil {
		return nil, nil
	}
	return k.ListPodHeapsterCPUUsageMetricsFn(namespace, name)
}

func (k *KubernetesClient) ListPodHeapsterRAMUsageMetrics(namespace string, name string) ([]*kubernetes.HeapsterMetric, error) {
	if k.ListPodHeapsterRAMUsageMetricsFn == nil {
		return nil, nil
	}
	return k.ListPodHeapsterRAMUsageMetricsFn(namespace, name)
}
