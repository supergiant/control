package storageclass

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
)

//TODO cleanup all fakerunners and move them to separate package
type fakeRunner struct {
	errMsg string
}

func (f *fakeRunner) Run(command *runner.Command) error {
	if len(f.errMsg) > 0 {
		return errors.New(f.errMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}

func TestStep_Run(t *testing.T) {
	err := templatemanager.Init("../../../../templates")
	require.NoError(t, err)

	tpl, _ := templatemanager.GetTemplate(StepName)
	output := new(bytes.Buffer)

	tt := []struct {
		provider clouds.Name
		output   string
	}{
		{
			clouds.AWS,
			"kubernetes.io/aws-ebs",
		}, {
			clouds.DigitalOcean,
			"local-storage",
		}, {
			clouds.GCE,
			"kubernetes.io/gce-pd",
		},
	}

	for _, tc := range tt {
		cfg := steps.NewConfig("", "", profile.Profile{})
		cfg.Provider = tc.provider
		cfg.Runner = &fakeRunner{}

		task := &Step{
			tpl,
		}

		err = task.Run(context.Background(), output, cfg)
		require.NoError(t, err)

		require.True(t, strings.Contains(output.String(), tc.output))
	}
}

func TestNew(t *testing.T) {
	s := New(&template.Template{})

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.script == nil {
		t.Errorf("Script must not be nil")
	}
}

func TestInit(t *testing.T) {
	templatemanager.SetTemplate(StepName, &template.Template{})
	Init()
	templatemanager.DeleteTemplate(StepName)

	s := steps.GetStep(StepName)

	if s == nil {
		t.Errorf("Step must not be nil")
	}
}

func TestInitPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("recover output must not be nil")
		}
	}()

	Init()

	s := steps.GetStep("not_found.sh.tpl")

	if s == nil {
		t.Error("Step not found")
	}
}

func TestStep_Rollback(t *testing.T) {
	s := &Step{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error when rollback %v", err)
	}
}

func TestStep_Depends(t *testing.T) {
	s := &Step{}

	if deps := s.Depends(); deps == nil || len(deps) != 1 || deps[0] != clustercheck.StepName {
		t.Errorf("Wrong dependencies list expected %v actual %v",
			[]string{clustercheck.StepName}, deps)
	}
}

func TestStep_Name(t *testing.T) {
	s := &Step{}

	if name := s.Name(); name != StepName {
		t.Errorf("Wrong name expected %s actual %s",
			StepName, name)
	}
}

func TestStep_Description(t *testing.T) {
	s := &Step{}

	if desc := s.Description(); desc != "add storage class to k8s" {
		t.Errorf("Wrong desc expected add storage class to "+
			"k8s actual %s", desc)
	}
}
