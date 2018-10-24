package steps

import "testing"

func TestRegisterStep(t *testing.T) {
	var (
		step     Step
		stepName = "test"
	)

	RegisterStep(stepName, step)

	if _, ok := stepMap[stepName]; !ok {
		t.Errorf("step %s not found in step map %v", stepName, stepMap)
	}
}

func TestGetStepNotFound(t *testing.T) {
	var (
		stepName = "not_found"
	)

	if _, ok := stepMap[stepName]; ok {
		t.Errorf("unexpected step %s", stepName)
	}

}
