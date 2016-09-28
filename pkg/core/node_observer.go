package core

import (
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type NodeObserver struct {
	core *Core
}

func (s *NodeObserver) Perform() error {
	var kubes []*model.Kube
	if err := s.core.DB.Where("ready = ?", true).Preload("CloudAccount").Preload("Nodes", "provider_id <> ?", "").Find(&kubes); err != nil {
		return err
	}

	for _, kube := range kubes {

		k8s := s.core.K8S(kube)

		metrics, err := k8s.ListNodeHeapsterStats()
		if err != nil {
			continue // very common to get an error here, and it's not critical
		}

		k8sNodes, err := k8s.ListNodes("")
		if err != nil {
			return err
		}

		for _, node := range kube.Nodes {

			var knode *kubernetes.Node
			for _, kn := range k8sNodes {
				if kn.Metadata.Name == node.Name {
					knode = kn
					break
				}
			}

			// Set ExternalIP
			for _, addr := range knode.Status.Addresses {
				if addr.Type == "ExternalIP" {
					node.ExternalIP = addr.Address
					break
				}
			}

			// Set OutOfDisk
			if len(knode.Status.Conditions) > 0 {
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

			var metric *kubernetes.HeapsterStats
			for _, mtrc := range metrics {
				if mtrc.Name == node.Name {
					metric = mtrc
					break
				}
			}

			// Set Stats
			if metric != nil {
				node.CPUUsage = metric.CPUUsage
				node.RAMUsage = metric.RAMUsage
				node.CPULimit = int64(nodeSize.CPUCores * 1000)
				node.RAMLimit = int64(nodeSize.RAMGIB * 1073741824)
			}

			if err := s.core.DB.Save(node); err != nil {
				return err
			}
		}
	}

	return nil
}
