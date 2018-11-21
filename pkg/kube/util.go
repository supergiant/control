package kube

import (
	"fmt"
	"strings"

	"github.com/supergiant/control/pkg/model"
)

func processAWSMetrics(k *model.Kube, metrics map[string]map[string]interface{}) {
	for _, masterNode := range k.Masters {
		key := ip2Host(masterNode.PrivateIp)
		value := metrics[key]
		delete(metrics, key)
		metrics[strings.ToLower(masterNode.Name)] = value
	}

	for _, workerNode := range k.Nodes {
		key := ip2Host(workerNode.PrivateIp)
		value := metrics[key]
		delete(metrics, key)
		metrics[strings.ToLower(workerNode.Name)] = value
	}
}

func ip2Host(ip string) string {
	return fmt.Sprintf("ip-%s",strings.Join(strings.Split(ip, "."), "-"))
}
