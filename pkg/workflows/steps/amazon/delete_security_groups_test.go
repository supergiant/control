package amazon

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"strings"
	"testing"
	"time"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockDeleteSecGroupSvc struct {
	mock.Mock
}

func (m *mockDeleteSecGroupSvc) DescribeSecurityGroups(
	input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DescribeSecurityGroupsOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockDeleteSecGroupSvc) RevokeSecurityGroupIngressWithContext(ctx aws.Context,
	input *ec2.RevokeSecurityGroupIngressInput, opts ...request.Option) (*ec2.RevokeSecurityGroupIngressOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*ec2.RevokeSecurityGroupIngressOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockDeleteSecGroupSvc) DeleteSecurityGroupWithContext(ctx aws.Context,
	input *ec2.DeleteSecurityGroupInput, opts ...request.Option) (*ec2.DeleteSecurityGroupOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*ec2.DeleteSecurityGroupOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteSecurityGroup_Run(t *testing.T) {
	testCases := []struct {
		description string

		masterSecGroupId string
		nodeSecGroupId   string

		getSvcErr error

		describeMasterOutput *ec2.DescribeSecurityGroupsOutput
		describeMasterErr    error

		describeNodeOutput *ec2.DescribeSecurityGroupsOutput
		describeNodeErr    error

		revokeMasterErr error
		revokeNodeErr   error

		deleteMasterErr error
		deleteNodeErr   error
		errMsg          string
	}{
		{
			description: "skip delete",
		},
		{
			description: "get service error",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			getSvcErr: errors.New("message1"),
			errMsg:    "message1",
		},
		{
			description: "describe master err",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterErr: errors.New("message2"),
			errMsg:            "message2",
		},
		{
			description: "describe master not found",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{},
			},
			errMsg: sgerrors.ErrNotFound.Error(),
		},
		{
			description: "describe node err",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeErr: errors.New("message3"),
			errMsg:          "message3",
		},
		{
			description: "describe node not found",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{},
			},
			errMsg: sgerrors.ErrNotFound.Error(),
		},
		{
			description: "revoke master error",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("node"),
					},
				},
			},

			revokeMasterErr: errors.New("message4"),
			errMsg: "message4",
		},
		{
			description: "revoke node error",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("node"),
					},
				},
			},

			revokeNodeErr: errors.New("message5"),
			errMsg: "message5",
		},
		{
			description: "delete master error",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("node"),
					},
				},
			},

			deleteMasterErr: errors.New("message6"),
			errMsg: "message6",
		},
		{
			description: "delete node error",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("node"),
					},
				},
			},

			deleteNodeErr: errors.New("message7"),
			errMsg: "message7",
		},
		{
			description: "success",

			masterSecGroupId: "1234",
			nodeSecGroupId:   "5678",

			describeMasterOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("master"),
					},
				},
			},

			describeNodeOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupName: aws.String("node"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockDeleteSecGroupSvc{}
		svc.On("DescribeSecurityGroups",
			mock.Anything).Return(testCase.describeMasterOutput,
			testCase.describeMasterErr).Once()
		svc.On("DescribeSecurityGroups",
			mock.Anything).Return(testCase.describeNodeOutput,
			testCase.describeNodeErr).Once()

		svc.On("RevokeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything,
			testCase.revokeMasterErr).Once()
		svc.On("RevokeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything,
			testCase.revokeNodeErr).Once()

		svc.On("DeleteSecurityGroupWithContext",
			mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything,
			testCase.deleteMasterErr).Once()
		svc.On("DeleteSecurityGroupWithContext",
			mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything,
			testCase.deleteNodeErr).Once()

		step := &DeleteSecurityGroup{
			getSvc: func(cfg steps.AWSConfig) (deleteSecurityGroupService, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				MastersSecurityGroupID: testCase.masterSecGroupId,
				NodesSecurityGroupID:   testCase.nodeSecGroupId,
			},
		}

		deleteSecGroupAttemptCount = 1
		deleteSecGroupTimeout = time.Nanosecond

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be empty")
			return
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %v must contain %s", err,
				testCase.errMsg)
		}
	}
}


func TestInitDeleteSecurityGroup(t *testing.T) {
	InitDeleteSecurityGroup(GetEC2)

	s := steps.GetStep(DeleteSecurityGroupsStepName)

	if s == nil {
		t.Errorf("step must not be nil")
	}
}

func TestNewDeleteSecurityGroupsStep(t *testing.T) {
	s := NewDeleteSecurityGroupService(GetEC2)

	if s == nil {
		t.Errorf("step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("get service func must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewDeleteSecurityGroupsStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDeleteSecurityGroupService(fn)

	if s == nil {
		t.Errorf("step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("get service func must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestDeleteSecurityGroup_Depends(t *testing.T) {
	s := &DeleteSecurityGroup{}

	if deps := s.Depends(); deps == nil || len(deps) != 1 || deps[0] != DeleteClusterMachinesStepName {
		t.Errorf("Wrong dependencies expected %v actual %v",
			[]string{DeleteClusterMachinesStepName}, deps)
	}
}

func TestDeleteSecurityGroup_Rollback(t *testing.T) {
	s := &DeleteSecurityGroup{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error while roll back %b", err)
	}
}


func TestDeleteSecurityGroup_Description(t *testing.T) {
	s := &DeleteSecurityGroup{}

	if desc := s.Description(); desc != "Deletes security groups" {
		t.Errorf("Wrong description value expected Deletes " +
			"security groups actual %s", desc)
	}
}

func TestDeleteSecurityGroup_Name(t *testing.T) {
	s := &DeleteSecurityGroup{}

	if name := s.Name(); name != DeleteSecurityGroupsStepName {
		t.Errorf("Wrong name expected %s actual %s",
			DeleteSecurityGroupsStepName, name)
	}
}
