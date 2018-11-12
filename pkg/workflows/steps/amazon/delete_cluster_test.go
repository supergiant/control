package amazon

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type FakeEC2DeleteCluster struct {
	ec2iface.EC2API

	instancesOutput  *ec2.DescribeInstancesOutput
	terminatedOutput *ec2.TerminateInstancesOutput
	err              error
}

func (f *FakeEC2DeleteCluster) DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return f.instancesOutput, f.err
}

func (f *FakeEC2DeleteCluster) TerminateInstancesWithContext(aws.Context, *ec2.TerminateInstancesInput, ...request.Option) (*ec2.TerminateInstancesOutput, error) {
	return f.terminatedOutput, f.err
}

func TestDeleteClusterStep_Run(t *testing.T) {
	//TBD
}
