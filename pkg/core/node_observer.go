package core

import (
	"time"

	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type NodeObserver struct {
	core *Core
}

// Perform - Gathers metric information on nodes.
func (s *NodeObserver) Perform() error {
	var kubes []*model.Kube
	if err := s.core.DB.Where("ready = ?", true).Preload("CloudAccount").Preload("Nodes", "provider_id <> ?", "").Find(&kubes); err != nil {
		return err
	}

	for _, kube := range kubes {
		k8s := s.core.K8S(kube)

		k8sNodes, err := k8s.ListNodes("")
		if err != nil {
			return err
		}

		// Kube level metrics
		kubemetrics, err := k8s.ListKubeHeapsterStats()
		if err != nil {
			continue // very common to get an error here, and it's not critical
		}

		if len(kube.ExtraData) == 0 {
			kube.ExtraData = map[string]interface{}{
				"metrics": map[string][]*kubernetes.HeapsterMetric{},
			}
		}
		for _, metric := range kubemetrics {
			mets, err := k8s.GetKubeHeapsterStats(metric)
			if err != nil {
				continue // very common to get an error here, and it's not critical()
			}
			kube.ExtraData[mets.MetricName] = mets.Metrics
		}
		if err := s.core.DB.Save(kube); err != nil {
			return err
		}

		if err := s.core.DB.Save(kube); err != nil {
			return err
		}

		kubeCPU := int64(0)
		kubeRAM := int64(0)
		for _, node := range kube.Nodes {

			// node level metrics
			metrics, err := k8s.ListNodeHeapsterStats(node.Name)
			if err != nil {
				continue // very common to get an error here, and it's not critical
			}

			node.ExtraData = map[string]interface{}{
				"metrics": map[string][]*kubernetes.HeapsterMetric{},
			}
			metData := map[string][]*kubernetes.HeapsterMetric{}
			for _, metric := range metrics {
				mets, err := k8s.GetNodeHeapsterStats(node.Name, metric)
				if err != nil {
					continue // very common to get an error here, and it's not critical
				}
				metData[mets.MetricName] = mets.Metrics
				node.ExtraData[mets.MetricName] = mets.Metrics
			}

			var knode *kubernetes.Node
			for _, kn := range k8sNodes {
				if kn.Metadata.Name == node.Name {
					knode = kn
					break
				}
			}

			if knode != nil {
				// Set ExternalIP
				for _, addr := range knode.Status.Addresses {
					if addr.Type == "ExternalIP" {
						node.ExternalIP = addr.Address
						break
					}
				}

				// Set OutOfDisk
				for _, condition := range knode.Status.Conditions {
					if condition.Type == "OutOfDisk" {
						node.OutOfDisk = condition.Status == "True"
						break
					}
				}
			}

			var nodeSize *NodeSize
			for _, ns := range s.core.NodeSizes[kube.CloudAccount.Provider] {
				if ns.Name == node.Size {
					nodeSize = ns
					break
				}
			}

			for metricType, metricValue := range metData {
				switch metricType {
				case "cpu_usage_rate":
					node.CPUUsage = metricValue[len(metricValue)-1].Value
				case "memory_usage":
					node.RAMUsage = metricValue[len(metricValue)-1].Value
				case "cpu_limit":
					node.CPULimit = int64(nodeSize.CPUCores * 1000)
				case "memory_limit":
					node.RAMLimit = int64(nodeSize.RAMGIB * 1073741824)
				case "cpu_node_capacity":
					kubeCPU = kubeCPU + metricValue[len(metricValue)-1].Value
				case "memory_node_capacity":
					kubeRAM = kubeRAM + metricValue[len(metricValue)-1].Value
				}
			}

			if err := s.core.DB.Save(node); err != nil {
				return err
			}
		}

		// Compile total kube cpu
		switch kube.ExtraData["kube_cpu_capacity"].(type) {
		case nil:
			kube.ExtraData["kube_cpu_capacity"] = []*kubernetes.HeapsterMetric{}
		case []interface{}:
			kube.ExtraData["kube_cpu_capacity"] = append(
				kube.ExtraData["kube_cpu_capacity"].([]interface{}),
				&kubernetes.HeapsterMetric{
					Timestamp: time.Now(),
					Value:     kubeCPU,
				})
			for len(kube.ExtraData["kube_cpu_capacity"].([]interface{})) > 45 {
				kube.ExtraData["kube_cpu_capacity"] = append(
					kube.ExtraData["kube_cpu_capacity"].([]interface{})[:0],
					kube.ExtraData["kube_cpu_capacity"].([]interface{})[1:]...)
			}
		}

		// compile total kube ram
		switch kube.ExtraData["kube_memory_capacity"].(type) {
		case nil:
			kube.ExtraData["kube_memory_capacity"] = []*kubernetes.HeapsterMetric{}
		case []interface{}:
			kube.ExtraData["kube_memory_capacity"] = append(
				kube.ExtraData["kube_memory_capacity"].([]interface{}),
				&kubernetes.HeapsterMetric{
					Timestamp: time.Now(),
					Value:     kubeRAM,
				})
			for len(kube.ExtraData["kube_memory_capacity"].([]interface{})) > 45 {
				kube.ExtraData["kube_memory_capacity"] = append(
					kube.ExtraData["kube_memory_capacity"].([]interface{})[:0],
					kube.ExtraData["kube_memory_capacity"].([]interface{})[1:]...)
			}
		}

		if err := s.core.DB.Save(kube); err != nil {
			return err
		}
	}
	return nil
}
