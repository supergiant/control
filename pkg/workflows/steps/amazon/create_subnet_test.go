package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockSubnetSvc struct {
	mock.Mock
}

func (m *mockSubnetSvc) CreateSubnetWithContext(ctx aws.Context,
	req *ec2.CreateSubnetInput, opts ...request.Option) (*ec2.CreateSubnetOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.CreateSubnetOutput)
	if !ok {
		return nil, args.Error(1)
	}

	return val, args.Error(1)
}

func (m *mockSubnetSvc) ModifySubnetAttributeWithContext(ctx aws.Context, req *ec2.ModifySubnetAttributeInput,
	opts ...request.Option) (*ec2.ModifySubnetAttributeOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.ModifySubnetAttributeOutput)
	if !ok {
		return nil, args.Error(1)
	}

	return val, args.Error(1)
}

type mockAccountGetter struct {
	mock.Mock
}

func (m *mockAccountGetter) Get(ctx context.Context, accName string) (*model.CloudAccount, error) {
	args := m.Called(ctx, accName)
	val, ok := args.Get(0).(*model.CloudAccount)
	if !ok {
		return nil, args.Error(1)
	}

	return val, args.Error(1)
}

type mockZoneGetter struct {
	zones []string
	err   error
}

func (m *mockZoneGetter) GetZones(context.Context, steps.Config) ([]string, error) {
	return m.zones, m.err
}

func TestCreateSubnetStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		getSvcErr     error
		zoneGetterErr error
		zoneGetter    account.ZonesGetter
		vpcCIDR       string

		getZoneErr  error
		getZoneResp []string

		createSubnetErr error
		createSubnet    *ec2.CreateSubnetOutput

		errMsg string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description:   "zone getter error",
			zoneGetterErr: errors.New("message2"),
			errMsg:        "message2",
		},
		{
			description: "invalid vpc cidr",
			getZoneResp: []string{"us-west-1a", "us-west-1b"},
			vpcCIDR:     "10.0.0.0/36",
			errMsg:      "Error parsing VPC",
		},
		{
			description: "calculating subnet",
			getZoneResp: []string{"us-west-1a", "us-west-1b"},
			vpcCIDR:     "10.3.5.1/32",
			errMsg:      "Calculating",
		},
		{
			description: "get zones error",
			vpcCIDR:     "10.0.0.0/16",
			getZoneErr:  errors.New("message3"),
			errMsg:      "message3",
		},
		{
			description:     "create subnets error",
			getZoneResp:     []string{"us-west-1a", "us-west-1b"},
			vpcCIDR:         "10.0.0.0/16",
			createSubnetErr: errors.New("message4"),
			errMsg:          "message4",
		},
		{
			description: "success",
			vpcCIDR:     "10.0.0.0/16",
			getZoneResp: []string{"us-west-1a", "us-west-1b"},
			createSubnet: &ec2.CreateSubnetOutput{
				Subnet: &ec2.Subnet{
					SubnetId: aws.String("1234"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockSubnetSvc{}
		svc.On("CreateSubnetWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.createSubnet, testCase.createSubnetErr)
		svc.On("ModifySubnetAttributeWithContext", mock.Anything, mock.Anything,
			mock.Anything).Return(nil, nil)

		zoneGetter := &mockZoneGetter{
			zones: testCase.getZoneResp,
			err:   testCase.getZoneErr,
		}

		step := &CreateSubnetsStep{
			getSvc: func(steps.AWSConfig) (subnetSvc, error) {
				return svc, testCase.getSvcErr
			},
			zoneGetterFactory: func(context.Context, accountGetter, *steps.Config) (account.ZonesGetter, error) {
				return zoneGetter, testCase.zoneGetterErr
			},
		}

		config, err := steps.NewConfig("clusterName", "", profile.Profile{})

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		config.AWSConfig.VPCCIDR = testCase.vpcCIDR

		err = step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Wrong error message %v expected to have %s",
				err, testCase.errMsg)
		}
	}
}

func TestInitCreateSubnet(t *testing.T) {
	InitCreateSubnet(GetEC2, nil)

	s := steps.GetStep(StepCreateSubnets)

	if s == nil {
		t.Errorf("Step %s not found", StepCreateSubnets)
	}
}

func TestNewCreateSubnetStep2(t *testing.T) {
	testCases := []struct {
		description string
		getAccErr   error
		acc         *model.CloudAccount
		errMSg      string
	}{
		{
			description: "Get cloud account error",
			getAccErr:   errors.New("message1"),
			errMSg:      "message1",
		},
		{
			description: "unsupported cloud provider",
			acc: &model.CloudAccount{
				Provider: clouds.Name("unsupported"),
			},
			errMSg: account.ErrUnsupportedProvider.Error(),
		},
		{
			description: "success",
			acc: &model.CloudAccount{
				Provider: clouds.AWS,
			},
		},
	}

	for _, testCase := range testCases {
		accGetter := &mockAccountGetter{}
		accGetter.On("Get", mock.Anything,
			mock.Anything).Return(testCase.acc, testCase.getAccErr)

		step := NewCreateSubnetStep(nil, accGetter)

		svc, err := step.zoneGetterFactory(context.Background(), accGetter,
			&steps.Config{})

		if testCase.errMSg != "" && err == nil {
			t.Error("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMSg) {
			t.Errorf("Wrong error message %v must contain %s",
				err, testCase.errMSg)
		}

		if testCase.errMSg == "" && svc == nil {
			t.Errorf("Service must not be nil")
		}
	}
}

func TestNewCreateSubnetStep(t *testing.T) {
	accSvc := &account.Service{}

	step := NewCreateSubnetStep(GetEC2, accSvc)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.accountGetter != accSvc {
		t.Errorf("account gettere value is wrong exepected %v actual %v",
			step.accountGetter, accSvc)
	}

	if step.getSvc == nil {
		t.Errorf("Wrong get EC2 function must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateSubnetStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	accSvc := &account.Service{}

	step := NewCreateSubnetStep(fn, accSvc)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.accountGetter != accSvc {
		t.Errorf("account getter value is wrong exepected %v actual %v",
			step.accountGetter, accSvc)
	}

	if step.getSvc == nil {
		t.Errorf("Wrong get getSvc function must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestCreateSubnetsStep_Name(t *testing.T) {
	step := &CreateSubnetsStep{}

	if step.Name() != StepCreateSubnets {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateSubnets, step.Name())
	}
}

func TestCreateSubnetsStep_Description(t *testing.T) {
	step := &CreateSubnetsStep{}

	if step.Description() != "Step create subnets in all availability zones for Region" {
		t.Errorf("Wrong step description expected Step create subnets in "+
			"all availability zones for Region actual %s", step.Description())
	}
}

func TestCreateSubnetsStep_Depends(t *testing.T) {
	step := &CreateSubnetsStep{}

	if deps := step.Depends(); deps != nil {
		t.Errorf("%s dependencies must be nil", StepCreateSubnets)
	}
}

func TestCreateSubnetsStep_Rollback(t *testing.T) {
	step := &CreateSubnetsStep{}

	if err := step.Rollback(context.Background(),
		&bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v when rollback", err)
	}
}
