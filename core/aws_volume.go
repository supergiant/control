package core

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
	vol := m.awsVolume()
	if vol == nil {
		panic(fmt.Errorf("Trying to access ID of nil volume %v#", vol))
	}
	return *vol.VolumeId
}

// simple memoization of aws vol record
func (m *AwsVolume) awsVolume() *ec2.Volume {
	if m.awsVol == nil {
		m.loadAwsVolume()
	}
	return m.awsVol
}

func (m *AwsVolume) loadAwsVolume() {
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
		panic(err) // TODO this isn't a 404, so we need to figure out what could happen; probably implement retry
	}

	if len(resp.Volumes) > 0 {
		m.awsVol = resp.Volumes[0]
	}
	// Volume does not exist otherwise and that's fine
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
	if m.awsVolume() == nil {
		return m.createAwsVolume()
	}
	return nil
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
	if m.awsVolume() == nil {
		return nil
	}
	if err := m.WaitForAvailable(); err != nil {
		return err
	}
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
