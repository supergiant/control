package kube

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
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

func kubeFromKubeConfig(kubeConfig clientcmddapi.Config) (*model.Kube, error) {
	currentCtxName := kubeConfig.CurrentContext
	currentContext := kubeConfig.Contexts[currentCtxName]

	if currentContext == nil {
		return nil, errors.Wrapf(sgerrors.ErrNilEntity, "current context %s not found in context map %v",
			currentCtxName, kubeConfig.Contexts)
	}

	authInfoName := currentContext.AuthInfo
	authInfo := kubeConfig.AuthInfos[authInfoName]

	if authInfo == nil {
		return nil, errors.Wrapf(sgerrors.ErrNilEntity, "authInfo %s not found in auth into auth map %v",
			authInfoName, kubeConfig.AuthInfos)
	}

	clusterName := currentContext.Cluster
	cluster := kubeConfig.Clusters[clusterName]

	if cluster == nil {
		return nil, errors.Wrapf(sgerrors.ErrNilEntity, "cluster %s not found in cluster map %v",
			clusterName, kubeConfig.Clusters)
	}

	return &model.Kube{
		Name:            currentContext.Cluster,
		ExternalDNSName: cluster.Server,
		Auth: model.Auth{
			CACert:    string(cluster.CertificateAuthorityData),
			AdminCert: string(authInfo.ClientCertificateData),
			AdminKey:  string(authInfo.ClientKeyData),
		},
	}, nil
}
