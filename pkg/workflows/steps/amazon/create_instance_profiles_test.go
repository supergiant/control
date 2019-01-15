package amazon

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	fakeErr = errors.New("fake error")

	awsNotFoundErr = awsError{code: iam.ErrCodeNoSuchEntityException}
)

type awsError struct {
	code string
}

func (e awsError) Error() string {
	return e.code
}
func (e awsError) Code() string {
	return e.code
}
func (e awsError) Message() string {
	return e.code
}
func (e awsError) OrigErr() error {
	return errors.New(e.code)
}

type fakeIAMClient struct {
	iamiface.IAMAPI

	getInstanceProfile    *iam.GetInstanceProfileOutput
	getInstanceProfileErr error

	createInstanceProfile    *iam.CreateInstanceProfileOutput
	createInstanceProfileErr error

	addRoleToInstanceProfile    *iam.AddRoleToInstanceProfileOutput
	addRoleToInstanceProfileErr error

	getRole    *iam.GetRoleOutput
	getRoleErr error

	createRole    *iam.CreateRoleOutput
	createRoleErr error

	getRolePolicy    *iam.GetRolePolicyOutput
	getRolePolicyErr error

	putRolePolicy    *iam.PutRolePolicyOutput
	putRolePolicyErr error
}

func (c *fakeIAMClient) GetInstanceProfileWithContext(aws.Context, *iam.GetInstanceProfileInput, ...request.Option) (*iam.GetInstanceProfileOutput, error) {
	return c.getInstanceProfile, c.getInstanceProfileErr
}
func (c *fakeIAMClient) CreateInstanceProfileWithContext(aws.Context, *iam.CreateInstanceProfileInput, ...request.Option) (*iam.CreateInstanceProfileOutput, error) {
	return c.createInstanceProfile, c.createInstanceProfileErr
}
func (c *fakeIAMClient) AddRoleToInstanceProfileWithContext(aws.Context, *iam.AddRoleToInstanceProfileInput, ...request.Option) (*iam.AddRoleToInstanceProfileOutput, error) {
	return c.addRoleToInstanceProfile, c.addRoleToInstanceProfileErr
}
func (c *fakeIAMClient) GetRoleWithContext(aws.Context, *iam.GetRoleInput, ...request.Option) (*iam.GetRoleOutput, error) {
	return c.getRole, c.getRoleErr
}
func (c *fakeIAMClient) CreateRoleWithContext(aws.Context, *iam.CreateRoleInput, ...request.Option) (*iam.CreateRoleOutput, error) {
	return c.createRole, c.createRoleErr
}
func (c *fakeIAMClient) GetRolePolicyWithContext(aws.Context, *iam.GetRolePolicyInput, ...request.Option) (*iam.GetRolePolicyOutput, error) {
	return c.getRolePolicy, c.getRolePolicyErr
}
func (c *fakeIAMClient) PutRolePolicyWithContext(aws.Context, *iam.PutRolePolicyInput, ...request.Option) (*iam.PutRolePolicyOutput, error) {
	return c.putRolePolicy, c.putRolePolicyErr
}

func TestCreateInstanceProfiles_Run(t *testing.T) {
	for _, tc := range []struct {
		name            string
		iamClientGetter GetIAMFn
		expectedErr     error
		cfg             steps.Config
	}{
		{
			name: "authorization error",
			iamClientGetter: func(config steps.AWSConfig) (iamiface.IAMAPI, error) {
				return &fakeIAMClient{}, fakeErr
			},
			expectedErr: fakeErr,
		},
		{
			name: "create iam profile error",
			iamClientGetter: func(config steps.AWSConfig) (iamiface.IAMAPI, error) {
				return &fakeIAMClient{
					getInstanceProfileErr: fakeErr,
				}, nil
			},
			expectedErr: fakeErr,
		},
		{
			name: "create iam profile",
			iamClientGetter: func(config steps.AWSConfig) (iamiface.IAMAPI, error) {
				return &fakeIAMClient{
					getInstanceProfile: &iam.GetInstanceProfileOutput{
						InstanceProfile: &iam.InstanceProfile{
							Roles: []*iam.Role{{RoleName: aws.String("someRole")}},
						},
					},
				}, nil
			},
			cfg: steps.Config{
				ClusterID: "42",
			},
		},
	} {
		step := NewCreateInstanceProfiles(tc.iamClientGetter)
		err := step.Run(context.Background(), ioutil.Discard, &tc.cfg)

		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
		if err == nil {
			require.Equalf(
				t,
				buildIAMName(tc.cfg.ClusterID, string(node.RoleMaster)),
				tc.cfg.AWSConfig.MastersInstanceProfile,
				"TC: %s", tc.name)

			require.Equalf(
				t,
				buildIAMName(tc.cfg.ClusterID, string(node.RoleNode)),
				tc.cfg.AWSConfig.NodesInstanceProfile,
				"TC: %s", tc.name)
		}
	}
}

func TestCreateIAMInstanceProfile(t *testing.T) {
	for _, tc := range []struct {
		name        string
		iamClient   iamiface.IAMAPI
		expectedErr error
	}{
		{
			name: "get instance profile error",
			iamClient: &fakeIAMClient{
				getInstanceProfileErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{
			name:        "get instance profile error: empty response",
			iamClient:   &fakeIAMClient{},
			expectedErr: ErrEmptyResponse,
		},
		{
			name: "get instance profile error: empty instance profile",
			iamClient: &fakeIAMClient{
				getInstanceProfile: &iam.GetInstanceProfileOutput{},
			},
			expectedErr: ErrEmptyResponse,
		},
		{
			name: "create instance profile error",
			iamClient: &fakeIAMClient{
				getInstanceProfileErr:    awsNotFoundErr,
				createInstanceProfileErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{
			name: "create instance profile",
			iamClient: &fakeIAMClient{
				getInstanceProfileErr: awsNotFoundErr,
				createInstanceProfile: &iam.CreateInstanceProfileOutput{
					InstanceProfile: &iam.InstanceProfile{
						Roles: []*iam.Role{{RoleName: aws.String("someRole")}},
					},
				},
			},
		},
		{
			name: "create instance profile error: empty response",
			iamClient: &fakeIAMClient{
				getInstanceProfileErr: awsNotFoundErr,
			},
			expectedErr: ErrEmptyResponse,
		},
		{
			name: "create instance profile error: empty instance profile",
			iamClient: &fakeIAMClient{
				getInstanceProfileErr: awsNotFoundErr,
				createInstanceProfile: &iam.CreateInstanceProfileOutput{},
			},
			expectedErr: ErrEmptyResponse,
		},
		{
			name: "add role to instance profile error",
			iamClient: &fakeIAMClient{
				getInstanceProfile: &iam.GetInstanceProfileOutput{
					InstanceProfile: &iam.InstanceProfile{},
				},
				addRoleToInstanceProfileErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{
			name: "add role to instance profile",
			iamClient: &fakeIAMClient{
				getInstanceProfile: &iam.GetInstanceProfileOutput{
					InstanceProfile: &iam.InstanceProfile{},
				},
			},
		},
	} {
		err := createIAMInstanceProfile(context.Background(), tc.iamClient, "test")
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}

func TestCreateIAMRole(t *testing.T) {
	for _, tc := range []struct {
		name        string
		iamClient   iamiface.IAMAPI
		expectedErr error
	}{
		{
			name:      "get role",
			iamClient: &fakeIAMClient{},
		},
		{
			name: "get role error",
			iamClient: &fakeIAMClient{
				getRoleErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{
			name: "create role error",
			iamClient: &fakeIAMClient{
				getRoleErr:    awsNotFoundErr,
				createRoleErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
	} {
		err := createIAMRole(context.Background(), tc.iamClient, "test", "policy")
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}

func TestCreateIAMRolePolicy(t *testing.T) {
	for _, tc := range []struct {
		name        string
		iamClient   iamiface.IAMAPI
		expectedErr error
	}{
		{
			name:      "get role policy",
			iamClient: &fakeIAMClient{},
		},
		{
			name: "get role policy error",
			iamClient: &fakeIAMClient{
				getRolePolicyErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
		{
			name: "put role policy error",
			iamClient: &fakeIAMClient{
				getRolePolicyErr: awsNotFoundErr,
				putRolePolicyErr: fakeErr,
			},
			expectedErr: fakeErr,
		},
	} {
		err := createIAMRolePolicy(context.Background(), tc.iamClient, "test", "policy")
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s", tc.name)
	}
}

func TestInitCreateInstanceProfiles(t *testing.T) {
	InitCreateInstanceProfiles(GetIAM)

	s := steps.GetStep(StepNameCreateInstanceProfiles)

	if s == nil {
		t.Errorf("Step %s must not be nil",
			StepNameCreateInstanceProfiles)
	}
}

func TestStepCreateInstanceProfiles_Depends(t *testing.T) {
	step := &StepCreateInstanceProfiles{}

	if deps := step.Depends(); deps != nil {
		t.Errorf("Unexpected deps value %v", deps)
	}
}

func TestStepCreateInstanceProfiles_Description(t *testing.T) {
	step := &StepCreateInstanceProfiles{}

	if desc := step.Description(); desc != "Create EC2 Instance master/node profiles" {
		t.Errorf("Unexpected description value %s "+
			"expected Create EC2 Instance master/node profiles", desc)
	}
}

func TestStepCreateInstanceProfiles_Rollback(t *testing.T) {
	step := &StepCreateInstanceProfiles{}

	if err := step.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error value %v", err)
	}
}

func TestStepCreateInstanceProfiles_Name(t *testing.T) {
	step := &StepCreateInstanceProfiles{}

	if name := step.Name(); name != StepNameCreateInstanceProfiles {
		t.Errorf("Wrong step name expected %s actual %s",
			StepNameCreateInstanceProfiles, step.Name())
	}
}

func TestNewCreateInstanceProfilesError(t *testing.T) {
	fn := func(steps.AWSConfig) (iamiface.IAMAPI, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewCreateInstanceProfiles(fn)

	if step == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if step.GetIAM == nil {
		t.Errorf("GetIAM must not be nil")
	}

	if api, err := step.GetIAM(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateInstanceProfiles(t *testing.T) {
	step := NewCreateInstanceProfiles(GetIAM)

	if step == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if step.GetIAM == nil {
		t.Errorf("GetIAM must not be nil")
	}

	if api, err := step.GetIAM(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}
