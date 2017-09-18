package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateKube creates a Kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	iamS := p.IAM(m)
	ec2S := p.EC2(m)
	s3S := p.S3(m)
	efsS := p.EFS(m)
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	// mock a fake ssh key if the user does not enter one. CoreOS may barf if we don't.
	// Don't worry. This key is a example key used in github doc.
	if m.SSHPubKey == "" {
		m.SSHPubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== schacon@mylaptop.local"
	}

	m.KubeProviderString = `
         --cloud-provider=aws \`

	m.ProviderString = `
          - --cloud-provider=aws`

	if m.KubernetesVersion == "" {
		m.KubernetesVersion = "1.5.7"
	}

	// Init our AZ setup
	// If zone configurtion is empty, the user did not pre-configure, so we need to set it up.
	if len(m.AWSConfig.PublicSubnetIPRange) == 0 {
		// map of default values.
		zoneDefaults := map[string]string{
			m.AWSConfig.Region + "a": "172.20.0.0/24",
			m.AWSConfig.Region + "b": "172.20.1.0/24",
			m.AWSConfig.Region + "c": "172.20.2.0/24",
			m.AWSConfig.Region + "d": "172.20.3.0/24",
			m.AWSConfig.Region + "e": "172.20.4.0/24",
		}

		// Get the available zones.
		availableZones, err := ec2S.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
		if err != nil {
			return err
		}

		// map available zones to default values.
		var tempZones []map[string]string
		for _, zone := range availableZones.AvailabilityZones {
			if *zone.State == "available" {
				subnetObj := map[string]string{
					"zone":      *zone.ZoneName,
					"ip_range":  zoneDefaults[*zone.ZoneName],
					"subnet_id": "",
				}
				tempZones = append(tempZones, subnetObj)
			}
		}

		// Are we multiAZ? if so add all available zones. If not default to first available.
		if m.AWSConfig.MultiAZ {
			m.AWSConfig.PublicSubnetIPRange = tempZones
		} else {
			m.AWSConfig.PublicSubnetIPRange = append(m.AWSConfig.PublicSubnetIPRange, tempZones[0])
		}
	}

	err := p.Core.DB.Save(m)
	if err != nil {
		return err
	}

	if m.AWSConfig.MasterRoleName == "" {
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
        "*"
      ]
    }
  ]
}`
			return createIAMRolePolicy(iamS, "kubernetes-master", policy)
		})

		procedure.AddStep("preparing IAM Instance Profile kubernetes-master", func() error {
			return createIAMInstanceProfile(iamS, "kubernetes-master")
		})
	}

	if m.AWSConfig.NodeRoleName == "" {
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
        "*"
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
	}
	if m.AWSConfig.BuildElasticFileSystem {

		procedure.AddStep("creating EFS share", func() error {

			input := &efs.CreateFileSystemInput{
				CreationToken: aws.String("cheese"),
			}
			resp, err := efsS.CreateFileSystem(input)
			if err != nil {
				return err
			}

			m.ServiceString = fmt.Sprintf(`
    - name: rpc-statd.service
      command: start
      enable: true
    - name: efs.service
      command: start
      content: |
        [Unit]
        After=network-online.target
        [Service]
        Type=oneshot
        ExecStartPre=-/usr/bin/mkdir -p /efs
        ExecStart=/bin/sh -c 'grep -qs /efs /proc/mounts || /usr/bin/mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $(/usr/bin/curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone).%s.efs.%s.amazonaws.com:/ /efs'
        ExecStop=/usr/bin/umount /efs
        RemainAfterExit=yes
        [Install]
        WantedBy=kubelet.service`, *resp.FileSystemId, m.AWSConfig.Region)

			m.AWSConfig.ElasticFileSystemID = *resp.FileSystemId

			err = p.Core.DB.Save(m)
			if err != nil {
				return err
			}

			return nil
		})
	}

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
		m.AWSConfig.BucketName = strings.ToLower("kubernetes-" + m.Name + "-" + util.RandomString(10))
		_, err := s3S.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(m.AWSConfig.BucketName),
		})
		if err != nil {
			return err
		}

		s3S.PutBucketPolicy(&s3.PutBucketPolicyInput{
			Bucket: aws.String(m.AWSConfig.BucketName), // Required
			Policy: aws.String(``),                     // Required
		})

		err = p.Core.DB.Save(m)
		if err != nil {
			return err
		}

		return nil
	})

	procedure.AddStep("upload assets to S3", func() error {
		// fetch the etcd clustering token
		if m.KubeMasterCount == 0 {
			m.KubeMasterCount = 1
		}
		url, err := etcdToken(strconv.Itoa(m.KubeMasterCount))
		if err != nil {
			return err
		}

		m.ETCDDiscoveryURL = url

		mversion := strings.Split(m.KubernetesVersion, ".")
		userdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/master.yaml")
		if err != nil {
			return err
		}
		template, err := template.New("minion_template").Parse(string(userdataTemplate))
		if err != nil {
			return err
		}

		var userdata bytes.Buffer
		if err = template.Execute(&userdata, m.AWSConfig.Tags); err != nil {
			return err
		}

		_, err = s3S.PutObject(&s3.PutObjectInput{
			Bucket:        aws.String(m.AWSConfig.BucketName), // Required
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

			m.AWSConfig.VPCMANAGED = true

			err := p.Core.DB.Save(m)
			if err != nil {
				return err
			}
			procedure.Core.Log.Info("This VPC is not managed. Using VPC ID supplied by user.")
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
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		return tagAWSResource(ec2S, m.AWSConfig.VPCID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-vpc",
		}, m.AWSConfig.Tags)
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
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
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
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		return tagAWSResource(ec2S, m.AWSConfig.InternetGatewayID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-ig",
		}, m.AWSConfig.Tags)
	})

	procedure.AddStep("attaching Internet Gateway to VPC", func() error {
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		input := &ec2.AttachInternetGatewayInput{
			VpcId:             aws.String(m.AWSConfig.VPCID),
			InternetGatewayId: aws.String(m.AWSConfig.InternetGatewayID),
		}
		if _, err := ec2S.AttachInternetGateway(input); err != nil && !strings.Contains(err.Error(), "already has an internet gateway attached") {
			return err
		}
		return nil
	})

	// Create Subnet

	procedure.AddStep("creating Subnet", func() error {
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		for idx, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if m.AWSConfig.PublicSubnetIPRange[idx]["subnet_id"] != "" {
				continue
			}

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
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				err := tagAWSResource(ec2S, subnet["subnet_id"], map[string]string{
					"KubernetesCluster": m.Name,
					"Name":              m.Name + "-" + subnet["zone"] + "-psub",
				}, m.AWSConfig.Tags)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	procedure.AddStep("enabling public IP assignment setting of Subnet", func() error {
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
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
		if m.AWSConfig.RouteTableID != "" || m.AWSConfig.VPCMANAGED == true {
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
		if m.AWSConfig.VPCMANAGED == true {
			return nil
		}
		return tagAWSResource(ec2S, m.AWSConfig.RouteTableID, map[string]string{
			"KubernetesCluster": m.Name,
			"Name":              m.Name + "-rt",
		}, m.AWSConfig.Tags)
	})

	procedure.AddStep("associating Route Table with Subnet", func() error {
		if len(m.AWSConfig.RouteTableSubnetAssociationID) != 0 || m.AWSConfig.VPCMANAGED == true {
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
			if strings.Contains(err.Error(), "already exists") {
				return nil
			}
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
		}, m.AWSConfig.Tags)
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
		}, m.AWSConfig.Tags)
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

	procedure.AddStep("setting EFS share target info", func() error {

		if m.AWSConfig.ElasticFileSystemID == "" {
			return nil
		}
		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			input := &efs.CreateMountTargetInput{
				FileSystemId:   aws.String(m.AWSConfig.ElasticFileSystemID),
				SecurityGroups: []*string{&m.AWSConfig.NodeSecurityGroupID},
				SubnetId:       aws.String(subnet["subnet_id"]),
			}
			target, err := efsS.CreateMountTarget(input)
			if err != nil {
				return err
			}

			m.AWSConfig.ElasticFileSystemTargets = append(m.AWSConfig.ElasticFileSystemTargets, *target.MountTargetId)
			err = p.Core.DB.Save(m)
			if err != nil {
				return err
			}

			return action.CancellableWaitFor("EFS target status", 20*time.Minute, time.Second, func() (bool, error) {
				resp, err := efsS.DescribeMountTargets(&efs.DescribeMountTargetsInput{
					MountTargetId: target.MountTargetId,
				})
				if err != nil {
					return false, nil
				}
				return *resp.MountTargets[0].LifeCycleState == "available", nil
			})

		}

		return nil
	})

	// Master Instance

	procedure.AddStep("creating Server for Kubernetes master(s)", func() error {
		if m.MasterID != "" {
			return nil
		}

		mversion := strings.Split(m.KubernetesVersion, ".")
		userdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/bootstrap.yaml")
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

		if m.KubeMasterCount == 0 {
			m.KubeMasterCount = 1
		}

		for i := 1; i <= m.KubeMasterCount; i++ {

			// Grab our subnets
			var subnets []string
			for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
				if subnet["subnet_id"] != "" {
					subnets = append(subnets, subnet["subnet_id"])
				}
			}
			// Randomly set a subnet
			var selectedSubnet string
			if len(subnets) == 1 {
				selectedSubnet = subnets[0]
			} else {
				selectedSubnet = subnets[(i-1)%len(m.AWSConfig.PublicSubnetIPRange)]
			}

			var pubNet bool
			if m.AWSConfig.PrivateNetwork {
				pubNet = false
			} else {
				pubNet = true
			}

			time.Sleep(5 * time.Second)
			procedure.Core.Log.Info("Building master #" + strconv.Itoa(i) + ", in subnet " + selectedSubnet + "...")

			var masterRole *string
			if m.AWSConfig.MasterRoleName != "" {
				masterRole = aws.String(m.AWSConfig.MasterRoleName)
			} else {
				masterRole = aws.String("kubernetes-master")
			}

			resp, err := ec2S.RunInstances(&ec2.RunInstancesInput{
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				ImageId:      aws.String(ami),
				InstanceType: aws.String(m.MasterNodeSize),
				KeyName:      aws.String(m.Name + "-key"),
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(0),
						AssociatePublicIpAddress: aws.Bool(pubNet),
						DeleteOnTermination:      aws.Bool(true),
						Groups: []*string{
							aws.String(m.AWSConfig.NodeSecurityGroupID),
						},
						SubnetId: aws.String(selectedSubnet),
						//PrivateIpAddress: aws.String(m.AWSConfig.MasterPrivateIP),
					},
				},
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: masterRole,
				},
				BlockDeviceMappings: []*ec2.BlockDeviceMapping{
					{
						DeviceName: aws.String("/dev/xvda"),
						Ebs: &ec2.EbsBlockDevice{
							DeleteOnTermination: aws.Bool(true),
							VolumeType:          aws.String("gp2"),
							VolumeSize:          aws.Int64(int64(m.AWSConfig.MasterVolumeSize)),
						},
					},
				},
				UserData: aws.String(encodedUserdata),
			})
			if err != nil {
				return err
			}

			m.MasterNodes = append(m.MasterNodes, *resp.Instances[0].InstanceId)

		}
		return nil
	})

	// Tag all of our masters.
	procedure.AddStep("tagging Kubernetes master", func() error {
		for _, master := range m.MasterNodes {
			err := tagAWSResource(ec2S, master, map[string]string{
				"KubernetesCluster": m.Name,
				"Name":              m.Name + "-master",
				"Role":              m.Name + "-master",
			}, m.AWSConfig.Tags)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// If we have more then one master... Lets get them on a internal loadbalancer.
	procedure.AddStep("creating master loadbalancer if needed", func() error {

		if m.KubeMasterCount == 1 {
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

		m.MasterPrivateIP = *loadbalancer.DNSName

		var elbInstances []*elb.Instance
		for _, node := range m.MasterNodes {
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
				aws.String(m.MasterNodes[0]),
			},
		}
		return action.CancellableWaitFor("Kubernetes master launch", 5*time.Minute, 3*time.Second, func() (bool, error) {
			resp, err := ec2S.DescribeInstances(input)
			if err != nil {
				return false, err
			}
			instance := resp.Reservations[0].Instances[0]

			//Always save private IP
			m.MasterPrivateIP = *instance.PrivateIpAddress
			if m.AWSConfig.PrivateNetwork {
				m.MasterPublicIP = *instance.PrivateIpAddress
			}
			p.Core.DB.Save(m)

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
