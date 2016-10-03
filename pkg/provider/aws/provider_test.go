package aws_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/test/fake_aws_provider"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSProviderValidateAccount(t *testing.T) {
	Convey("AWS Provider ValidateAccount works correctly", t, func() {
		table := []struct {
			// Input
			cloudAccount *model.CloudAccount
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				cloudAccount: &model.CloudAccount{},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						// NOTE won't matter until error is mocked
						DescribeKeyPairsFn: func(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
							output := &ec2.DescribeKeyPairsOutput{}
							return output, nil
						},
					}
				},
			}

			err := provider.ValidateAccount(item.cloudAccount)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

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
			mockDescribeInstancesError             error
			mockGetInstanceProfileError            error
			mockCreateInstanceProfileError         error
			mockAddRoleToInstanceProfileError      error
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
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
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
												PublicIpAddress: awssdk.String("0.0.0.0"),
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
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateKube(item.kube, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderDeleteKube(t *testing.T) {
	Convey("AWS Provider DeleteKube works correctly", t, func() {
		table := []struct {
			// Input
			kube *model.Kube
			// Mocks
			mockDeleteKeyPairError          error
			mockDeleteVpcError              error
			mockDeleteInternetGatewayError  error
			mockDeleteSubnetError           error
			mockDeleteRouteTableError       error
			mockDisassociateRouteTableError error
			mockDeleteSecurityGroupError    error
			mockTerminateInstancesError     error
			mockDescribeInstancesError      error
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
			},

			// A successful example (where there's nothing to delete)
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{},
				},
			},

			// On TerminateInstances error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockTerminateInstancesError: errors.New("TerminateInstances ERROR"),
				err: errors.New("TerminateInstances ERROR"),
			},

			// On DescribeInstances error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDescribeInstancesError: errors.New("DescribeInstances ERROR"),
				err: errors.New("DescribeInstances ERROR"),
			},

			// On DisassociateRouteTable error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDisassociateRouteTableError: errors.New("DisassociateRouteTable ERROR"),
				err: errors.New("DisassociateRouteTable ERROR"),
			},

			// On DescribeInstances error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDescribeInstancesError: errors.New("DescribeInstances ERROR"),
				err: errors.New("DescribeInstances ERROR"),
			},

			// On DeleteInternetGateway error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDeleteInternetGatewayError: errors.New("DeleteInternetGateway ERROR"),
				err: errors.New("DeleteInternetGateway ERROR"),
			},

			// On DeleteRouteTable error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDeleteRouteTableError: errors.New("DeleteRouteTable ERROR"),
				err: errors.New("DeleteRouteTable ERROR"),
			},

			// On DeleteSecurityGroup error
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						PublicSubnetID:                "PublicSubnetID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: "RouteTableSubnetAssociationID",
						ELBSecurityGroupID:            "ELBSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
						MasterID:                      "MasterID",
					},
				},
				mockDeleteSecurityGroupError: errors.New("DeleteSecurityGroup ERROR"),
				err: errors.New("DeleteSecurityGroup ERROR"),
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
						DeleteKeyPairFn: func(input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
							output := &ec2.DeleteKeyPairOutput{}
							return output, item.mockDeleteKeyPairError
						},
						DeleteVpcFn: func(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
							output := &ec2.DeleteVpcOutput{}
							return output, item.mockDeleteVpcError
						},
						DeleteInternetGatewayFn: func(input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
							output := &ec2.DeleteInternetGatewayOutput{}
							return output, item.mockDeleteInternetGatewayError
						},
						DeleteSubnetFn: func(input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
							output := &ec2.DeleteSubnetOutput{}
							return output, item.mockDeleteSubnetError
						},
						DeleteRouteTableFn: func(input *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
							output := &ec2.DeleteRouteTableOutput{}
							return output, item.mockDeleteRouteTableError
						},
						DisassociateRouteTableFn: func(input *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
							output := &ec2.DisassociateRouteTableOutput{}
							return output, item.mockDisassociateRouteTableError
						},
						DeleteSecurityGroupFn: func(input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
							output := &ec2.DeleteSecurityGroupOutput{}
							return output, item.mockDeleteSecurityGroupError
						},
						TerminateInstancesFn: func(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
							output := &ec2.TerminateInstancesOutput{}
							return output, item.mockTerminateInstancesError
						},
						DescribeInstancesFn: func(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
							output := &ec2.DescribeInstancesOutput{
								Reservations: []*ec2.Reservation{
									{
										Instances: []*ec2.Instance{
											{
												InstanceId: awssdk.String("instance-id"),
												State: &ec2.InstanceState{
													Name: awssdk.String("terminated"),
												},
											},
										},
									},
								},
							}
							return output, item.mockDescribeInstancesError
						},
					}
				},
			}

			err := provider.DeleteKube(item.kube)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderCreateNode(t *testing.T) {
	Convey("AWS Provider CreateNode works correctly", t, func() {
		table := []struct {
			// Input
			node *model.Node
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				node: &model.Node{
					Kube: &model.Kube{
						NodeSizes: []string{"m4.large"},
						AWSConfig: &model.AWSKubeConfig{},
						// Relations
						Entrypoints: []*model.Entrypoint{
							{
								Name: "my-entrypoint",
							},
						},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						RunInstancesFn: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
							output := &ec2.Reservation{
								Instances: []*ec2.Instance{
									{
										InstanceId:     awssdk.String("instance-id"),
										PrivateDnsName: awssdk.String("private.dns"),
										InstanceType:   awssdk.String("m4.large"),
										LaunchTime:     awssdk.Time(time.Now()),
									},
								},
							}
							return output, nil
						},
					}
				},
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						RegisterInstancesWithLoadBalancerFn: func(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
							output := &elb.RegisterInstancesWithLoadBalancerOutput{}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateNode(item.node, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderDeleteNode(t *testing.T) {
	Convey("AWS Provider DeleteNode works correctly", t, func() {
		table := []struct {
			// Input
			node *model.Node
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				node: &model.Node{
					Kube: &model.Kube{
						NodeSizes: []string{"m4.large"},
						AWSConfig: &model.AWSKubeConfig{},
						// Relations
						Entrypoints: []*model.Entrypoint{
							{
								Name: "my-entrypoint",
							},
						},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						TerminateInstancesFn: func(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
							output := &ec2.TerminateInstancesOutput{}
							return output, nil
						},
					}
				},
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						DeregisterInstancesFromLoadBalancerFn: func(input *elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {
							output := &elb.DeregisterInstancesFromLoadBalancerOutput{}
							return output, nil
						},
					}
				},
			}

			err := provider.DeleteNode(item.node)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderCreateVolume(t *testing.T) {
	Convey("AWS Provider CreateVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						CreateVolumeFn: func(input *ec2.CreateVolumeInput) (*ec2.Volume, error) {
							output := &ec2.Volume{
								VolumeId: awssdk.String("VolumeId"),
								Size:     awssdk.Int64(10),
							}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateVolume(item.volume, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderKubernetesVolumeDefinition(t *testing.T) {
	Convey("AWS Provider KubernetesVolumeDefinition works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			kubeVol *kubernetes.Volume
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					Name:       "test",
					ProviderID: "provider-ID",
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
					},
				},
				// Expectations
				kubeVol: &kubernetes.Volume{
					Name: "test",
					AwsElasticBlockStore: &kubernetes.AwsElasticBlockStore{
						VolumeID: "provider-ID",
						FSType:   "ext4",
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
			}

			kubeVol := provider.KubernetesVolumeDefinition(item.volume)

			So(kubeVol, ShouldResemble, item.kubeVol)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderResizeVolume(t *testing.T) {
	Convey("AWS Provider ResizeVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					ProviderID: "ProviderID",
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						CreateSnapshotFn: func(input *ec2.CreateSnapshotInput) (*ec2.Snapshot, error) {
							output := &ec2.Snapshot{
								SnapshotId: awssdk.String("SnapshotId"),
							}
							return output, nil
						},
						DescribeSnapshotsFn: func(input *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
							output := &ec2.DescribeSnapshotsOutput{
								Snapshots: []*ec2.Snapshot{
									{
										SnapshotId: awssdk.String("SnapshotId"),
										State:      awssdk.String("completed"),
									},
								},
							}
							return output, nil
						},
						CreateVolumeFn: func(input *ec2.CreateVolumeInput) (*ec2.Volume, error) {
							output := &ec2.Volume{
								VolumeId: awssdk.String("VolumeId"),
								Size:     awssdk.Int64(10),
							}
							return output, nil
						},
						DeleteVolumeFn: func(input *ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {
							output := &ec2.DeleteVolumeOutput{}
							return output, nil
						},
						DescribeVolumesFn: func(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
							output := &ec2.DescribeVolumesOutput{
								Volumes: []*ec2.Volume{
									{
										State: awssdk.String("available"),
									},
								},
							}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.ResizeVolume(item.volume, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderWaitForVolumeAvailable(t *testing.T) {
	Convey("AWS Provider WaitForVolumeAvailable works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						DescribeVolumesFn: func(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
							output := &ec2.DescribeVolumesOutput{
								Volumes: []*ec2.Volume{
									{
										State: awssdk.String("available"),
									},
								},
							}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.WaitForVolumeAvailable(item.volume, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderDeleteVolume(t *testing.T) {
	Convey("AWS Provider DeleteVolume works correctly", t, func() {
		table := []struct {
			// Input
			volume *model.Volume
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				volume: &model.Volume{
					ProviderID: "ProviderID",
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				EC2: func(kube *model.Kube) ec2iface.EC2API {
					return &fake_aws_provider.EC2{
						DescribeVolumesFn: func(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
							output := &ec2.DescribeVolumesOutput{
								Volumes: []*ec2.Volume{
									{
										State: awssdk.String("available"),
									},
								},
							}
							return output, nil
						},
						DeleteVolumeFn: func(input *ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {
							output := &ec2.DeleteVolumeOutput{}
							return output, nil
						},
					}
				},
			}

			err := provider.DeleteVolume(item.volume)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderCreateEntrypoint(t *testing.T) {
	Convey("AWS Provider CreateEntrypoint works correctly", t, func() {
		table := []struct {
			// Input
			entrypoint *model.Entrypoint
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				entrypoint: &model.Entrypoint{
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
						// Relations
						Nodes: []*model.Node{
							{
								Name: "a-node.host",
							},
						},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
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
						ConfigureHealthCheckFn: func(input *elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error) {
							output := &elb.ConfigureHealthCheckOutput{}
							return output, nil
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.CreateEntrypoint(item.entrypoint, action)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderDeleteEntrypoint(t *testing.T) {
	Convey("AWS Provider DeleteEntrypoint works correctly", t, func() {
		table := []struct {
			// Input
			entrypoint *model.Entrypoint
			// Mocks
			mockDeleteLoadBalancerError error
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				entrypoint: &model.Entrypoint{
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
						// Relations
						Nodes: []*model.Node{
							{
								Name: "a-node.host",
							},
						},
					},
				},
				err: nil,
			},

			// On DeleteLoadBalancer error
			{
				// Input
				entrypoint: &model.Entrypoint{
					Kube: &model.Kube{
						AWSConfig: &model.AWSKubeConfig{},
						// Relations
						Nodes: []*model.Node{
							{
								Name: "a-node.host",
							},
						},
					},
				},
				mockDeleteLoadBalancerError: errors.New("ERROR"),
				err: errors.New("ERROR"),
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						DeleteLoadBalancerFn: func(input *elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error) {
							output := &elb.DeleteLoadBalancerOutput{}
							return output, item.mockDeleteLoadBalancerError
						},
					}
				},
			}

			err := provider.DeleteEntrypoint(item.entrypoint)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderCreateEntrypointListener(t *testing.T) {
	Convey("AWS Provider CreateEntrypointListener works correctly", t, func() {
		table := []struct {
			// Input
			entrypointListener *model.EntrypointListener
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				entrypointListener: &model.EntrypointListener{
					Entrypoint: &model.Entrypoint{
						Kube: &model.Kube{
							AWSConfig: &model.AWSKubeConfig{},
						},
					},
				},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						// NOTE won't matter until error is mocked
						CreateLoadBalancerListenersFn: func(input *elb.CreateLoadBalancerListenersInput) (*elb.CreateLoadBalancerListenersOutput, error) {
							output := &elb.CreateLoadBalancerListenersOutput{}
							return output, nil
						},
					}
				},
			}

			err := provider.CreateEntrypointListener(item.entrypointListener)

			So(err, ShouldResemble, item.err)
		}
	})
}

//------------------------------------------------------------------------------

func TestAWSProviderDeleteEntrypointListener(t *testing.T) {
	Convey("AWS Provider DeleteEntrypointListener works correctly", t, func() {
		table := []struct {
			// Input
			entrypointListener *model.EntrypointListener
			// Mocks
			mockDeleteEntrypointListenerError error
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				entrypointListener: &model.EntrypointListener{
					Entrypoint: &model.Entrypoint{
						Kube: &model.Kube{
							AWSConfig: &model.AWSKubeConfig{},
						},
					},
				},
			},

			// On non-404 error
			{
				// Input
				entrypointListener: &model.EntrypointListener{
					Entrypoint: &model.Entrypoint{
						Kube: &model.Kube{
							AWSConfig: &model.AWSKubeConfig{},
						},
					},
				},
				mockDeleteEntrypointListenerError: errors.New("ERROR"),
				err: errors.New("ERROR"),
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &aws.Provider{
				Core: c,
				ELB: func(kube *model.Kube) elbiface.ELBAPI {
					return &fake_aws_provider.ELB{
						// NOTE won't matter until error is mocked
						DeleteLoadBalancerListenersFn: func(input *elb.DeleteLoadBalancerListenersInput) (*elb.DeleteLoadBalancerListenersOutput, error) {
							output := &elb.DeleteLoadBalancerListenersOutput{}
							return output, item.mockDeleteEntrypointListenerError
						},
					}
				},
			}

			err := provider.DeleteEntrypointListener(item.entrypointListener)

			So(err, ShouldResemble, item.err)
		}
	})
}
