package kube

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"

	"k8s.io/client-go/kubernetes"

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

func syncMachines(ctx context.Context, k *model.Kube, account *model.CloudAccount) error {
	config := &steps.Config{}
	if err := util.FillCloudAccountCredentials(account, config); err != nil {
		return errors.Wrap(err, "error fill cloud account credentials")
	}

	config.AWSConfig.Region = k.Region
	EC2, err := amazon.GetEC2(config.AWSConfig)

	if err != nil {
		return errors.Wrap(sgerrors.ErrInvalidCredentials, err.Error())
	}

	describeInstanceOutput, err := EC2.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", clouds.TagClusterID)),
				Values: aws.StringSlice([]string{k.ID}),
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "describe instances")
	}

	for _, res := range describeInstanceOutput.Reservations {
		for _, instance := range res.Instances {
			node := &model.Machine{
				Size:   *instance.InstanceType,
				State:  model.MachineStateActive,
				Role:   model.RoleNode,
				Region: k.Region,
			}

			if instance.PublicIpAddress != nil {
				node.PublicIp = *instance.PublicIpAddress
			}

			if instance.PrivateIpAddress != nil {
				node.PrivateIp = *instance.PrivateIpAddress
			}

			for _, tag := range instance.Tags {
				if tag.Key != nil && *tag.Key == clouds.TagNodeName {
					node.Name = *tag.Value
				}
			}

			isFound := false

			for _, machine := range k.Nodes {
				if instance.PrivateIpAddress != nil && machine.PrivateIp == *instance.PrivateIpAddress {
					isFound = true
				}
			}

			var state int64

			if instance.State != nil && instance.State.Code != nil {
				state = *instance.State.Code
			}

			// If node is new in workers and it is not a master
			if !isFound && k.Masters[node.Name] == nil && state == 16 {
				logrus.Debugf("Add new node %v", node)
				k.Nodes[node.Name] = node
			}
		}
	}

	return nil
}

func findNextMinorVersion(current string, versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	for i := 0; i < len(versions)-1; i++ {
		if (len(versions[i]) > 3 && len(current) > 3) && strings.EqualFold(versions[i][:4], current[:4]) {
			return versions[i+1]
		}
	}

	return ""
}

func discoverK8SVersion(kubeConfig *clientcmddapi.Config) (string, error) {
	restConf, err := clientcmd.NewNonInteractiveClientConfig(
		*kubeConfig,
		kubeConfig.CurrentContext,
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()

	if err != nil {
		return "", errors.Wrapf(err, "create rest config")
	}

	restConf.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if len(restConf.UserAgent) == 0 {
		restConf.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConf)

	if err != nil {
		return "", errors.Wrapf(err, "error create discovery client")
	}

	serverVersion, err := discoveryClient.ServerVersion()

	if err != nil {
		return "", errors.Wrapf(err, "error getting server version")
	}

	return strings.TrimPrefix(serverVersion.GitVersion, "v"), nil
}

func discoverHelmVersion(kubeConfig *clientcmddapi.Config) (string, error) {
	restConf, err := clientcmd.NewNonInteractiveClientConfig(
		*kubeConfig,
		kubeConfig.CurrentContext,
		&clientcmd.ConfigOverrides{},
		nil,
	).ClientConfig()

	if err != nil {
		return "", errors.Wrapf(err, "create rest config")
	}

	restConf.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if len(restConf.UserAgent) == 0 {
		restConf.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	clientSet, err := kubernetes.NewForConfig(restConf)

	if err != nil {
		return "", errors.Wrapf(err, "get client set")
	}

	deploymentList, err := clientSet.AppsV1().Deployments("kube-system").List(v1.ListOptions{})

	if err != nil {
		return "", errors.Wrapf(err, "list deployments")
	}

	for _, deployment := range deploymentList.Items {
		if strings.Contains(deployment.Name, "tiller") {
			for _, container := range deployment.Spec.Template.Spec.Containers {
				slice := strings.Split(container.Image, ":")

				if len(slice) > 1 {
					return strings.Trim(slice[1], "v"), nil
				}
			}
		}
	}

	return "", nil
}
