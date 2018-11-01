package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	errFake = errors.New("aws fake error")

	awsTagName = "aws:name"
	awsTagKey  = "aws:key"
	awsTagVal  = "aws:val"

	ec2t2MicroID   = "m4.micro.id"
	ec2t2MicroName = "m4.micro.name"
	ec2t2MicroInst = &ec2.Instance{
		InstanceId: aws.String(ec2t2MicroID),
		Tags: []*ec2.Tag{
			{Key: aws.String(awsTagName), Value: aws.String(ec2t2MicroName)},
			{Key: aws.String(awsTagKey), Value: aws.String(awsTagVal)},
		},
	}

	ec2m4LargeID   = "m4.large.id"
	ec2m4LargeName = "m4.large.name"
	ec2m4LargeInst = &ec2.Instance{
		InstanceId: aws.String(ec2m4LargeID),
		Tags: []*ec2.Tag{
			{Key: aws.String(awsTagName), Value: aws.String(ec2m4LargeName)},
			{Key: aws.String(awsTagKey), Value: aws.String(awsTagVal)},
		},
	}

	ec2Reservation = &ec2.Reservation{
		Instances: []*ec2.Instance{ec2t2MicroInst, ec2m4LargeInst},
	}
)

// https://aws.amazon.com/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/
type fakeEC2Service struct {
	ec2iface.EC2API
	ec2Reservation *ec2.Reservation
	err            error
}

func (m *fakeEC2Service) RunInstancesWithContext(ctx aws.Context, input *ec2.RunInstancesInput, opts ...request.Option) (*ec2.Reservation, error) {
	return m.ec2Reservation, m.err
}
func (m *fakeEC2Service) CreateTags(tags *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return nil, m.err
}
func (m *fakeEC2Service) TerminateInstancesWithContext(aws.Context, *ec2.TerminateInstancesInput, ...request.Option) (*ec2.TerminateInstancesOutput, error) {
	return nil, m.err
}
func (m *fakeEC2Service) DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	if m.ec2Reservation == nil {
		return &ec2.DescribeInstancesOutput{}, m.err
	}
	return &ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{m.ec2Reservation}}, m.err
}

func TestNewClient(t *testing.T) {
	tcs := []struct {
		id, secret  string
		expectedErr error
	}{
		{"", "", ErrInvalidKeys},
		{"   ", "adfoenkfad", ErrInvalidKeys},
		{"123", "123", nil},
	}

	for i, tc := range tcs {
		_, err := New(tc.id, tc.secret, nil)
		if err != tc.expectedErr {
			t.Fatalf("TC#%d: new client: %v", i+1, err)
		}
	}
}

func TestClient_CreateInstance(t *testing.T) {
	tcs := []struct {
		name   string
		config InstanceConfig

		ec2Res *ec2.Reservation
		ec2Err error

		expectedErr error
	}{
		// TC#1
		{
			name:        "no region provided",
			expectedErr: ErrNoRegionProvided,
		},
		// TC#2
		{
			name:        "failed to run instance",
			config:      InstanceConfig{Region: "us1"},
			ec2Err:      errFake,
			expectedErr: errFake,
		},
		// TC#3
		{
			name:        "no instances created",
			config:      InstanceConfig{Region: "us1"},
			ec2Res:      &ec2.Reservation{},
			expectedErr: ErrNoInstancesCreated,
		},
		// TC#5
		{
			config: InstanceConfig{
				Region:        "us1",
				UsedData:      "userdata",
				HasPublicAddr: true,
			},
			ec2Res: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{InstanceId: aws.String("1")},
				},
			},
			name: "create and tag instance",
		},
	}

	ec2Fake := &fakeEC2Service{}
	for i, tc := range tcs {
		ec2Fake.ec2Reservation, ec2Fake.err = tc.ec2Res, tc.ec2Err
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Fake
			},
		}

		_, err := c.CreateInstance(context.Background(), tc.config)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: %s", i+1, tc.name)
	}
}

func TestClient_DeleteInstance(t *testing.T) {
	tcs := []struct {
		name       string
		instanceID string
		region     string

		ec2DeleteInsErr error

		expectedErr error
	}{
		// TC#1
		{
			name:        "no region provided",
			expectedErr: ErrNoRegionProvided,
		},
		// TC#2
		{
			name:        "empty instance id",
			region:      "us1",
			expectedErr: ErrInstanceIDEmpty,
		},
		// TC#3
		{
			name:            "failed to delete instance",
			region:          "us1",
			instanceID:      "123",
			ec2DeleteInsErr: errFake,
			expectedErr:     errFake,
		},
		// TC#4
		{
			name:       "failed to delete instance",
			region:     "us1",
			instanceID: "123",
		},
	}

	ec2Fake := &fakeEC2Service{}
	for i, tc := range tcs {
		ec2Fake.err = tc.ec2DeleteInsErr
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Fake
			},
		}

		_, err := c.DeleteInstance(context.Background(), tc.region, tc.instanceID)
		if err != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: %s", i+1, tc.name)
		}
	}
}

func TestClient_GetInstance(t *testing.T) {
	tcs := []struct {
		region      string
		instanceID  string
		expectedRes *ec2.Instance
		ec2Err      error
		expectedErr error
	}{
		{ // TC#1
			expectedErr: ErrNoRegionProvided,
		},
		{ // TC#2
			region:      "us1",
			expectedErr: ErrInstanceIDEmpty,
		},
		{ // TC#3
			region:      "us1",
			instanceID:  "123",
			ec2Err:      errFake,
			expectedErr: errFake,
		},
		{ // TC#4
			region:      "us1",
			instanceID:  "123",
			expectedErr: ErrInstanceNotFound,
		},
		{ // TC#5
			region:      "us1",
			instanceID:  ec2t2MicroID,
			expectedRes: ec2t2MicroInst,
		},
	}

	ec2Fake := &fakeEC2Service{
		ec2Reservation: ec2Reservation,
	}
	for i, tc := range tcs {
		ec2Fake.err = tc.ec2Err
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Fake
			},
		}

		res, err := c.GetInstance(context.Background(), tc.region, tc.instanceID)
		if err != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d", i+1)
		}
		if res != nil {
			require.Equalf(t, tc.expectedRes, res, "TC#%d", i+1)
		}
	}
}

func TestClient_ListRegionInstances(t *testing.T) {
	tcs := []struct {
		region      string
		reservation *ec2.Reservation
		expectedRes []*ec2.Instance
		ec2Err      error
		expectedErr error
	}{
		{ // TC#1
			expectedErr: ErrNoRegionProvided,
		},
		{ // TC#2
			region:      "us1",
			ec2Err:      errFake,
			expectedErr: errFake,
		},
		{ // TC#3
			region:      "us1",
			expectedRes: []*ec2.Instance{},
		},
		{ // TC#4
			region:      "us1",
			reservation: ec2Reservation,
			expectedRes: []*ec2.Instance{ec2t2MicroInst, ec2m4LargeInst},
		},
	}

	ec2Fake := &fakeEC2Service{}
	for i, tc := range tcs {
		ec2Fake.ec2Reservation, ec2Fake.err = tc.reservation, tc.ec2Err
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Fake
			},
		}

		res, err := c.ListRegionInstances(context.Background(), tc.region, nil)
		if err != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d", i+1)
		}
		if res != nil {
			require.Equalf(t, tc.expectedRes, res, "TC#%d", i+1)
		}
	}
}

func TestEC2Svc(t *testing.T) {
	svc := ec2Svc(session.New(&aws.Config{}), "us1")
	_, ok := svc.(ec2iface.EC2API)
	if !ok {
		t.Fatal("ec2Svc func should return a EC2API interface object")
	}
}
