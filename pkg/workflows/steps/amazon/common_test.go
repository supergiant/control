package amazon

import (
	"testing"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestGetEC2(t *testing.T) {
	api, err := GetEC2(steps.AWSConfig{})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if api == nil {
		t.Errorf("Api must not be nil")
	}
}

func TestGetIAM(t *testing.T) {
	api, err := GetIAM(steps.AWSConfig{})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if api == nil {
		t.Errorf("Api must not be nil")
	}
}
