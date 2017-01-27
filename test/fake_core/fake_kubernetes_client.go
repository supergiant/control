package fake_core

import "github.com/supergiant/supergiant/pkg/kubernetes"

type KubernetesClient struct {
	EnsureNamespaceFn                func(name string) error
	GetResourceFn                    func(apiVersion, kind, namespace, name string, out interface{}) error
	CreateResourceFn                 func(apiVersion, kind, namespace string, in, out interface{}) error
	UpdateResourceFn                 func(apiVersion, kind, namespace, name string, objIn interface{}, out interface{}) error
	DeleteResourceFn                 func(apiVersion, kind, namespace, name string) error
	ListNamespacesFn                 func(query string) ([]*kubernetes.Namespace, error)
	ListEventsFn                     func(query string) ([]*kubernetes.Event, error)
	ListNodesFn                      func(query string) ([]*kubernetes.Node, error)
	ListPodsFn                       func(query string) ([]*kubernetes.Pod, error)
	ListServicesFn                   func(query string) ([]*kubernetes.Service, error)
	ListPersistentVolumesFn          func(query string) ([]*kubernetes.PersistentVolume, error)
	ListNodeHeapsterStatsFn          func() ([]*kubernetes.HeapsterStats, error)
	ListPodHeapsterCPUUsageMetricsFn func(namespace, name string) ([]*kubernetes.HeapsterMetric, error)
	ListPodHeapsterRAMUsageMetricsFn func(namespace, name string) ([]*kubernetes.HeapsterMetric, error)
	GetPodLogFn                      func(namespace, name string) (string, error)
}

func (k *KubernetesClient) EnsureNamespace(name string) error {
	if k.EnsureNamespaceFn == nil {
		return nil
	}
	return k.EnsureNamespaceFn(name)
}

func (k *KubernetesClient) GetResource(apiVersion string, kind string, namespace string, name string, out interface{}) error {
	if k.GetResourceFn == nil {
		return nil
	}
	return k.GetResourceFn(apiVersion, kind, namespace, name, out)
}

func (k *KubernetesClient) CreateResource(apiVersion string, kind string, namespace string, objIn interface{}, out interface{}) error {
	if k.CreateResourceFn == nil {
		return nil
	}
	return k.CreateResourceFn(apiVersion, kind, namespace, objIn, out)
}

func (k *KubernetesClient) UpdateResource(apiVersion string, kind string, namespace string, name string, objIn interface{}, out interface{}) error {
	if k.UpdateResourceFn == nil {
		return nil
	}
	return k.UpdateResourceFn(apiVersion, kind, namespace, name, objIn, out)
}

func (k *KubernetesClient) DeleteResource(apiVersion string, kind string, namespace string, name string) error {
	if k.DeleteResourceFn == nil {
		return nil
	}
	return k.DeleteResourceFn(apiVersion, kind, namespace, name)
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

func (k *KubernetesClient) ListServices(query string) ([]*kubernetes.Service, error) {
	if k.ListServicesFn == nil {
		return nil, nil
	}
	return k.ListServicesFn(query)
}

func (k *KubernetesClient) ListPersistentVolumes(query string) ([]*kubernetes.PersistentVolume, error) {
	if k.ListPersistentVolumesFn == nil {
		return nil, nil
	}
	return k.ListPersistentVolumesFn(query)
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

func (k *KubernetesClient) GetPodLog(namespace string, name string) (string, error) {
	if k.GetPodLogFn == nil {
		return "", nil
	}
	return k.GetPodLogFn(namespace, name)
}
