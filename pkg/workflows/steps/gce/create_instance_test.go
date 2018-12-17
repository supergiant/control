package gce

import (
	"bytes"
	"context"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
	"strings"
	"testing"
	"time"
)

func TestCreateInstanceStep_Run(t *testing.T) {
	testCases := []struct {
		description string
		getSvcErr   error

		image       *compute.Image
		getImageErr error

		machineType       *compute.MachineType
		getMachineTypeErr error

		insertErr error

		instance       *compute.Instance
		getInstanceErr error

		setMetadataErr error

		errMsg string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description: "get image error",
			getImageErr: errors.New("message2"),
			errMsg:      "message2",
		},
		{
			description: "get machine type error",
			image: &compute.Image{
				Id: 1234,
			},
			getMachineTypeErr: errors.New("message3"),
			errMsg:            "message3",
		},
		{
			description: "insert instance error",
			image: &compute.Image{
				Id: 1234,
			},
			machineType: &compute.MachineType{
				SelfLink: "https://itsme.com",
			},
			insertErr: errors.New("message4"),
			errMsg:    "message4",
		},
		{
			description: "get instance error",
			image: &compute.Image{
				Id: 1234,
			},
			machineType: &compute.MachineType{
				SelfLink: "https://itsme.com",
			},
			getInstanceErr: errors.New("message5"),
			errMsg:         "message5",
		},
		{
			description: "set metadata error",
			image: &compute.Image{
				Id: 1234,
			},
			machineType: &compute.MachineType{
				SelfLink: "https://itsme.com",
			},
			instance: &compute.Instance{
				Status:   "RUNNING",
				Metadata: &compute.Metadata{},
			},
			setMetadataErr: errors.New("message6"),
			errMsg:         "message6",
		},
		{
			description: "timeout",
			image: &compute.Image{
				Id: 1234,
			},
			machineType: &compute.MachineType{
				SelfLink: "https://itsme.com",
			},
			instance: &compute.Instance{
				Metadata: &compute.Metadata{},
			},
			errMsg: sgerrors.ErrTimeoutExceeded.Error(),
		},
		{
			description: "success",
			image: &compute.Image{
				Id: 1234,
			},
			machineType: &compute.MachineType{
				SelfLink: "https://itsme.com",
			},
			instance: &compute.Instance{
				Status:   "RUNNING",
				Metadata: &compute.Metadata{},
				NetworkInterfaces: []*compute.NetworkInterface{
					{
						NetworkIP: "10.20.30.40",
						AccessConfigs: []*compute.AccessConfig{
							{
								NatIP: "11.22.33.44",
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)

		step := &CreateInstanceStep{
			checkPeriod:     time.Nanosecond,
			instanceTimeout: time.Millisecond * 1,
			getComputeSvc: func(ctx context.Context,
				config steps.GCEConfig) (*computeService, error) {
				return &computeService{
					getFromFamily: func(context.Context, steps.GCEConfig) (*compute.Image, error) {
						return testCase.image, testCase.getImageErr
					},
					getMachineTypes: func(context.Context, steps.GCEConfig) (*compute.MachineType, error) {
						return testCase.machineType, testCase.getMachineTypeErr
					},
					insertInstance: func(context.Context, steps.GCEConfig, *compute.Instance) (*compute.Operation, error) {
						return nil, testCase.insertErr
					},
					getInstance: func(context.Context, steps.GCEConfig, string) (*compute.Instance, error) {
						return testCase.instance, testCase.getInstanceErr
					},
					setInstanceMetadata: func(context.Context, steps.GCEConfig, string, *compute.Metadata) (*compute.Operation, error) {
						return nil, testCase.setMetadataErr
					},
				}, testCase.getSvcErr
			},
		}

		for _, role := range []bool{true, false} {
			config := steps.NewConfig("", "",
				"", profile.Profile{})
			config.TaskID = uuid.New()
			config.ClusterName = util.RandomString(8)
			config.ClusterID = uuid.New()[:8]
			config.IsMaster = role

			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				for {
					select {
					case <-config.NodeChan():
					case <-ctx.Done():
					}
				}
			}()

			err := step.Run(ctx, &bytes.Buffer{}, config)
			cancel()

			if err == nil && testCase.errMsg != "" {
				t.Errorf("error must not be nil")
			}

			if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
				t.Errorf("Error message %s does not contain %s",
					err.Error(), testCase.errMsg)
			}
		}
	}
}

func TestNewCreateInstanceStep(t *testing.T) {
	period := time.Second * 1
	timeout := time.Second * 2

	s, err := NewCreateInstanceStep(period, timeout)

	if s == nil {
		t.Error("Step must not be nil")
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if s.checkPeriod != period {
		t.Errorf("Wrong period expected %v actual %v",
			period, s.checkPeriod)
	}

	if s.instanceTimeout != timeout {
		t.Errorf("Wrong timeout expected %v actual %v",
			timeout, s.instanceTimeout)
	}

	if s.getComputeSvc == nil {
		t.Errorf("get compute service must not be nil")
	}
}

func TestCreateInstanceStep_Depends(t *testing.T) {
	s := CreateInstanceStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("dependencies must be nil")
	}
}

func TestCreateInstanceStep_Name(t *testing.T) {
	s := CreateInstanceStep{}

	if name := s.Name(); name != CreateInstanceStepName {
		t.Errorf("Wrong name expected %s actual %s",
			CreateInstanceStepName, name)
	}
}

func TestCreateInstanceStep_Rollback(t *testing.T) {
	s := CreateInstanceStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestCreateInstanceStep_Description(t *testing.T) {
	s := CreateInstanceStep{}

	if desc := s.Description(); desc != "Google compute engine step for creating instance" {
		t.Errorf("Wrong description expected %s actual "+
			"Google compute engine step for creating instance", desc)
	}
}
