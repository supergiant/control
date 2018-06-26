package aws_test

import (
	"errors"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/efs/efsiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/test/fake_aws_provider"
	"github.com/supergiant/supergiant/test/fake_core"
)

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
			mockCreateBucketError           error
			mockPutBucketPolicyError        error
			mockListObjectsError            error
			mockDeleteObjectError           error
			mockPutObjectError              error
			mockDeleteBucketError           error
			mockDeleteLoadBalancerError     error
			mockDeleteMountTargetError      error
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				kube: &model.Kube{
					NodeSizes: []string{"m4.large"},
					MasterID:  "MasterID",
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
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
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockTerminateInstancesError: errors.New("TerminateInstances ERROR"),
				err: errors.New("TerminateInstances ERROR"),
			},

			// On DescribeInstances error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockDescribeInstancesError: errors.New("DescribeInstances ERROR"),
				err: errors.New("DescribeInstances ERROR"),
			},

			// On DisassociateRouteTable error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockDisassociateRouteTableError: errors.New("DisassociateRouteTable ERROR"),
				err: errors.New("DisassociateRouteTable ERROR"),
			},

			// On DescribeInstances error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockDescribeInstancesError: errors.New("DescribeInstances ERROR"),
				err: errors.New("DescribeInstances ERROR"),
			},

			// On DeleteInternetGateway error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockDeleteInternetGatewayError: errors.New("DeleteInternetGateway ERROR"),
				err: errors.New("DeleteInternetGateway ERROR"),
			},

			// On DeleteRouteTable error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
					},
				},
				mockDeleteRouteTableError: errors.New("DeleteRouteTable ERROR"),
				err: errors.New("DeleteRouteTable ERROR"),
			},

			// On DeleteSecurityGroup error
			{
				// Input
				kube: &model.Kube{
					NodeSizes:   []string{"m4.large"},
					MasterID:    "MasterID",
					MasterNodes: []string{"MasterNode"},
					AWSConfig: &model.AWSKubeConfig{
						VPCID:                         "VPCID",
						InternetGatewayID:             "InternetGatewayID",
						RouteTableID:                  "RouteTableID",
						RouteTableSubnetAssociationID: []string{"RouteTableSubnetAssociationID"},
						MasterSecurityGroupID:         "MasterSecurityGroupID",
						NodeSecurityGroupID:           "NodeSecurityGroupID",
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
						DeleteLoadBalancerFn: func(input *elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error) {
							output := &elb.DeleteLoadBalancerOutput{}
							return output, item.mockDeleteLoadBalancerError
						},
					}
				},
				EFS: func(kube *model.Kube) efsiface.EFSAPI {
					return &fake_aws_provider.EFS{
						DeleteMountTargetfn: func(input *efs.DeleteMountTargetInput) (*efs.DeleteMountTargetOutput, error) {
							output := &efs.DeleteMountTargetOutput{}
							return output, item.mockDeleteMountTargetError
						},
					}
				},
			}

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.DeleteKube(item.kube, action)

			So(err, ShouldResemble, item.err)
		}
	})
}
