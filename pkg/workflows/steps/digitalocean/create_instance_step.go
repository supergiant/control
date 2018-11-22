package digitalocean

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestCreateInstanceStep_Rollback(t *testing.T) {
	s := CreateInstanceStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
