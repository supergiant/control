package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateNode creates a Kubernetes minion.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
	m.Name = m.Kube.Name + "-minion-" + util.RandomString(5)

	// For the testzzzz
	if m.Kube.KubernetesVersion == "" {
		m.Kube.KubernetesVersion = "1.5.7"
	}
	// TODO move to init outside of func
	mversion := strings.Split(m.Kube.KubernetesVersion, ".")
	userdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/minion.yaml")
	if err != nil {
		return err
	}
	template, err := template.New("minion_template").Parse(string(userdataTemplate))
	if err != nil {
		return err
	}
	var userdata bytes.Buffer
	if err = template.Execute(&userdata, m); err != nil {
		return err
	}
	encodedUserdata := base64.StdEncoding.EncodeToString(userdata.Bytes())

	ec2S := p.EC2(m.Kube)
	ami, err := getAMI(ec2S)
	if err != nil {
		return err
	}

	var subnets []string
	for _, subnet := range m.Kube.AWSConfig.PublicSubnetIPRange {
		if subnet["subnet_id"] != "" {
			subnets = append(subnets, subnet["subnet_id"])
		}
	}

	var selectedSubnet string
	if len(subnets) == 1 {
		selectedSubnet = subnets[0]
	} else {
		fmt.Println("Number of nodes:", len(m.Kube.Nodes))
		selectedSubnet = subnets[(len(m.Kube.Nodes)-1)%len(m.Kube.AWSConfig.PublicSubnetIPRange)]
	}

	var nodeRole *string
	if m.Kube.AWSConfig.NodeRoleName != "" {
		nodeRole = aws.String(m.Kube.AWSConfig.NodeRoleName)
	} else {
		nodeRole = aws.String("kubernetes-minion")
	}

	resp, err := ec2S.RunInstances(&ec2.RunInstancesInput{
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		InstanceType: aws.String(m.Size),
		ImageId:      aws.String(ami),
		EbsOptimized: aws.Bool(true),
		KeyName:      aws.String(m.Kube.Name + "-key"),
		SecurityGroupIds: []*string{
			aws.String(m.Kube.AWSConfig.NodeSecurityGroupID),
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: nodeRole,
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(int64(m.Kube.AWSConfig.NodeVolumeSize)),
				},
			},
		},
		UserData: aws.String(encodedUserdata),
		SubnetId: aws.String(selectedSubnet),
	})
	if err != nil {
		return err
	}

	server := resp.Instances[0]

	err = tagAWSResource(ec2S, *server.InstanceId, map[string]string{
		"KubernetesCluster": m.Kube.Name,
		"Name":              m.Name,
		"Role":              m.Kube.Name + "-minion",
	}, m.Kube.AWSConfig.Tags)
	if err != nil {
		// TODO
		p.Core.Log.Error("Failed to tag EC2 Instance " + *server.InstanceId)
	}
	m.ProviderID = *server.InstanceId
	m.Name = *server.PrivateDnsName
	m.ProviderCreationTimestamp = time.Now()
	return p.Core.DB.Save(m)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
