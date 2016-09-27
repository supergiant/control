package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/supergiant/guber"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

var AWSMasterAMIs = map[string]string{
	"ap-northeast-1": "ami-907fa690",
	"ap-southeast-1": "ami-b4a79de6",
	"eu-central-1":   "ami-e8635bf5",
	"eu-west-1":      "ami-0fd0ae78",
	"sa-east-1":      "ami-f9f675e4",
	"us-east-1":      "ami-f57b8f9e",
	"us-west-1":      "ami-87b643c3",
	"cn-north-1":     "ami-3abf2203",
	"ap-southeast-2": "ami-1bb9c221",
	"us-west-2":      "ami-33566d03",
}

// TODO this and the similar concept in Kubes should be moved to core, not global vars
var globalAWSSession = session.New()

type Provider struct {
	Core        *core.Core
	Credentials map[string]string
}

func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	// Doesn't really matter what we do here, as long as it works
	_, err := p.ec2("us-east-1").DescribeKeyPairs(new(ec2.DescribeKeyPairsInput))
	return err
}

func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	return p.createKube(m, action)
}

func (p *Provider) DeleteKube(m *model.Kube) error {
	return p.deleteKube(m)
}

func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
	return p.createNode(m)
}

func (p *Provider) DeleteNode(m *model.Node) error {
	return p.deleteServer(m)
}

func (p *Provider) CreateVolume(m *model.Volume, action *core.Action) error {
	return p.createVolume(m, nil)
}

func (p *Provider) KubernetesVolumeDefinition(m *model.Volume) *guber.Volume {
	return &guber.Volume{
		Name: m.Name,
		AwsElasticBlockStore: &guber.AwsElasticBlockStore{
			VolumeID: m.ProviderID,
			FSType:   "ext4",
		},
	}
}

func (p *Provider) ResizeVolume(m *model.Volume, action *core.Action) error {
	return p.resizeVolume(m) // TODO pass action for cancellableWaitFor
}

func (p *Provider) WaitForVolumeAvailable(m *model.Volume, action *core.Action) error {
	return p.waitForAvailable(m)
}

func (p *Provider) DeleteVolume(m *model.Volume) error {
	return p.deleteVolume(m)
}

func (p *Provider) CreateEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return p.createELB(m)
}

func (p *Provider) AddPortToEntrypoint(m *model.Entrypoint, lbPort int64, nodePort int64) error {
	return p.addPortToEntrypoint(m, lbPort, nodePort)
}

func (p *Provider) RemovePortFromEntrypoint(m *model.Entrypoint, lbPort int64) error {
	return p.removePortFromEntrypoint(m, lbPort)
}

func (p *Provider) DeleteEntrypoint(m *model.Entrypoint) error {
	return p.deleteELB(m)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (p *Provider) awsConfig(region string) *aws.Config {
	creds := credentials.NewStaticCredentials(p.Credentials["access_key"], p.Credentials["secret_key"], "")
	creds.Get()
	return aws.NewConfig().WithRegion(region).WithCredentials(creds)
}

func (p *Provider) ec2(region string) *ec2.EC2 {
	return ec2.New(globalAWSSession, p.awsConfig(region))
}

func (p *Provider) iam(region string) *iam.IAM {
	return iam.New(globalAWSSession, p.awsConfig(region))
}

func (p *Provider) elb(region string) *elb.ELB {
	return elb.New(globalAWSSession, p.awsConfig(region))
}

//------------------------------------------------------------------------------

func (p *Provider) createKube(m *model.Kube, action *core.Action) error {
	iamS := p.iam(m.AWSConfig.Region)
	ec2S := p.ec2(m.AWSConfig.Region)
	procedure := &core.Procedure{
		Core:  p.Core,
		Name:  "Create Kube",
		Model: m,
	}

	procedure.AddStep("preparing IAM Role kubernetes-master", func() error {
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Principal": {"AWS": "*"},
					"Effect": "Allow",
					"Sid": ""
				}
			]
		}`
		return createIAMRole(iamS, "kubernetes-master", policy)
	})

	procedure.AddStep("preparing IAM Role Policy kubernetes-master", func() error {
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": ["ec2:*"],
					"Resource": ["*"]
				},
				{
					"Effect": "Allow",
					"Action": ["elasticloadbalancing:*"],
					"Resource": ["*"]
				},
				{
					"Effect": "Allow",
					"Action": "s3:*",
					"Resource": [
						"arn:aws:s3:::kubernetes-*"
					]
				}
			]
		}`
		return createIAMRolePolicy(iamS, "kubernetes-master", policy)
	})

	procedure.AddStep("preparing IAM Instance Profile kubernetes-master", func() error {
		return createIAMInstanceProfile(iamS, "kubernetes-master")
	})

	procedure.AddStep("preparing IAM Role kubernetes-minion", func() error {
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Principal": {"AWS": "*"},
					"Effect": "Allow",
					"Sid": ""
				}
			]
		}`
		return createIAMRole(iamS, "kubernetes-minion", policy)
	})

	procedure.AddStep("preparing IAM Role Policy kubernetes-minion", func() error {
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "s3:*",
					"Resource": [
						"arn:aws:s3:::kubernetes-*"
					]
				},
				{
					"Effect": "Allow",
					"Action": "ec2:Describe*",
					"Resource": "*"
				},
				{
					"Effect": "Allow",
					"Action": "ec2:AttachVolume",
					"Resource": "*"
				},
				{
					"Effect": "Allow",
					"Action": "ec2:DetachVolume",
					"Resource": "*"
				}
			]
		}`
		return createIAMRolePolicy(iamS, "kubernetes-minion", policy)
	})

	procedure.AddStep("preparing IAM Instance Profile kubernetes-minion", func() error {
		return createIAMInstanceProfile(iamS, "kubernetes-minion")
	})

	procedure.AddStep("creating SSH Key Pair", func() error {
		if m.AWSConfig.PrivateKey != "" {
			return nil
		}
		input := &ec2.CreateKeyPairInput{
			KeyName: aws.String(m.Name + "-key"),
		}
		resp, err := ec2S.CreateKeyPair(input)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate") {

				delInput := &ec2.DeleteKeyPairInput{
					KeyName: aws.String(m.Name + "-key"),
				}
				_, _ = ec2S.DeleteKeyPair(delInput)
				return errors.New("KeyPair existed, but key material was not captured. Deleted KeyPair... will retry")

			}
			return err
		}
		m.AWSConfig.PrivateKey = *resp.KeyMaterial
		return nil
	})

	procedure.AddStep("creating VPC", func() error {
		if m.AWSConfig.VPCID != "" {
			return nil
		}
		input := &ec2.CreateVpcInput{
			CidrBlock: aws.String(m.AWSConfig.VPCIPRange),
		}
		resp, err := ec2S.CreateVpc(input)
		if err != nil {
			return err
		}
		m.AWSConfig.VPCID = *resp.Vpc.VpcId
		return nil
	})

	procedure.AddStep("tagging VPC", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.VPCID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-vpc",
		})
	})

	procedure.AddStep("enabling VPC DNS", func() error {
		input := &ec2.ModifyVpcAttributeInput{
			VpcId:              aws.String(m.AWSConfig.VPCID),
			EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		_, err := ec2S.ModifyVpcAttribute(input)
		return err
	})

	// Create Internet Gateway

	procedure.AddStep("creating Internet Gateway", func() error {
		if m.AWSConfig.InternetGatewayID != "" {
			return nil
		}
		resp, err := ec2S.CreateInternetGateway(new(ec2.CreateInternetGatewayInput))
		if err != nil {
			return err
		}
		m.AWSConfig.InternetGatewayID = *resp.InternetGateway.InternetGatewayId
		return nil
	})

	procedure.AddStep("tagging Internet Gateway", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.InternetGatewayID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-ig",
		})
	})

	procedure.AddStep("attaching Internet Gateway to VPC", func() error {
		input := &ec2.AttachInternetGatewayInput{
			VpcId:             aws.String(m.AWSConfig.VPCID),
			InternetGatewayId: aws.String(m.AWSConfig.InternetGatewayID),
		}
		if _, err := ec2S.AttachInternetGateway(input); err != nil && !strings.Contains(err.Error(), "already attached") {
			return err
		}
		return nil
	})

	// Create Subnet

	procedure.AddStep("creating Subnet", func() error {
		if m.AWSConfig.PublicSubnetID != "" {
			return nil
		}
		input := &ec2.CreateSubnetInput{
			VpcId:            aws.String(m.AWSConfig.VPCID),
			CidrBlock:        aws.String(m.AWSConfig.PublicSubnetIPRange),
			AvailabilityZone: aws.String(m.AWSConfig.AvailabilityZone),
		}
		resp, err := ec2S.CreateSubnet(input)
		if err != nil {
			return err
		}
		m.AWSConfig.PublicSubnetID = *resp.Subnet.SubnetId
		return nil
	})

	procedure.AddStep("tagging Subnet", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.PublicSubnetID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-psub",
		})
	})

	procedure.AddStep("enabling public IP assignment setting of Subnet", func() error {
		input := &ec2.ModifySubnetAttributeInput{
			SubnetId:            aws.String(m.AWSConfig.PublicSubnetID),
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		_, err := ec2S.ModifySubnetAttribute(input)
		return err
	})

	// Route Table

	procedure.AddStep("creating Route Table", func() error {
		if m.AWSConfig.RouteTableID != "" {
			return nil
		}
		input := &ec2.CreateRouteTableInput{
			VpcId: aws.String(m.AWSConfig.VPCID),
		}
		resp, err := ec2S.CreateRouteTable(input)
		if err != nil {
			return err
		}
		m.AWSConfig.RouteTableID = *resp.RouteTable.RouteTableId
		return nil
	})

	procedure.AddStep("tagging Route Table", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.RouteTableID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-rt",
		})
	})

	procedure.AddStep("associating Route Table with Subnet", func() error {
		if m.AWSConfig.RouteTableSubnetAssociationID != "" {
			return nil
		}
		input := &ec2.AssociateRouteTableInput{
			RouteTableId: aws.String(m.AWSConfig.RouteTableID),
			SubnetId:     aws.String(m.AWSConfig.PublicSubnetID),
		}
		resp, err := ec2S.AssociateRouteTable(input)
		if err != nil {
			return err
		}
		m.AWSConfig.RouteTableSubnetAssociationID = *resp.AssociationId
		return nil
	})

	procedure.AddStep("creating Route for Internet Gateway", func() error {
		input := &ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			RouteTableId:         aws.String(m.AWSConfig.RouteTableID),
			GatewayId:            aws.String(m.AWSConfig.InternetGatewayID),
		}
		if _, err := ec2S.CreateRoute(input); err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return err
		}
		return nil
	})

	// Create Security Groups

	procedure.AddStep("creating ELB Security Group", func() error {
		if m.AWSConfig.ELBSecurityGroupID != "" {
			return nil
		}
		input := &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String(m.Name + "_elb_sg"),
			Description: aws.String("Allow any external port through to internal 30-40k range"),
			VpcId:       aws.String(m.AWSConfig.VPCID),
		}
		resp, err := ec2S.CreateSecurityGroup(input)
		if err != nil {
			return err
		}
		m.AWSConfig.ELBSecurityGroupID = *resp.GroupId
		return nil
	})

	procedure.AddStep("tagging ELB Security Group", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.ELBSecurityGroupID, map[string]string{
			"KubernetesCluster": m.Name,
		})
	})

	procedure.AddStep("creating ELB Security Group ingress rules", func() error {
		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(m.AWSConfig.ELBSecurityGroupID),
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(0),
					ToPort:     aws.Int64(0),
					IpProtocol: aws.String("-1"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		}
		if _, err := ec2S.AuthorizeSecurityGroupIngress(input); err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return err
		}
		return nil
	})

	procedure.AddStep("creating ELB Security Group egress rules", func() error {
		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(m.AWSConfig.ELBSecurityGroupID),
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(30000),
					ToPort:     aws.Int64(40000),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				{
					FromPort:   aws.Int64(10250),
					ToPort:     aws.Int64(10250),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		}
		if _, err := ec2S.AuthorizeSecurityGroupIngress(input); err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return err
		}
		return nil
	})

	procedure.AddStep("creating Node Security Group", func() error {
		if m.AWSConfig.NodeSecurityGroupID != "" {
			return nil
		}
		input := &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String(m.Name + "_sg"),
			Description: aws.String("Allow any traffic to 443 and 22, but only traffic from ELB for 10250 and 30k-40k"),
			VpcId:       aws.String(m.AWSConfig.VPCID),
		}
		resp, err := ec2S.CreateSecurityGroup(input)
		if err != nil {
			return err
		}
		m.AWSConfig.NodeSecurityGroupID = *resp.GroupId
		return nil
	})

	procedure.AddStep("tagging Node Security Group", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.NodeSecurityGroupID, map[string]string{
			"KubernetesCluster": m.Name,
		})
	})

	procedure.AddStep("creating Node Security Group ingress rules", func() error {
		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(m.AWSConfig.NodeSecurityGroupID),
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(0),
					ToPort:     aws.Int64(0),
					IpProtocol: aws.String("-1"),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: aws.String(m.AWSConfig.NodeSecurityGroupID), // ?? TODO is this correct? -- https://github.com/supergiant/terraform-assets/blob/master/aws/1.1.7/security_groups.tf#L39
						},
					},
				},
				{
					FromPort:   aws.Int64(22),
					ToPort:     aws.Int64(22),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				{
					FromPort:   aws.Int64(443),
					ToPort:     aws.Int64(443),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				{
					FromPort:   aws.Int64(30000),
					ToPort:     aws.Int64(40000),
					IpProtocol: aws.String("tcp"),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: aws.String(m.AWSConfig.ELBSecurityGroupID),
						},
					},
				},
				{
					FromPort:   aws.Int64(10250),
					ToPort:     aws.Int64(10250),
					IpProtocol: aws.String("tcp"),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: aws.String(m.AWSConfig.ELBSecurityGroupID),
						},
					},
				},
			},
		}
		if _, err := ec2S.AuthorizeSecurityGroupIngress(input); err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return err
		}
		return nil
	})

	procedure.AddStep("creating Node Security Group egress rules", func() error {
		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(m.AWSConfig.NodeSecurityGroupID),
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(0),
					ToPort:     aws.Int64(0),
					IpProtocol: aws.String("-1"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		}
		if _, err := ec2S.AuthorizeSecurityGroupIngress(input); err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return err
		}
		return nil
	})

	// Master Instance

	procedure.AddStep("creating Server for Kubernetes master", func() error {
		if m.AWSConfig.MasterID != "" {
			return nil
		}

		userdataTemplate, err := Asset("config/providers/aws/master_userdata.txt")
		if err != nil {
			return err
		}
		template, err := template.New("master_template").Parse(string(userdataTemplate))
		if err != nil {
			return err
		}
		var userdata bytes.Buffer
		if err = template.Execute(&userdata, m); err != nil {
			return err
		}
		encodedUserdata := base64.StdEncoding.EncodeToString(userdata.Bytes())

		input := &ec2.RunInstancesInput{
			MinCount:     aws.Int64(1),
			MaxCount:     aws.Int64(1),
			ImageId:      aws.String(AWSMasterAMIs[m.AWSConfig.Region]),
			InstanceType: aws.String(m.MasterNodeSize),
			KeyName:      aws.String(m.Name + "-key"),
			NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
				{
					DeviceIndex:              aws.Int64(0),
					AssociatePublicIpAddress: aws.Bool(true),
					DeleteOnTermination:      aws.Bool(true),
					Groups: []*string{
						aws.String(m.AWSConfig.NodeSecurityGroupID),
					},
					SubnetId:         aws.String(m.AWSConfig.PublicSubnetID),
					PrivateIpAddress: aws.String(m.AWSConfig.MasterPrivateIP),
				},
			},
			IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
				Name: aws.String("kubernetes-master"),
			},
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/xvdb"),
					Ebs: &ec2.EbsBlockDevice{
						DeleteOnTermination: aws.Bool(true),
						VolumeType:          aws.String("gp2"),
						VolumeSize:          aws.Int64(20),
					},
				},
			},
			UserData: aws.String(encodedUserdata),
		}
		resp, err := ec2S.RunInstances(input)
		if err != nil {
			return err
		}

		instance := resp.Instances[0]

		m.AWSConfig.MasterID = *instance.InstanceId
		return nil
	})

	procedure.AddStep("tagging Kubernetes master", func() error {
		return tagAWSResource(ec2S, m.AWSConfig.MasterID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-master",
			"Role":              m.Name + "-master",
		})
	})

	// Wait for server to be ready

	procedure.AddStep("waiting for Kubernetes master to launch", func() error {
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(m.AWSConfig.MasterID),
			},
		}

		return action.CancellableWaitFor("Kubernetes master launch", 5*time.Minute, 3*time.Second, func() (bool, error) {
			resp, err := ec2S.DescribeInstances(input)
			if err != nil {
				return false, err
			}

			instance := resp.Reservations[0].Instances[0]

			// Save IP when ready
			if m.MasterPublicIP == "" {
				if ip := instance.PublicIpAddress; ip != nil {
					m.MasterPublicIP = *ip
					if err := p.Core.DB.Save(m); err != nil {
						return false, err
					}
				}
			}

			return *instance.State.Name == "running", nil
		})
	})

	// Create route for master

	procedure.AddStep("creating Route for Kubernetes master", func() error {
		input := &ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("10.246.0.0/24"),
			RouteTableId:         aws.String(m.AWSConfig.RouteTableID),
			InstanceId:           aws.String(m.AWSConfig.MasterID),
		}
		_, err := ec2S.CreateRoute(input)
		return err
	})

	// Create first minion

	procedure.AddStep("creating Kubernetes minion", func() error {
		// TODO repeated in DO provider
		if err := p.Core.DB.Model(m).Association("Nodes").Find(&m.Nodes).Error; err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		}
		node := &model.Node{
			KubeID: m.ID,
			Size:   m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})

	procedure.AddStep("waiting for Kubernetes", func() error {
		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, 3*time.Second, func() (bool, error) {
			nodes, err := p.Core.K8S(m).Nodes().List()
			if err != nil {
				return false, nil
			}
			return len(nodes.Items) > 0, nil
		})
	})

	return procedure.Run()
}

func (p *Provider) deleteKube(m *model.Kube) error {
	ec2S := p.ec2(m.AWSConfig.Region)
	procedure := &core.Procedure{
		Core:  p.Core,
		Name:  "Delete Kube",
		Model: m,
	}

	procedure.AddStep("deleting master", func() error {
		if m.AWSConfig.MasterID == "" {
			return nil
		}

		input := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(m.AWSConfig.MasterID),
			},
		}
		if _, err := ec2S.TerminateInstances(input); isErrAndNotAWSNotFound(err) {
			return err
		}

		// Wait for termination
		descinput := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(m.AWSConfig.MasterID),
			},
		}
		waitErr := util.WaitFor("Kubernetes master termination", 5*time.Minute, 3*time.Second, func() (bool, error) { // TODO --------- use server() method
			resp, err := ec2S.DescribeInstances(descinput)
			if err != nil {
				return false, err
			}
			if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
				return true, nil
			}
			instance := resp.Reservations[0].Instances[0]
			return *instance.State.Name == "terminated", nil
		})
		// Done waiting
		if waitErr != nil {
			return waitErr
		}

		m.AWSConfig.MasterID = ""
		return nil
	})

	procedure.AddStep("disassociating Route Table from Subnet", func() error {
		if m.AWSConfig.RouteTableSubnetAssociationID == "" {
			return nil
		}
		input := &ec2.DisassociateRouteTableInput{
			AssociationId: aws.String(m.AWSConfig.RouteTableSubnetAssociationID),
		}
		if _, err := ec2S.DisassociateRouteTable(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.RouteTableSubnetAssociationID = ""
		return nil
	})

	procedure.AddStep("deleting Internet Gateway", func() error {
		if m.AWSConfig.InternetGatewayID == "" {
			return nil
		}
		diginput := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: aws.String(m.AWSConfig.InternetGatewayID),
			VpcId:             aws.String(m.AWSConfig.VPCID),
		}

		// NOTE we do this (maybe we should just describe, not spam detach) because
		// we can't wait directly on minions to terminate (we can, but I'm lazy rn)
		waitErr := util.WaitFor("Internet Gateway to detach", 5*time.Minute, 5*time.Second, func() (bool, error) {
			if _, err := ec2S.DetachInternetGateway(diginput); err != nil && !strings.Contains(err.Error(), "not attached") {

				p.Core.Log.Warn(err.Error())

				return false, nil
			}
			return true, nil
		})
		if waitErr != nil {
			return waitErr
		}

		input := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: aws.String(m.AWSConfig.InternetGatewayID),
		}
		if _, err := ec2S.DeleteInternetGateway(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.InternetGatewayID = ""
		return nil
	})

	procedure.AddStep("deleting Route Table", func() error {
		if m.AWSConfig.RouteTableID == "" {
			return nil
		}
		input := &ec2.DeleteRouteTableInput{
			RouteTableId: aws.String(m.AWSConfig.RouteTableID),
		}
		if _, err := ec2S.DeleteRouteTable(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.RouteTableID = ""
		return nil
	})

	procedure.AddStep("deleting public Subnet", func() error {
		if m.AWSConfig.PublicSubnetID == "" {
			return nil
		}
		input := &ec2.DeleteSubnetInput{
			SubnetId: aws.String(m.AWSConfig.PublicSubnetID),
		}

		waitErr := util.WaitFor("Public Subnet to delete", 2*time.Minute, 5*time.Second, func() (bool, error) {
			if _, err := ec2S.DeleteSubnet(input); isErrAndNotAWSNotFound(err) {
				return false, nil
			}
			return true, nil
		})
		if waitErr != nil {
			return waitErr
		}

		m.AWSConfig.PublicSubnetID = ""
		return nil
	})

	procedure.AddStep("deleting Node Security Group", func() error {
		if m.AWSConfig.NodeSecurityGroupID == "" {
			return nil
		}
		input := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(m.AWSConfig.NodeSecurityGroupID),
		}
		if _, err := ec2S.DeleteSecurityGroup(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.NodeSecurityGroupID = ""
		return nil
	})

	procedure.AddStep("deleting ELB Security Group", func() error {
		if m.AWSConfig.ELBSecurityGroupID == "" {
			return nil
		}
		input := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(m.AWSConfig.ELBSecurityGroupID),
		}
		if _, err := ec2S.DeleteSecurityGroup(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.ELBSecurityGroupID = ""
		return nil
	})

	procedure.AddStep("deleting VPC", func() error {
		if m.AWSConfig.VPCID == "" {
			return nil
		}
		input := &ec2.DeleteVpcInput{
			VpcId: aws.String(m.AWSConfig.VPCID),
		}
		if _, err := ec2S.DeleteVpc(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		m.AWSConfig.VPCID = ""
		return nil
	})

	procedure.AddStep("deleting SSH Key Pair", func() error {
		input := &ec2.DeleteKeyPairInput{
			KeyName: aws.String(m.Name + "-key"),
		}
		if _, err := ec2S.DeleteKeyPair(input); isErrAndNotAWSNotFound(err) {
			return err
		}
		return nil
	})

	return procedure.Run()
}

func (p *Provider) createNode(m *model.Node) error {
	server, err := p.createServer(m)
	if err != nil {
		return err
	}
	p.setAttrsFromServer(m, server)
	if err := p.Core.DB.Save(m); err != nil {
		return err
	}
	for _, entrypoint := range m.Kube.Entrypoints {
		if err := p.registerNodes(entrypoint, m); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) servers(m *model.Kube) (instances []*ec2.Instance, err error) {
	return p.filteredServers(m, map[string][]string{
		"tag:Name": []string{
			m.Name + "-minion",
		},
	})
}

func (p *Provider) server(m *model.Kube, id string) (*ec2.Instance, error) {
	instances, err := p.filteredServers(m, map[string][]string{
		"instance-id": []string{id},
	})
	if err != nil {
		return nil, err
	}
	return instances[0], nil
}

func (p *Provider) filteredServers(m *model.Kube, filters map[string][]string) (instances []*ec2.Instance, err error) {
	filters["instance-state-name"] = []string{
		"pending",
		"running",
	}
	var ec2Filters []*ec2.Filter
	for key, vals := range filters {
		ec2Filter := &ec2.Filter{
			Name: aws.String(key),
		}
		for _, val := range vals {
			ec2Filter.Values = append(ec2Filter.Values, aws.String(val))
		}
		ec2Filters = append(ec2Filters, ec2Filter)
	}
	input := &ec2.DescribeInstancesInput{
		Filters: ec2Filters,
	}
	resp, err := p.ec2(m.AWSConfig.Region).DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}
	return instances, nil
}

func (p *Provider) setAttrsFromServer(m *model.Node, server *ec2.Instance) {
	m.ProviderID = *server.InstanceId
	m.Name = *server.PrivateDnsName
	m.Size = *server.InstanceType
	m.ProviderCreationTimestamp = *server.LaunchTime
}

func (p *Provider) createServer(m *model.Node) (*ec2.Instance, error) {

	// TODO move to init outside of func
	userdataTemplate, err := Asset("config/providers/aws/minion_userdata.txt")
	if err != nil {
		return nil, err
	}
	template, err := template.New("minion_template").Parse(string(userdataTemplate))
	if err != nil {
		return nil, err
	}
	var userdata bytes.Buffer
	if err = template.Execute(&userdata, m.Kube); err != nil {
		return nil, err
	}
	encodedUserdata := base64.StdEncoding.EncodeToString(userdata.Bytes())

	input := &ec2.RunInstancesInput{
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		InstanceType: aws.String(m.Size),
		ImageId:      aws.String(AWSMasterAMIs[m.Kube.AWSConfig.Region]),
		EbsOptimized: aws.Bool(true),
		KeyName:      aws.String(m.Kube.Name + "-key"),
		SecurityGroupIds: []*string{
			aws.String(m.Kube.AWSConfig.NodeSecurityGroupID),
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("kubernetes-minion"),
		},
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			// root device ? TODO (do we need the one after this?)
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(80),
					DeleteOnTermination: aws.Bool(true),
				},
			},
			&ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/xvdb"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType:          aws.String("gp2"),
					VolumeSize:          aws.Int64(80),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
		UserData: aws.String(encodedUserdata),
		SubnetId: aws.String(m.Kube.AWSConfig.PublicSubnetID),
	}

	ec2S := p.ec2(m.Kube.AWSConfig.Region)

	resp, err := ec2S.RunInstances(input)
	if err != nil {
		return nil, err
	}

	server := resp.Instances[0]

	err = tagAWSResource(ec2S, *server.InstanceId, map[string]string{
		"KubernetesCluster": m.Kube.Name,
		"Name":              m.Kube.Name + "-minion",
		"Role":              m.Kube.Name + "-minion",
	})
	if err != nil {
		// TODO
		p.Core.Log.Error("Failed to tag EC2 Instance " + *server.InstanceId)
	}

	return server, nil
}

func (p *Provider) deleteServer(m *model.Node) error {

	// TODO move out of here
	if m.Kube == nil {
		p.Core.Log.Warnf("Deleting Node %d without deleting server because Kube is nil", *m.ID)
		return nil
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(m.ProviderID)},
	}
	_, err := p.ec2(m.Kube.AWSConfig.Region).TerminateInstances(input)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) createELB(m *model.Entrypoint) error {
	params := &elb.CreateLoadBalancerInput{
		Listeners: []*elb.Listener{ // NOTE we must provide at least 1 listener, it is currently arbitrary
			{
				InstancePort:     aws.Int64(420),
				LoadBalancerPort: aws.Int64(420),
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: aws.String(m.ProviderID),
		Scheme:           aws.String("internet-facing"),
		SecurityGroups: []*string{
			aws.String(m.Kube.AWSConfig.ELBSecurityGroupID),
		},
		Subnets: []*string{
			aws.String(m.Kube.AWSConfig.PublicSubnetID),
		},
	}
	resp, err := p.elb(m.Kube.AWSConfig.Region).CreateLoadBalancer(params)
	if err != nil {
		return err
	}

	// Save Address
	m.Address = *resp.DNSName
	if err := p.Core.DB.Save(m); err != nil {
		return err
	}

	if err := p.registerNodes(m, m.Kube.Nodes...); err != nil {
		return err
	}

	// Configure health check
	healthParams := &elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(m.ProviderID),
		HealthCheck: &elb.HealthCheck{
			Target:             aws.String("HTTPS:10250/healthz"),
			HealthyThreshold:   aws.Int64(2),
			UnhealthyThreshold: aws.Int64(10),
			Interval:           aws.Int64(30),
			Timeout:            aws.Int64(5),
		},
	}
	_, err = p.elb(m.Kube.AWSConfig.Region).ConfigureHealthCheck(healthParams)
	return err
}

func (p *Provider) addPortToEntrypoint(m *model.Entrypoint, lbPort int64, nodePort int64) error {
	input := &elb.CreateLoadBalancerListenersInput{
		LoadBalancerName: aws.String(m.ProviderID),
		Listeners: []*elb.Listener{
			{
				LoadBalancerPort: aws.Int64(lbPort),
				InstancePort:     aws.Int64(nodePort),
				Protocol:         aws.String("TCP"),
			},
		},
	}
	_, err := p.elb(m.Kube.AWSConfig.Region).CreateLoadBalancerListeners(input)
	return err
}

func (p *Provider) removePortFromEntrypoint(m *model.Entrypoint, lbPort int64) error {
	params := &elb.DeleteLoadBalancerListenersInput{
		LoadBalancerName: aws.String(m.ProviderID),
		LoadBalancerPorts: []*int64{
			aws.Int64(lbPort),
		},
	}
	_, err := p.elb(m.Kube.AWSConfig.Region).DeleteLoadBalancerListeners(params)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) registerNodes(m *model.Entrypoint, nodes ...*model.Node) error {
	var elbInstances []*elb.Instance
	for _, node := range nodes {
		elbInstances = append(elbInstances, &elb.Instance{
			InstanceId: aws.String(node.ProviderID),
		})
	}
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(m.ProviderID),
		Instances:        elbInstances,
	}
	_, err := p.elb(m.Kube.AWSConfig.Region).RegisterInstancesWithLoadBalancer(input)
	return err
}

func (p *Provider) deleteELB(m *model.Entrypoint) error {
	// Delete ELB
	params := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(m.ProviderID),
	}
	_, err := p.elb(m.Kube.AWSConfig.Region).DeleteLoadBalancer(params)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) createVolume(volume *model.Volume, snapshotID *string) error {
	volInput := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(volume.Kube.AWSConfig.AvailabilityZone),
		VolumeType:       aws.String(volume.Type),
		Size:             aws.Int64(int64(volume.Size)),
		SnapshotId:       snapshotID,
	}
	awsVol, err := p.ec2(volume.Kube.AWSConfig.Region).CreateVolume(volInput)
	if err != nil {
		return err
	}

	volume.ProviderID = *awsVol.VolumeId
	volume.Size = int(*awsVol.Size)
	if err := p.Core.DB.Save(volume); err != nil {
		return err
	}

	tagsInput := &ec2.CreateTagsInput{
		Resources: []*string{
			awsVol.VolumeId,
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(volume.Name),
			},
		},
	}
	_, err = p.ec2(volume.Kube.AWSConfig.Region).CreateTags(tagsInput)
	return err
}

func (p *Provider) resizeVolume(m *model.Volume) error {
	snapshot, err := p.createSnapshot(m)
	if err != nil {
		return err
	}
	if err := p.deleteVolume(m); err != nil {
		return err
	}
	if err := p.createVolume(m, snapshot.SnapshotId); err != nil {
		return err
	}
	if err := p.deleteSnapshot(m, snapshot); err != nil {
		p.Core.Log.Errorf("Error deleting snapshot %s: %s", *snapshot.SnapshotId, err.Error())
	}
	return nil
}

func (p *Provider) deleteVolume(volume *model.Volume) error {
	if err := p.waitForAvailable(volume); err != nil {
		return err
	}
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volume.ProviderID),
	}
	if _, err := p.ec2(volume.Kube.AWSConfig.Region).DeleteVolume(input); isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) waitForAvailable(volume *model.Volume) error {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(volume.ProviderID),
				},
			},
		},
	}

	resp, err := p.ec2(volume.Kube.AWSConfig.Region).DescribeVolumes(input)
	if err != nil {
		return err
	}
	if len(resp.Volumes) == 0 {
		return nil
	}

	p.Core.Log.Debugf("Waiting for EBS volume %s to be available", volume.Name)
	return p.ec2(volume.Kube.AWSConfig.Region).WaitUntilVolumeAvailable(input)
}

func (p *Provider) forceDetachVolume(volume *model.Volume) error {
	input := &ec2.DetachVolumeInput{
		VolumeId: aws.String(volume.ProviderID),
		Force:    aws.Bool(true),
	}
	if _, err := p.ec2(volume.Kube.AWSConfig.Region).DetachVolume(input); isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) createSnapshot(volume *model.Volume) (*ec2.Snapshot, error) {
	input := &ec2.CreateSnapshotInput{
		Description: aws.String(fmt.Sprintf("%s-%d", volume.Name, volume.Instance.ReleaseID)),
		VolumeId:    aws.String(volume.ProviderID),
	}
	snapshot, err := p.ec2(volume.Kube.AWSConfig.Region).CreateSnapshot(input)
	if err != nil {
		return nil, err
	}
	waitInput := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{snapshot.SnapshotId},
	}
	//
	// TODO
	//
	// We should use action.CancellableWaitFor here
	//
	if err := p.ec2(volume.Kube.AWSConfig.Region).WaitUntilSnapshotCompleted(waitInput); err != nil {
		return nil, err // TODO destroy snapshot that failed to complete?
	}
	return snapshot, nil
}

func (p *Provider) deleteSnapshot(volume *model.Volume, snapshot *ec2.Snapshot) error {
	input := &ec2.DeleteSnapshotInput{
		SnapshotId: snapshot.SnapshotId,
	}
	if _, err := p.ec2(volume.Kube.AWSConfig.Region).DeleteSnapshot(input); isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

// is it NOT Not Found
func isErrAndNotAWSNotFound(err error) bool {
	return err != nil && !regexp.MustCompile(`([Nn]ot *[Ff]ound|404)`).MatchString(err.Error())
}

func createIAMRole(iamS *iam.IAM, name string, policy string) error {
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	_, err := iamS.GetRole(getInput)
	if err == nil {
		return nil
	} else if isErrAndNotAWSNotFound(err) {
		return err
	}
	input := &iam.CreateRoleInput{
		RoleName: aws.String(name),
		Path:     aws.String("/"),
		AssumeRolePolicyDocument: aws.String(policy),
	}
	_, err = iamS.CreateRole(input)
	return err
}

func createIAMRolePolicy(iamS *iam.IAM, name string, policy string) error {
	getInput := &iam.GetRolePolicyInput{
		RoleName:   aws.String(name),
		PolicyName: aws.String(name),
	}
	_, err := iamS.GetRolePolicy(getInput)
	if err == nil {
		return nil
	} else if isErrAndNotAWSNotFound(err) {
		return err
	}

	putRoleInput := &iam.PutRolePolicyInput{
		RoleName:       aws.String(name),
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(policy),
	}
	_, err = iamS.PutRolePolicy(putRoleInput)
	return err
}

func createIAMInstanceProfile(iamS *iam.IAM, name string) error {
	getInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	var instanceProfile *iam.InstanceProfile

	resp, err := iamS.GetInstanceProfile(getInput)
	if err != nil {
		if isErrAndNotAWSNotFound(err) {
			return err
		}

		// Create
		input := &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Path:                aws.String("/"),
		}
		createResp, createErr := iamS.CreateInstanceProfile(input)
		if createErr != nil {
			return createErr
		}
		instanceProfile = createResp.InstanceProfile

	} else {
		instanceProfile = resp.InstanceProfile
	}

	if len(instanceProfile.Roles) == 0 {
		addInput := &iam.AddRoleToInstanceProfileInput{
			RoleName:            aws.String(name),
			InstanceProfileName: aws.String(name),
		}
		if _, err = iamS.AddRoleToInstanceProfile(addInput); err != nil {
			return err
		}
	}

	return nil
}

func tagAWSResource(ec2S *ec2.EC2, idstr string, tags map[string]string) error {
	var ec2Tags []*ec2.Tag
	for key, val := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(val),
		})
	}
	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(idstr)},
		Tags:      ec2Tags,
	}
	_, err := ec2S.CreateTags(input)
	return err
}
