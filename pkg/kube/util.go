package kube

import (
	"fmt"
	"strings"

	"github.com/supergiant/control/pkg/model"
)

func processAWSMetrics(k *model.Kube, metrics map[string]map[string]interface{}) {
	for _, masterNode := range k.Masters {
		// After some amount of time prometheus start using region in metric name
		prefix := ip2Host(masterNode.PrivateIp)
		for metricKey := range metrics {
			if strings.Contains(metricKey, prefix) {
				value := metrics[metricKey]
				delete(metrics, metricKey)
				metrics[strings.ToLower(masterNode.Name)] = value
			}
		}
	}

	for _, workerNode := range k.Nodes {
		prefix := ip2Host(workerNode.PrivateIp)

		for metricKey := range metrics {
			if strings.Contains(metricKey, prefix) {
				value := metrics[metricKey]
				delete(metrics, metricKey)
				metrics[strings.ToLower(workerNode.Name)] = value
			}
		}
	}
}

func ip2Host(ip string) string {
	return fmt.Sprintf("ip-%s", strings.Join(strings.Split(ip, "."), "-"))
}
