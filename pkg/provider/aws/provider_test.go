package aws_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/test/fake_aws_provider"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

type mockS3Client struct {
	s3iface.S3API
}

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

			action := &core.Action{Status: new(model.ActionStatus)}
			err := provider.DeleteNode(item.node, action)

			So(err, ShouldResemble, item.err)
		}
	})
}
