package digitalocean

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap/buffer"

	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestNewCreateInstanceStep(t *testing.T) {
	timeout := time.Second * 1
	period := time.Minute * 1
	step := NewCreateInstanceStep(timeout, period)

	if step.DropletTimeout != timeout {
		t.Errorf("Wrong timeout expected %v actual %v",
			timeout, step.DropletTimeout)
	}

	if step.CheckPeriod != period {
		t.Errorf("wrong check period value expected %v actual %v",
			step.CheckPeriod, period)
	}

	if step.getServices == nil {
		t.Errorf("get services must not be nil")
	}

	if client, err := step.getServices("access token"); err == nil {
		t.Errorf("Unexpected values %v %v", client, err)
	}
}

func TestCreateInstanceStep_Rollback(t *testing.T) {
	s := CreateInstanceStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}

func TestCreateInstanceStep_Depends(t *testing.T) {
	step := NewCreateInstanceStep(time.Second, time.Second)

	if step.Depends() != nil {
		t.Error("Create instance step depends must be nil")
	}
}

func TestCreateInstanceStep_Name(t *testing.T) {
	step := NewCreateInstanceStep(time.Second, time.Second)

	if step.Name() != CreateMachineStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			CreateMachineStepName, step.Name())
	}
}

func TestCreateInstanceStep_Description(t *testing.T) {
	step := NewCreateInstanceStep(time.Second, time.Second)

	if step.Description() != "Create instance in Digital Ocean" {
		t.Errorf("Wrong step description expectetd Create instance"+
			" in Digital Ocean actual %s", step.Description())
	}
}

func TestCreateInstanceStep_Run(t *testing.T) {
	testCases := []struct {
		description  string
		createKey    *godo.Key
		createKeyErr error
		droplet      *godo.Droplet
		createErr    error
		getErr       error

		dropletTimeout time.Duration
		period         time.Duration

		isMaster    bool
		expectError bool
	}{
		{
			description:    "create key error",
			dropletTimeout: time.Minute,
			period:         time.Second,
			createKeyErr:   errors.New("error creating key"),
			expectError:    true,
		},
		{
			description:    "create droplet error",
			dropletTimeout: time.Minute,
			period:         time.Second,
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			createErr:    errors.New("error creating droplet"),
			expectError:  true,
		},
		{
			description:    "timeout exceed",
			dropletTimeout: time.Nanosecond,
			period:         time.Nanosecond * 2,
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			droplet: &godo.Droplet{
				ID:     5678,
				Status: "active",
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						{
							Type:      "public",
							IPAddress: "1.2.3.4",
						},
						{
							Type:      "private",
							IPAddress: "5.6.7.8",
						},
					},
				},
			},
			expectError: true,
		},
		{
			description:    "fail on master get",
			dropletTimeout: time.Minute,
			period:         time.Second,
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			droplet: &godo.Droplet{
				ID:     5678,
				Status: "active",
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						{
							Type:      "public",
							IPAddress: "1.2.3.4",
						},
						{
							Type:      "private",
							IPAddress: "5.6.7.8",
						},
					},
				},
			},
			isMaster:    true,
			getErr:      errors.New("get error"),
			expectError: true,
		},
		{
			description:    "fail on master create",
			dropletTimeout: time.Minute,
			period:         time.Second,
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			droplet: &godo.Droplet{
				ID:     5678,
				Status: "active",
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						{
							Type:      "public",
							IPAddress: "1.2.3.4",
						},
						{
							Type:      "private",
							IPAddress: "5.6.7.8",
						},
					},
				},
			},
			isMaster:    true,
			createErr:   nil,
			expectError: false,
		},
		{
			description:    "success",
			dropletTimeout: time.Minute,
			period:         time.Second,
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			droplet: &godo.Droplet{
				ID:     5678,
				Status: "active",
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						{
							Type:      "public",
							IPAddress: "1.2.3.4",
						},
						{
							Type:      "private",
							IPAddress: "5.6.7.8",
						},
					},
				},
			},
			createErr:   nil,
			expectError: false,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		keySvc := &mockKeyService{
			key: testCase.createKey,
			resp: &godo.Response{
				Response: &http.Response{},
			},
			err: testCase.createKeyErr,
		}

		dropletSvc := &mockDropletService{
			droplet: testCase.droplet,
			resp: &godo.Response{
				Response: &http.Response{},
			},
			createErr: testCase.createErr,
			getErr:    testCase.getErr,
		}

		step := &CreateInstanceStep{
			CheckPeriod:    testCase.period,
			DropletTimeout: testCase.dropletTimeout,
			getServices: func(accessToken string) (DropletService, KeyService) {
				return dropletSvc, keySvc
			},
		}

		cfg := steps.NewConfig("", "", profile.Profile{
			MasterProfiles: make([]profile.NodeProfile, 10),
		})
		cfg.ClusterID = uuid.New()
		cfg.TaskID = uuid.New()
		cfg.IsMaster = testCase.isMaster
		err := step.Run(context.Background(), &buffer.Buffer{}, cfg)

		if testCase.expectError && err == nil {
			t.Errorf("Error not must be nil")
		}

		if !testCase.expectError && err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
}
