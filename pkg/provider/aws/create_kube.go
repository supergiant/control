package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// CreateKube creates a Kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	iamS := p.IAM(m)
	ec2S := p.EC2(m)
	s3S := p.S3(m)
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	// mock a fake ssh key if the user does not enter one. CoreOS may barf if we don't.
	// Don't worry. This key is a example key used in github doc.
	if m.AWSConfig.SSHPubKey == "" {
		m.AWSConfig.SSHPubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== schacon@mylaptop.local"
	}

	// Init our AZ setup
	// If zone configurtion is empty, the user did not pre-configure, so we need to set it up.
	if len(m.AWSConfig.PublicSubnetIPRange) == 0 && m.AWSConfig.MultiAZ {
		zones := map[string]string{
			m.AWSConfig.Region + "a": "172.20.0.0/24",
			m.AWSConfig.Region + "b": "172.20.1.0/24",
			m.AWSConfig.Region + "c": "172.20.2.0/24",
			m.AWSConfig.Region + "d": "172.20.3.0/24",
			m.AWSConfig.Region + "e": "172.20.4.0/24",
		}

		for zone, ipRange := range zones {
			subnetObj := map[string]string{
				"zone":      zone,
				"ip_range":  ipRange,
				"subnet_id": "",
			}

			m.AWSConfig.PublicSubnetIPRange = append(m.AWSConfig.PublicSubnetIPRange, subnetObj)
		}
	} else if len(m.AWSConfig.PublicSubnetIPRange) == 0 && !m.AWSConfig.MultiAZ {
		zones := map[string]string{
			m.AWSConfig.AvailabilityZone: "172.20.0.0/24",
		}

		for zone, ipRange := range zones {
			subnetObj := map[string]string{
				"zone":      zone,
				"ip_range":  ipRange,
				"subnet_id": "",
			}

			m.AWSConfig.PublicSubnetIPRange = append(m.AWSConfig.PublicSubnetIPRange, subnetObj)
		}
	}

	m.AWSConfig.AvailabilityZone = m.AWSConfig.PublicSubnetIPRange[0]["zone"]

	err := p.Core.DB.Save(m)
	if err != nil {
		return err
	}

	procedure.AddStep("preparing IAM Role kubernetes-master", func() error {
		policy := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": { "Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
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
      "Action": ["route53:*"],
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
      "Effect": "Allow",
      "Principal": { "Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
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
    },
    {
      "Effect": "Allow",
      "Action": ["route53:*"],
      "Resource": ["*"]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:DescribeRepositories",
        "ecr:ListImages",
        "ecr:BatchGetImage"
      ],
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

	procedure.AddStep("creating S3 bucket", func() error {
		_, err := s3S.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String("kubernetes-" + m.Name),
		})
		if err != nil {
			return err
		}

		s3S.PutBucketPolicy(&s3.PutBucketPolicyInput{
			Bucket: aws.String("kubernetes-" + m.Name), // Required
			Policy: aws.String(``),                     // Required
		})

		return nil
	})

	procedure.AddStep("upload assets to S3", func() error {
		// fetch the etcd clustering token
		if m.AWSConfig.KubeMasterCount == 0 {
			m.AWSConfig.KubeMasterCount = 1
		}
		url, err := etcdToken(strconv.Itoa(m.AWSConfig.KubeMasterCount))
		if err != nil {
			return err
		}

		m.AWSConfig.ETCDDiscoveryURL = url

		userdataTemplate, err := bindata.Asset("config/providers/aws/master.yaml")
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

		_, err = s3S.PutObject(&s3.PutObjectInput{
			Bucket:        aws.String("kubernetes-" + m.Name), // Required
			Key:           aws.String("build/master.yaml"),    // Required
			Body:          bytes.NewReader(userdata.Bytes()),
			ContentLength: aws.Int64(int64(userdata.Len())),
		})
		if err != nil {
			return err
		}
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

		for idx, subnet := range m.AWSConfig.PublicSubnetIPRange {

			resp, err := ec2S.CreateSubnet(&ec2.CreateSubnetInput{
				VpcId:            aws.String(m.AWSConfig.VPCID),
				CidrBlock:        aws.String(subnet["ip_range"]),
				AvailabilityZone: aws.String(subnet["zone"]),
			})
			if err != nil {
				//Can we build here?
				if strings.Contains(err.Error(), "Subnets can currently only be created in the following") {
					continue
				}
				return err
			}

			m.AWSConfig.PublicSubnetIPRange[idx]["subnet_id"] = *resp.Subnet.SubnetId

		}

		err := p.Core.DB.Save(m)
		if err != nil {
			return err
		}

		return nil
	})

	procedure.AddStep("tagging Subnet", func() error {

		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				err := tagAWSResource(ec2S, subnet["subnet_id"], map[string]string{
					"KubernetesCluster": m.Name,
					"Name":              m.Name + "-" + subnet["zone"] + "-psub",
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	procedure.AddStep("enabling public IP assignment setting of Subnet", func() error {

		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				_, err := ec2S.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					SubnetId:            aws.String(subnet["subnet_id"]),
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	// Route Table

	procedure.AddStep("creating Route Table", func() error {
		if m.AWSConfig.RouteTableID != "" {
			return nil
		}

		resp, err := ec2S.CreateRouteTable(&ec2.CreateRouteTableInput{
			VpcId: aws.String(m.AWSConfig.VPCID),
		})
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
		if len(m.AWSConfig.RouteTableSubnetAssociationID) != 0 {
			return nil
		}

		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				resp, err := ec2S.AssociateRouteTable(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String(m.AWSConfig.RouteTableID),
					SubnetId:     aws.String(subnet["subnet_id"]),
				})
				if err != nil {
					return err
				}
				m.AWSConfig.RouteTableSubnetAssociationID = append(m.AWSConfig.RouteTableSubnetAssociationID, *resp.AssociationId)
			}
		}

		err := p.Core.DB.Save(m)
		if err != nil {
			return err
		}

		return nil
	})

	procedure.AddStep("creating Route for Internet Gateway", func() error {
		_, err := ec2S.CreateRoute(&ec2.CreateRouteInput{
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			RouteTableId:         aws.String(m.AWSConfig.RouteTableID),
			GatewayId:            aws.String(m.AWSConfig.InternetGatewayID),
		})
		if err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
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

	procedure.AddStep("creating Server for Kubernetes master(s)", func() error {
		if m.AWSConfig.MasterID != "" {
			return nil
		}

		userdataTemplate, err := bindata.Asset("config/providers/aws/bootstrap.yaml")
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

		ami, err := getAMI(ec2S)
		if err != nil {
			return err
		}

		if m.AWSConfig.KubeMasterCount == 0 {
			m.AWSConfig.KubeMasterCount = 1
		}

		masterCount := m.AWSConfig.KubeMasterCount
		loopCount := 0

		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {

			if subnet["subnet_id"] != "" {
				if loopCount >= masterCount {
					break
				}

				time.Sleep(5 * time.Second)

				resp, err := ec2S.RunInstances(&ec2.RunInstancesInput{
					MinCount:     aws.Int64(1),
					MaxCount:     aws.Int64(1),
					ImageId:      aws.String(ami),
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
							SubnetId: aws.String(subnet["subnet_id"]),
							//PrivateIpAddress: aws.String(m.AWSConfig.MasterPrivateIP),
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
				})
				if err != nil {
					return err
				}

				m.AWSConfig.MasterNodes = append(m.AWSConfig.MasterNodes, *resp.Instances[0].InstanceId)
				loopCount++

			} else {
				loopCount++
				continue
			}

		}

		return nil
	})

	// Tag all of our masters.
	procedure.AddStep("tagging Kubernetes master", func() error {
		for _, master := range m.AWSConfig.MasterNodes {
			err := tagAWSResource(ec2S, master, map[string]string{
				"KubernetesCluster": m.Name,
				"Name":              m.Name + "-master",
				"Role":              m.Name + "-master",
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	// If we have more then one master... Lets get them on a internal loadbalancer.
	procedure.AddStep("creating master loadbalancer if needed", func() error {

		if m.AWSConfig.KubeMasterCount == 1 {
			return nil
		}

		var subnets []*string
		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				subnets = append(subnets, aws.String(subnet["subnet_id"]))
			}
		}

		loadbalancer, err := p.ELB(m).CreateLoadBalancer(&elb.CreateLoadBalancerInput{
			Listeners: []*elb.Listener{ // NOTE we must provide at least 1 listener, it is currently arbitrary
				{
					InstancePort:     aws.Int64(443),
					LoadBalancerPort: aws.Int64(443),
					Protocol:         aws.String("TCP"),
				},
				{
					InstancePort:     aws.Int64(2379),
					LoadBalancerPort: aws.Int64(2379),
					Protocol:         aws.String("TCP"),
				},
			},
			LoadBalancerName: aws.String(m.Name + "-api"),
			Scheme:           aws.String("internal"),
			SecurityGroups: []*string{
				aws.String(m.AWSConfig.NodeSecurityGroupID),
			},
			Subnets: subnets,
		})
		if err != nil {
			return err
		}

		m.AWSConfig.MasterPrivateIP = *loadbalancer.DNSName

		var elbInstances []*elb.Instance
		for _, node := range m.AWSConfig.MasterNodes {
			elbInstances = append(elbInstances, &elb.Instance{
				InstanceId: aws.String(node),
			})
		}

		_, err = p.ELB(m).RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{
			LoadBalancerName: aws.String(m.Name + "-api"),
			Instances:        elbInstances,
		})
		if err != nil {
			return err
		}
		return nil
	})

	// Wait for server to be ready

	procedure.AddStep("waiting for Kubernetes master to launch", func() error {
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(m.AWSConfig.MasterNodes[0]),
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
					if m.AWSConfig.MasterPrivateIP == "" {
						m.AWSConfig.MasterPrivateIP = *instance.PrivateIpAddress
					}
					if err := p.Core.DB.Save(m); err != nil {
						return false, err
					}
				}
			}

			return *instance.State.Name == "running", nil
		})
	})

	// Create first minion//
	procedure.AddStep("creating Kubernetes minion", func() error {
		// TODO repeated in DO provider
		if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		} //
		node := &model.Node{
			KubeName: m.Name,
			Size:     m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})

	procedure.AddStep("waiting for Kubernetes", func() error {

		k8s := p.Core.K8S(m)

		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, time.Second, func() (bool, error) {
			nodes, err := k8s.ListNodes("")
			if err != nil {
				return false, nil
			}
			return len(nodes) > 0, nil
		})
	})

	return procedure.Run()
}
