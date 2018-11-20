package kube

import (
	"testing"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
)

func TestIp2Host(t *testing.T) {
	testCases := []struct{
		ip string
		host string
	}{
		{
			ip: "172.16.10.2",
			host: "ip-172-16-10-2",
		},
		{
			ip: "",
			host: "ip-",
		},
	}

	for _, testCase := range testCases {
		host := ip2Host(testCase.ip)

		if host != testCase.host {
			t.Errorf("Wrong host expected %s actual %s",
				testCase.host, host)
		}
	}
}

func TestProcessAWSMetrics(t *testing.T) {
	masters := map[string]*node.Node{
		"master-1": {
			Name: "master-1",
			PrivateIp: "10.20.30.40",
		},
	}

	nodes := map[string]*node.Node{
		"node-1": {
			Name: "node-1",
			PrivateIp: "172.16.0.1",
		},
		"node-2": {
			Name: "node-2",
			PrivateIp: "172.16.0.2",
		},
	}

	k := &model.Kube{
		Provider: clouds.AWS,
		Masters:  masters,
		Nodes:    nodes,
	}

	metrics := map[string]map[string]interface{}{
		"ip-10-20-30-40": {},
		"ip-172-16-0-1": {},
		"ip-172-16-0-2": {},
	}

	processAWSMetrics(k, metrics)

	for _, masterNode := range masters {
		if _, ok := metrics[masterNode.Name]; !ok {
			t.Errorf("Node %s not found in %v", masterNode.Name, metrics)
		}
	}

	for _, workerNode := range nodes {
		if _, ok := metrics[workerNode.Name]; !ok {
			t.Errorf("Node %s not found in %v", workerNode.Name, metrics)
		}
	}
}