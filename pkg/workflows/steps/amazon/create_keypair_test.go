package amazon

import (
	"testing"
)

func TestCreateKeyPairStepName(t *testing.T) {
	s := KeyPairStep{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s", StepName, s.Name())
	}
}

func TestCreateKeyPairDepends(t *testing.T) {
	s := KeyPairStep{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}
