package amazon

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type FakeEC2KDeleteCluster struct {
	ec2iface.EC2API

	instancesOutput  *ec2.DescribeInstancesOutput
	terminatedOutput *ec2.TerminateInstancesOutput

	err error
}

func TestDeleteClusterStep_Run(t *testing.T) {

}
