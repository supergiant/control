package model

import (
	"fmt"
	"guber"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsVolume struct {
	c         *Client
	Blueprint *VolumeBlueprint
	Instance  *Instance

	awsVol *ec2.Volume // used internally to store record of AWS vol
}

func (m *AwsVolume) name() string {
	return fmt.Sprintf("%s-%s", m.Instance.Name(), m.Blueprint.Name)
}

func (m *AwsVolume) id() string {
	return *m.awsVolume().VolumeId
}

// simple memoization of aws vol record
func (m *AwsVolume) awsVolume() *ec2.Volume {
	if m.awsVol == nil {
		if err := m.loadAwsVolume(); err != nil {
			panic(err) // TODO
		}
	}
	return m.awsVol
}

func (m *AwsVolume) loadAwsVolume() error {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(m.name()),
				},
			},
		},
	}
	resp, err := m.c.EC2.DescribeVolumes(input)
	if err != nil {
		return err // this could presumably be a rate-limit based error...
	}

	if len(resp.Volumes) == 0 {
		return fmt.Errorf("Volume with name %s not found", m.name())
	}
	m.awsVol = resp.Volumes[0]
	return nil
}

func (m *AwsVolume) createAwsVolume() error {
	volInput := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(AwsAZ),
		VolumeType:       aws.String(m.Blueprint.Type),
		Size:             aws.Int64(int64(m.Blueprint.Size)),
	}

	awsVol, err := m.c.EC2.CreateVolume(volInput)
	if err != nil {
		return err
	}
	tagsInput := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(*awsVol.VolumeId),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(m.name()),
			},
		},
	}
	if _, err = m.c.EC2.CreateTags(tagsInput); err != nil {
		return err // TODO an error here means we create a hanging volume, since it does not get named
	}
	m.awsVol = awsVol
	return nil
}

func (m *AwsVolume) Provision() error {
	if err := m.loadAwsVolume(); err != nil {
		return err
	} else if m.awsVolume != nil {
		return nil
	}
	return m.createAwsVolume()
}

func (m *AwsVolume) WaitForAvailable() error {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(m.id()),
				},
			},
		},
	}
	return m.c.EC2.WaitUntilVolumeAvailable(input)
}

func (m *AwsVolume) Destroy() error {
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(m.id()),
	}
	_, err := m.c.EC2.DeleteVolume(input)
	return err
}

func (m *AwsVolume) AsKubeVolume() *guber.Volume {
	return &guber.Volume{
		Name: m.Blueprint.Name, // NOTE this is not the physical volume name
		AwsElasticBlockStore: &guber.AwsElasticBlockStore{
			VolumeID: m.id(),
			FSType:   "ext4",
		},
	}
}
