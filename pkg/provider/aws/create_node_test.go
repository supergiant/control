package aws_test

import (
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/test/fake_aws_provider"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

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
						AWSConfig: &model.AWSKubeConfig{
							PublicSubnetIPRange: []map[string]string{
								{
									"subnet_id": "test",
								},
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
						DescribeImagesFn: func(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
							output := &ec2.DescribeImagesOutput{
								Images: []*ec2.Image{
									{
										ImageId:      awssdk.String("ami-1234"),
										CreationDate: awssdk.String("August 24, 2016 at 4:36:22 PM UTC-5"),
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
