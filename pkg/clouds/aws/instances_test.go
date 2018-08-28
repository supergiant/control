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

// https://aws.amazon.com/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/
type mockedEC2Service struct {
	ec2iface.EC2API
	CreateInstRes *ec2.Reservation
	RunInstErr    error
	DelInstErr    error
	TagResErr     error
}

func (m *mockedEC2Service) RunInstancesWithContext(ctx aws.Context, input *ec2.RunInstancesInput, opts ...request.Option) (*ec2.Reservation, error) {
	return m.CreateInstRes, m.RunInstErr
}
func (m *mockedEC2Service) CreateTags(tags *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return nil, m.TagResErr
}
func (m *mockedEC2Service) TerminateInstancesWithContext(aws.Context, *ec2.TerminateInstancesInput, ...request.Option) (*ec2.TerminateInstancesOutput, error) {
	return nil, m.DelInstErr
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
	fakeRunInstErr := errors.New("run isntance error")

	tcs := []struct {
		name   string
		config InstanceConfig

		ec2Res        *ec2.Reservation
		ec2RunInstErr error
		ec2TagResErr  error

		expectedErr error
	}{
		// TC#1
		{
			name:        "no region provided",
			expectedErr: ErrNoRegionProvided,
		},
		// TC#2
		{
			name:          "failed to run instance",
			config:        InstanceConfig{Region: "us1"},
			ec2RunInstErr: fakeRunInstErr,
			expectedErr:   fakeRunInstErr,
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

	ec2Mock := &mockedEC2Service{}
	for i, tc := range tcs {
		ec2Mock.CreateInstRes, ec2Mock.RunInstErr, ec2Mock.TagResErr = tc.ec2Res, tc.ec2RunInstErr, tc.ec2TagResErr
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Mock
			},
		}

		_, err := c.CreateInstance(context.Background(), tc.config)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: %s", i+1, tc.name)
	}
}

func TestClient_DeleteInstance(t *testing.T) {
	fakeDelInstErr := errors.New("delete isntance error")

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
			ec2DeleteInsErr: fakeDelInstErr,
			expectedErr:     fakeDelInstErr,
		},
		// TC#4
		{
			name:       "failed to delete instance",
			region:     "us1",
			instanceID: "123",
		},
	}

	ec2Mock := &mockedEC2Service{}
	for i, tc := range tcs {
		ec2Mock.DelInstErr = tc.ec2DeleteInsErr
		c := &Client{
			ec2SvcFn: func(s *session.Session, region string) ec2iface.EC2API {
				return ec2Mock
			},
		}

		_, err := c.DeleteInstance(context.Background(), tc.region, tc.instanceID)
		if err != nil {
			require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC#%d: %s", i+1, tc.name)
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
