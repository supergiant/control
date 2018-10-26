package amazon

import (
	"testing"
)

func TestCreateInstanceStepName(t *testing.T) {
	s := StepCreateInstance{}

	if s.Name() != StepNameCreateEC2Instance {
		t.Errorf("Unexpected step name expected %s actual %s", StepNameCreateEC2Instance, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := StepCreateInstance{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}
