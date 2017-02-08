package aws_test

import (
	"errors"
	"testing"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/test/fake_aws_provider"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSProviderCreateKube(t *testing.T) {
	Convey("AWS Provider CreateKube works correctly", t, func() {
		table := []struct {
			// Input
			kube *model.Kube
			// Mocks
			mockCreateKeyPairError                 error
			mockCreateVpcError                     error
			mockCreateInternetGatewayError         error
			mockAttachInternetGatewayError         error
			mockCreateSubnetError                  error
			mockCreateRouteTableError              error
			mockCreateRouteError                   error
			mockAssociateRouteTableError           error
			mockCreateSecurityGroupError           error
			mockAuthorizeSecurityGroupIngressError error
			mockRunInstancesError                  error
			mockDescribeAvailabilityZonesError     error
			mockDescribeInstancesError             error
			mockGetInstanceProfileError            error
			mockCreateInstanceProfileError         error
			mockAddRoleToInstanceProfileError      error
			mockCreateBucketError                  error
			mockPutBucketPolicyError               error
			mockListObjectsError                   error
			mockDeleteObjectError                  error
			mockPutObjectError                     error
			mockDeleteBucketError                  error
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockGetInstanceProfileError: errors.New("404"), // doesn't exist on first run
				err: nil,
			},

			// When there's an error getting existing InstanceProfile
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockGetInstanceProfileError: errors.New("500: bad thing"),
				err: errors.New("500: bad thing"),
			},

			// When there's an error creating InstanceProfile
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockGetInstanceProfileError:    errors.New("404"),
				mockCreateInstanceProfileError: errors.New("uh oh"),
				err: errors.New("uh oh"),
			},

			// When there's an error on AddRoleToInstanceProfile
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockGetInstanceProfileError:       errors.New("404"),
				mockAddRoleToInstanceProfileError: errors.New("bad"),
				err: errors.New("bad"),
			},

			// A successful example (when everything already exists)
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						PrivateKey:                    "PrivateKey",
						SSHPubKey:                     "SSHPubKey",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterNodes:                   []string{"MasterNode"},
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				err: nil,
			},

			// A successful example of a user selecting a multi-az deployment, with default values for subnets.
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						PrivateKey:        "PrivateKey",
						SSHPubKey:         "SSHPubKey",
						InternetGatewayID: "InternetGatewayID",
						RouteTableID:      "RouteTableID",
						MultiAZ:           true,
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterNodes:                   []string{"MasterNode"},
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				err: nil,
			},
			// When there's a duplicate private key
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateKeyPairError: errors.New("Duplicate"),
				err: errors.New("KeyPair existed, but key material was not captured. Deleted KeyPair... will retry"),
			},

			// When there's an unexpected error with CreateKeyPair
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateKeyPairError: errors.New("CreateKeyPair ERROR"),
				err: errors.New("CreateKeyPair ERROR"),
			},

			// When there's an error on CreateVpc
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateVpcError: errors.New("CreateVpc ERROR"),
				err:                errors.New("CreateVpc ERROR"),
			},

			// When there's an error on CreateInternetGateway
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateInternetGatewayError: errors.New("CreateInternetGateway ERROR"),
				err: errors.New("CreateInternetGateway ERROR"),
			},

			// When there's an error on AttachInternetGateway
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockAttachInternetGatewayError: errors.New("AttachInternetGateway ERROR"),
				err: errors.New("AttachInternetGateway ERROR"),
			},

			// When there's an error on CreateSubnet
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateSubnetError: errors.New("CreateSubnet ERROR"),
				err: errors.New("CreateSubnet ERROR"),
			},

			// When there's an error on CreateRouteTable
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateRouteTableError: errors.New("CreateRouteTable ERROR"),
				err: errors.New("CreateRouteTable ERROR"),
			},

			// When there's an error on CreateRoute
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateRouteError: errors.New("CreateRoute ERROR"),
				err:                  errors.New("CreateRoute ERROR"),
			},

			// When there's an error on AssociateRouteTable
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockAssociateRouteTableError: errors.New("AssociateRouteTable ERROR"),
				err: errors.New("AssociateRouteTable ERROR"),
			},

			// When there's an error on CreateSecurityGroup
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockCreateSecurityGroupError: errors.New("CreateSecurityGroup ERROR"),
				err: errors.New("CreateSecurityGroup ERROR"),
			},

			// When there's an error on AuthorizeSecurityGroupIngress
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockAuthorizeSecurityGroupIngressError: errors.New("AuthorizeSecurityGroupIngress ERROR"),
				err: errors.New("AuthorizeSecurityGroupIngress ERROR"),
			},

			// When there's an error on RunInstances
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockRunInstancesError: errors.New("RunInstances ERROR"),
				err: errors.New("RunInstances ERROR"),
			},

			// When there's an error on DescribeInstances
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
				mockDescribeInstancesError: errors.New("DescribeInstances ERROR"),
				err: errors.New("DescribeInstances ERROR"),
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),

				K8S: func(*model.Kube) kubernetes.ClientInterface {
					return &fake_core.KubernetesClient{
						ListNodesFn: func(query string) ([]*kubernetes.Node, error) {
							return []*kubernetes.Node{
								{
									Metadata: kubernetes.Metadata{
										Name: "created-node",
									},
								},
							}, nil
						},
					}
				},

				Nodes: new(fake_core.Nodes),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						CreateKeyPairFn: func(input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
							output := &ec2.CreateKeyPairOutput{
								KeyMaterial: awssdk.String("bleep-blorp"),
							}
							return output, item.mockCreateKeyPairError
						},
						CreateVpcFn: func(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
							output := &ec2.CreateVpcOutput{
								Vpc: &ec2.Vpc{
									VpcId: awssdk.String("bloop"),
								},
							}
							return output, item.mockCreateVpcError
						},
						CreateInternetGatewayFn: func(input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
							output := &ec2.CreateInternetGatewayOutput{
								InternetGateway: &ec2.InternetGateway{
									InternetGatewayId: awssdk.String("bloop"),
								},
							}
							return output, item.mockCreateInternetGatewayError
						},
						AttachInternetGatewayFn: func(input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
							return nil, item.mockAttachInternetGatewayError
						},
						CreateSubnetFn: func(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
							output := &ec2.CreateSubnetOutput{
								Subnet: &ec2.Subnet{
									SubnetId: awssdk.String("bloop"),
								},
							}
							return output, item.mockCreateSubnetError
						},
						CreateRouteTableFn: func(input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
							output := &ec2.CreateRouteTableOutput{
								RouteTable: &ec2.RouteTable{
									RouteTableId: awssdk.String("bloop"),
								},
							}
							return output, item.mockCreateRouteTableError
						},
						CreateRouteFn: func(input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
							return nil, item.mockCreateRouteError
						},
						AssociateRouteTableFn: func(input *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
							output := &ec2.AssociateRouteTableOutput{
								AssociationId: awssdk.String("spicy-boy"),
							}
							return output, item.mockAssociateRouteTableError
						},
						CreateSecurityGroupFn: func(input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
							output := &ec2.CreateSecurityGroupOutput{
								GroupId: awssdk.String("bloop"),
							}
							return output, item.mockCreateSecurityGroupError
						},
						DescribeImagesFn: func(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
							output := &ec2.DescribeImagesOutput{
								Images: []*ec2.Image{
									&ec2.Image{
										ImageId:      awssdk.String("ami-1234"),
										CreationDate: awssdk.String("August 24, 2016 at 4:36:22 PM UTC-5"),
									},
								},
							}
							return output, item.mockCreateSecurityGroupError
						},
						AuthorizeSecurityGroupIngressFn: func(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
							return nil, item.mockAuthorizeSecurityGroupIngressError
						},
						RunInstancesFn: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
							output := &ec2.Reservation{
								Instances: []*ec2.Instance{
									{
										InstanceId: awssdk.String("instance-id"),
									},
								},
							}
							return output, item.mockRunInstancesError
						},
						DescribeAvailabilityZonesFn: func(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
							output := &ec2.DescribeAvailabilityZonesOutput{
								AvailabilityZones: []*ec2.AvailabilityZone{
									{
										ZoneName: awssdk.String("aws-zone"),
										State:    awssdk.String("available"),
									},
								},
							}
							return output, item.mockRunInstancesError
						},
						DescribeInstancesFn: func(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
							output := &ec2.DescribeInstancesOutput{
								Reservations: []*ec2.Reservation{
									{
										Instances: []*ec2.Instance{
											{
												InstanceId: awssdk.String("instance-id"),
												State: &ec2.InstanceState{
													Name: awssdk.String("running"),
												},
												PublicIpAddress:  awssdk.String("0.0.0.0"),
												PrivateIpAddress: awssdk.String("0.0.0.0"),
											},
										},
									},
								},
							}
							return output, item.mockDescribeInstancesError
						},
					}
				},
				IAM: func(kube *model.Kube) iamiface.IAMAPI {
					return &fake_aws_provider.IAM{
						GetInstanceProfileFn: func(input *iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
							output := &iam.GetInstanceProfileOutput{
								InstanceProfile: &iam.InstanceProfile{
									Roles: []*iam.Role{},
								},
							}
							return output, item.mockGetInstanceProfileError
						},
						CreateInstanceProfileFn: func(input *iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
							output := &iam.CreateInstanceProfileOutput{
								InstanceProfile: &iam.InstanceProfile{
									Roles: []*iam.Role{},
								},
							}
							return output, item.mockCreateInstanceProfileError
						},
						AddRoleToInstanceProfileFn: func(input *iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
							output := &iam.AddRoleToInstanceProfileOutput{}
							return output, item.mockAddRoleToInstanceProfileError
						},
					}
				},
				S3: func(kube *model.Kube) s3iface.S3API {
					return &fake_aws_provider.S3{
						CreateBucketFn: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
							output := &s3.CreateBucketOutput{
								Location: awssdk.String("test"),
							}
							return output, item.mockCreateBucketError
						},
						PutBucketPolicyFn: func(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
							output := &s3.PutBucketPolicyOutput{}
							return output, item.mockPutBucketPolicyError
						},
						PutObjectFn: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
							output := &s3.PutObjectOutput{}
							return output, item.mockPutObjectError
						},
						ListObjectsFn: func(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
							output := &s3.ListObjectsOutput{}
							return output, item.mockListObjectsError
						},
						DeleteObjectFn: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
							output := &s3.DeleteObjectOutput{}
							return output, item.mockDeleteObjectError
						},
						DeleteBucketFn: func(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
							output := &s3.DeleteBucketOutput{}
							return output, item.mockDeleteBucketError
						},
					}
				},
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						CreateLoadBalancerFn: func(input *elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error) {
							output := &elb.CreateLoadBalancerOutput{
								DNSName: awssdk.String("elb.dns.blah.blah"),
							}
							return output, nil
						},
						// NOTE for both of the following, it doesn't really matter due to
						// how we mock (we don't care about the return, just the error
						// here). It will only matter once we mock an error.
						RegisterInstancesWithLoadBalancerFn: func(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
							output := &elb.RegisterInstancesWithLoadBalancerOutput{}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateKube(item.kube, action)

			So(err, ShouldResemble, item.err)
		}
	})
}
