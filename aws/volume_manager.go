package aws

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	region = os.Getenv("REGION")
	az     = os.Getenv("AZ")
	svc    = ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
)

// Volume is a convenience wrapper around ec2.Volume to hold Name
type Volume struct {
	*ec2.Volume
	BaseName string // name without the instance number
	Name     string
}

type VolumeManager struct {
}

func (v *VolumeManager) Create(name string, ebsType string, size int) (*ec2.Volume, error) {
	volInput := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(az),
		VolumeType:       aws.String(ebsType),
		Size:             aws.Int64(int64(size)),
	}
	volume, err := svc.CreateVolume(volInput)
	if err != nil {
		return nil, err
	}
	tagsInput := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(*volume.VolumeId),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	}
	if _, err = svc.CreateTags(tagsInput); err != nil {
		return nil, err
	}
	return volume, nil
}

// Find is used instead of "Get" since it is a tag-based search
func (v *VolumeManager) Find(name string) (*ec2.Volume, error) {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	resp, err := svc.DescribeVolumes(input)
	if err != nil {
		return nil, err
	}
	if len(resp.Volumes) > 0 {
		return resp.Volumes[0], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Volume with name %s not found", name))
	}
}

func (v *VolumeManager) Delete(id string) error {
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	}
	_, err := svc.DeleteVolume(input)
	return err
}

func (v *VolumeManager) WaitForAvailable(id string) error {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(id),
				},
			},
		},
	}
	return svc.WaitUntilVolumeAvailable(input)
}
