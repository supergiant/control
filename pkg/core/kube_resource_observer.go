package core

import (
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type KubeResourceObserver struct {
	core *Core
}

func (s *KubeResourceObserver) Perform() (err error) {
	var kubeResources []*model.KubeResource
	if err = s.core.DB.Preload("Kube").Find(&kubeResources); err != nil {
		return err
	}
	for _, kubeResource := range kubeResources {
		// TODO this is a bit.. freestyle right now.
		// We set i there and let the Save in Refresh() persist it.
		// Don't care about errors from Heapster.
		if kubeResource.Kind == "Pod" {
			cpuMetrics, _ := s.core.K8S(kubeResource.Kube).ListPodHeapsterCPUUsageMetrics(kubeResource.Namespace, kubeResource.Name)
			ramMetrics, _ := s.core.K8S(kubeResource.Kube).ListPodHeapsterRAMUsageMetrics(kubeResource.Namespace, kubeResource.Name)
			if cpuMetrics != nil && ramMetrics != nil {
				kubeResource.ExtraData = map[string]interface{}{
					"metrics": map[string][]*kubernetes.HeapsterMetric{
						"cpu_usage": cpuMetrics,
						"ram_usage": ramMetrics,
					},
				}
			}
		}

		if err := s.core.KubeResources.Refresh(kubeResource); err != nil {
			return err
		}
	}
	return nil
}
