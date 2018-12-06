package clusterservices

import (
	"bytes"
	"context"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/flannel"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
	"testing"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	fakeErr = errors.New("fake err")
)

type fakeRunner struct {
	err error
}
func (r *fakeRunner) Run(command *runner.Command) error {
	return r.err
}

func TestStepName(t *testing.T) {
	s := Step{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s", StepName, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := Step{}
	require.Nil(t, s.Depends())
}

func TestStep_Rollback(t *testing.T) {
	s := Step{}
	require.Nil(t, s.Depends())
}

func TestNew(t *testing.T) {
	tpl := template.New("test")
	s := New(tpl)

	if s.script != tpl {
		t.Errorf("Wrong template expected %v actual %v", tpl, s.script)
	}
}

func TestInit(t *testing.T) {
	templatemanager.SetTemplate(StepName, &template.Template{})
	Init()
	templatemanager.DeleteTemplate(StepName)

	s := steps.GetStep(StepName)

	if s == nil {
		t.Error("Step not found")
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

func TestStep_Run(t *testing.T) {
	err := templatemanager.Init("../../../../templates")
	if err != nil {
		t.Fatal(err)
	}
	authorizedKeys.Init()
	downloadk8sbinary.Init()
	certificates.Init()
	manifest.Init()
	flannel.Init()
	docker.Init()
	kubelet.Init()
	cni.Init()

	for _, tc := range []struct{
		name string
		cfg *steps.Config
		expectedErr error
	}{
		{
			name: "runner: error",
			cfg: &steps.Config{
				ClusterName: "test",
				IsMaster: true,
				AWSConfig: steps.AWSConfig{
					Region: "1",
				},
				Runner: &fakeRunner{
					err: fakeErr,
				},
			},
			expectedErr: fakeErr,
		},
		{
			name: "success",
			cfg: &steps.Config{
				ClusterName: "test",
				IsMaster: true,
				AWSConfig: steps.AWSConfig{
					Region: "1",
				},
				Runner: &fakeRunner{},
			},
		},
	}{
		tpl, _ := templatemanager.GetTemplate(StepName)
		require.NotNilf(t, tpl, "TC: %q: template %s not found", tc.name, StepName)

		s := Step{script:tpl}
		out := &bytes.Buffer{}

		err := s.Run(context.Background(), out, tc.cfg)
		require.Equalf(t, tc.expectedErr, errors.Cause(err), "TC: %s: %s", tc.name, err)
	}
}
