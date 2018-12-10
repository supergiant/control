package digitalocean

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"

	"github.com/pborman/uuid"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
	"go.uber.org/zap/buffer"
	"net/http"
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

func TestDeleteMachinesStep_Description(t *testing.T) {
	step := NewCreateInstanceStep(time.Second, time.Second)

	if step.Description() != "Create instance in Digital Ocean" {
		t.Errorf("Wrong step description expectetd Create instance"+
			" in Digital Ocean actual %s", step.Description())
	}
}

func TestCreateInstanceStep_Run(t *testing.T) {
	testCases := []struct {
		createKey    *godo.Key
		createKeyErr error
		droplet      *godo.Droplet
		dropletErr   error

		expectError bool
	}{
		{
			createKeyErr: errors.New("error creating key"),
			expectError:  true,
		},
		{
			createKey: &godo.Key{
				ID: 1234,
			},
			createKeyErr: nil,
			dropletErr:   errors.New("error creating droplet"),
			expectError:  true,
		},
		{
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
			dropletErr:  nil,
			expectError: false,
		},
	}

	for _, testCase := range testCases {
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
			err: testCase.dropletErr,
		}

		step := &CreateInstanceStep{
			CheckPeriod:    time.Second * 1,
			DropletTimeout: time.Minute * 1,
			getServices: func(accessToken string) (DropletService, KeyService) {
				return dropletSvc, keySvc
			},
		}

		cfg := steps.NewConfig("", "",
			"", profile.Profile{
				MasterProfiles: make([]profile.NodeProfile, 10),
			})
		cfg.ClusterID = uuid.New()
		cfg.TaskID = uuid.New()
		cfg.IsMaster = true
		err := step.Run(context.Background(), &buffer.Buffer{}, cfg)

		if testCase.expectError && err == nil {
			t.Errorf("Error not must be nil")
		}

		if !testCase.expectError && err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
}
