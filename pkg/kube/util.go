package kube

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"

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

func createSpotInstance(req *SpotRequest, config *steps.Config) error {
	switch config.Provider {
	case clouds.AWS:
		return createAwsSpotInstance(req, config)
	}

	return sgerrors.ErrUnsupportedProvider
}

func createAwsSpotInstance(req *SpotRequest, config *steps.Config) error {
	svc, err := amazon.GetEC2(config.AWSConfig)

	if err != nil {
		return errors.Wrap(err, "get EC2 client")
	}

	config.AWSConfig.InstanceType = req.MachineType
	config.AWSConfig.VolumeSize = "10"
	volumeSize, err := strconv.ParseInt(config.AWSConfig.VolumeSize, 10, 64)

	if err != nil {
		return errors.Wrapf(err, "parse volume size %s", config.AWSConfig.VolumeSize)
	}

	input := &ec2.RequestSpotInstancesInput{
		Type: aws.String("persistent"),
		LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
			IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
				Name: aws.String(config.AWSConfig.NodesInstanceProfile),
			},
			SubnetId:         aws.String(config.AWSConfig.Subnets[req.AvailabilityZone]),
			SecurityGroupIds: []*string{aws.String(config.AWSConfig.NodesSecurityGroupID)},
			ImageId:          aws.String(config.AWSConfig.ImageID),
			InstanceType:     aws.String(config.AWSConfig.InstanceType),
			KeyName:          aws.String(config.AWSConfig.KeyPairName),
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/sda1"),
					Ebs: &ec2.EbsBlockDevice{
						DeleteOnTermination: aws.Bool(false),
						VolumeType:          aws.String("gp2"),
						VolumeSize:          aws.Int64(volumeSize),
					},
				},
			},
			UserData: aws.String(config.ConfigMap.Data),
		},
		SpotPrice:     aws.String(req.SpotPrice),
		ClientToken:   aws.String(uuid.New()),
		InstanceCount: aws.Int64(1),
		DryRun:        aws.Bool(config.DryRun),
		ValidFrom:     aws.Time(time.Now().Add(time.Second * 10)),
		// TODO(stgleb): pass this as a parameter
		ValidUntil: aws.Time(time.Now().Add(time.Duration(24*365) * time.Hour)),
	}

	result, err := svc.RequestSpotInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				logrus.Errorf("request spot instance caused %s", aerr.Message())
			}
		} else {
			logrus.Errorf("Error %v", err)
		}
		return errors.Wrap(err, "request spot instance")
	}

	go func() {
		requestIds := make([]*string, 0)

		for _, spot := range result.SpotInstanceRequests {
			requestIds = append(requestIds, spot.SpotInstanceRequestId)
		}

		describeReq := &ec2.DescribeSpotInstanceRequestsInput{
			DryRun:                 aws.Bool(false),
			SpotInstanceRequestIds: requestIds,
		}

		err = svc.WaitUntilSpotInstanceRequestFulfilled(describeReq)

		if err != nil {
			logrus.Errorf("wait until request full filled %v", err)
		}

		spotRequests, err := svc.DescribeSpotInstanceRequests(describeReq)

		if err != nil {
			logrus.Errorf("describe spot instance requests %v", err)
		}

		logrus.Debugf("Tag spot instance requests and spot instances")
		for _, instance := range spotRequests.SpotInstanceRequests {

			ec2Tags := []*ec2.Tag{
				{
					Key:   aws.String("KubernetesCluster"),
					Value: aws.String(config.ClusterName),
				},
				{
					Key:   aws.String(clouds.TagClusterID),
					Value: aws.String(config.ClusterID),
				},
				{
					Key: aws.String("Name"),
					Value: aws.String(fmt.Sprintf("%s-node-%s",
						config.ClusterName, uuid.New()[:4])),
				},
				{
					Key:   aws.String("Role"),
					Value: aws.String(util.MakeRole(config.IsMaster)),
				},
			}

			tagInput := &ec2.CreateTagsInput{
				Resources: []*string{},
				Tags:      ec2Tags,
			}

			logrus.Infof("Tag instance %s and request id %s",
				*instance.InstanceId, *instance.SpotInstanceRequestId)
			tagInput.Resources = append(tagInput.Resources, instance.InstanceId)
			tagInput.Resources = append(tagInput.Resources, instance.SpotInstanceRequestId)

			_, err = svc.CreateTags(tagInput)

			if err != nil {
				logrus.Errorf("tagging spot instances %v", err)
			}
		}
	}()

	return nil
}
