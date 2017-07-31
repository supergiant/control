package aws

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// DeleteKube deletes a Kubernetes cluster.
func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
	ec2S := p.EC2(m)
	s3S := p.S3(m)
	efS := p.EFS(m)
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Delete Kube",
		Model:  m,
		Action: action,
	}

	procedure.AddStep("deleting master(s)", func() error {
		if len(m.MasterNodes) == 0 {
			return nil
		}

		for _, master := range m.MasterNodes {
			input := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{
					aws.String(master),
				},
			}
			if _, err := ec2S.TerminateInstances(input); isErrAndNotAWSNotFound(err) {
				return err
			}

			// Wait for termination
			descinput := &ec2.DescribeInstancesInput{
				InstanceIds: []*string{
					aws.String(master),
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

		}
		m.MasterID = ""
		return nil
	})

	procedure.AddStep("destroying api loadbalancer if it exists", func() error {
		// Delete ELB
		_, err := p.ELB(m).DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
			LoadBalancerName: aws.String(m.Name + "-api"),
		})
		if isErrAndNotAWSNotFound(err) {
			return err
		}
		return nil
	})

	procedure.AddStep("delete EFS targets", func() error {
		// Delete EFS targets
		if len(m.AWSConfig.ElasticFileSystemTargets) == 0 {
			return nil
		}
		for _, target := range m.AWSConfig.ElasticFileSystemTargets {
			_, err := efS.DeleteMountTarget(&efs.DeleteMountTargetInput{
				MountTargetId: aws.String(target),
			})
			if isErrAndNotAWSNotFound(err) {
				return err
			}
		}
		return nil
	})

	procedure.AddStep("disassociating Route Table from Subnet(s)", func() error {
		if len(m.AWSConfig.RouteTableSubnetAssociationID) == 0 || m.AWSConfig.VPCMANAGED == true {
			return nil
		}

		for _, assID := range m.AWSConfig.RouteTableSubnetAssociationID {
			_, err := ec2S.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
				AssociationId: aws.String(assID),
			})
			if isErrAndNotAWSNotFound(err) {
				return err
			}
		}
		m.AWSConfig.RouteTableSubnetAssociationID = []string{}
		return nil
	})

	procedure.AddStep("deleting Internet Gateway", func() error {
		if m.AWSConfig.InternetGatewayID == "" || m.AWSConfig.VPCMANAGED == true {
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
				if strings.Contains(err.Error(), "does not exist") {
					// it does not exist,
					return true, nil
				}
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
			if strings.Contains(err.Error(), "does not exist") {
				// it does not exist,
				return nil
			}
			return err
		}
		m.AWSConfig.InternetGatewayID = ""
		return nil
	})

	procedure.AddStep("deleting Route Table", func() error {
		if m.AWSConfig.VPCMANAGED == true || m.AWSConfig.RouteTableID == "" {
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

	procedure.AddStep("deleting public Subnet(s)", func() error {
		if len(m.AWSConfig.PublicSubnetIPRange) == 0 || m.AWSConfig.VPCMANAGED == true {
			return nil
		}

		for _, subnet := range m.AWSConfig.PublicSubnetIPRange {
			if subnet["subnet_id"] != "" {
				input := &ec2.DeleteSubnetInput{
					SubnetId: aws.String(subnet["subnet_id"]),
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
			}
		}

		m.AWSConfig.PublicSubnetIPRange = []map[string]string{}
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

	procedure.AddStep("deleting S3 bucket", func() error {

		// if bucket does not exist skip.
		if !bucketExist(s3S, m.AWSConfig.BucketName) {
			return nil
		}

		objects, err := s3S.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(m.AWSConfig.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "The authorization header is malformed") {
				// it does not exist,
				return nil
			}
			return err
		}

		for _, object := range objects.Contents {
			_, err = s3S.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(m.AWSConfig.BucketName),
				Key:    aws.String(*object.Key),
			})
			if err != nil {
				if strings.Contains(err.Error(), "The specified method is not allowed against this resource.") {
					// it does not exist,
					return nil
				}
				return err
			}
		}

		_, err = s3S.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(m.AWSConfig.BucketName),
		})
		if err != nil {
			if strings.Contains(err.Error(), "The specified method is not allowed against this resource.") {
				// it does not exist,
				return nil
			}
			return err
		}
		return nil
	})

	procedure.AddStep("deleting VPC", func() error {
		if m.AWSConfig.VPCMANAGED {
			procedure.Core.Log.Info("This VPC is not managed. It is NOT being deleted.")
			return nil
		}
		// Delete the VPC
		_, err := ec2S.DeleteVpc(&ec2.DeleteVpcInput{
			VpcId: aws.String(m.AWSConfig.VPCID),
		})
		if isErrAndNotAWSNotFound(err) {
			if strings.Contains(err.Error(), "The request must contain the parameter") {
				// it does not exist,
				return nil
			}
			return err
		}
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
