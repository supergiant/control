package certificates

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/pki"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

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

func TestWriteCertificates(t *testing.T) {
	var (
		kubernetesConfigDir = "/etc/kubernetes"
		privateIP           = "10.20.30.40"
		publicIP            = "22.33.44.55"
		userName            = "user"
		password            = "1234"

		r runner.Runner = &fakeRunner{}
	)

	err := templatemanager.Init("../../../../templates")
	if err != nil {
		t.Fatal(err)
	}

	tpl, _ := templatemanager.GetTemplate(StepName)
	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)

	caPair, err := pki.NewCAPair(nil)

	if err != nil {
		t.Errorf("unexpected error creating PKI bundle %v", err)
	}

	cfg := steps.NewConfig("", "", "", profile.Profile{
		K8SServicesCIDR: "10.3.0.0/16",
	})
	// TODO: update tests
	cfg.CertificatesConfig = steps.CertificatesConfig{
		KubernetesConfigDir: kubernetesConfigDir,
		PrivateIP:           privateIP,
		StaticAuth: profile.StaticAuth{
			BasicAuth: []profile.BasicAuthUser{
				{
					Password: "42",
					Name:     "john.doe@sg.io",
					ID:       "john.doe",
					Groups:   []string{"systems:masters", "other:group"},
				},
			},
			Tokens: []profile.TokenAuthUser{
				{
					Token:  "1234",
					Name:   "user@sg.io",
					ID:     "user",
					Groups: []string{"systems:masters", "other:group"},
				},
			},
		},
		CAKey:  string(caPair.Key),
		CACert: string(caPair.Cert),
	}

	cfg.Runner = r
	cfg.Node = model.Machine{
		State:     model.MachineStateActive,
		PrivateIp: privateIP,
		PublicIp:  publicIP,
	}

	nodeRoles := []bool{true, false}

	for _, isMaster := range nodeRoles {
		cfg.IsMaster = isMaster
		task := &Step{
			tpl,
		}

		err = task.Run(context.Background(), output, cfg)

		if err != nil {
			t.Errorf("Unpexpected error while  provision node %v", err)
		}

		if !strings.Contains(output.String(), kubernetesConfigDir) {
			t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
		}

		if !strings.Contains(output.String(), userName) {
			t.Errorf("username %s not found in %s", userName, output.String())
		}

		if !strings.Contains(output.String(), password) {
			t.Errorf("password %s not found in %s", password, output.String())
		}

		if !strings.Contains(output.String(), string(caPair.Key)) {
			t.Errorf("CA key not found in %s", output.String())
		}

		if !strings.Contains(output.String(), string(caPair.Cert)) {
			t.Errorf("CA cert not found in %s", output.String())
		}

		if !strings.Contains(output.String(), privateIP) {
			t.Errorf("Master private ip %s not found in %s",
				privateIP, output.String())
		}

		if !strings.Contains(output.String(), publicIP) {
			t.Errorf("Master public ip %s not found in %s",
				publicIP, output.String())
		}

		if isMaster {
			if !strings.Contains(output.String(), "apiserver-key.pem") {
				t.Errorf("apiserver-key.pem not found in %s",
					output.String())
			}

			if !strings.Contains(output.String(), "apiserver.pem") {
				t.Errorf("apiserver-key.pem not found in %s",
					output.String())
			}

			if !strings.Contains(output.String(), "worker-key.pem") {
				t.Errorf("worker-key.pem not found in %s",
					output.String())
			}

			if !strings.Contains(output.String(), "worker.pem") {
				t.Errorf("worker.pem not found in %s",
					output.String())
			}
		} else {
			if !strings.Contains(output.String(), "worker-key.pem") {
				t.Errorf("worker-key.pem %s not found in %s",
					publicIP, output.String())
			}

			if !strings.Contains(output.String(), "worker.pem") {
				t.Errorf("worker.pem %s not found in %s",
					publicIP, output.String())
			}

			if strings.Contains(output.String(), "apiserver-key.pem") {
				t.Errorf("apiserver-key.pem must not be in in %s",
					output.String())
			}

			if strings.Contains(output.String(), "apiserver.pem") {
				t.Errorf("apiserver.pem must not be in in %s",
					output.String())
			}
		}
		output.Reset()
	}
}

func TestWriteCertificatesError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	task := &Step{
		proxyTemplate,
	}

	cfg := steps.NewConfig("", "", "", profile.Profile{
		K8SServicesCIDR: "10.3.0.0/16",
	})
	cfg.Runner = r
	cfg.AddMaster(&model.Machine{
		State:     model.MachineStateActive,
		PrivateIp: "10.20.30.40",
	})
	err = task.Run(context.Background(), output, cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}

func TestStepName(t *testing.T) {
	s := Step{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s", StepName, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := Step{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}

func TestStep_Rollback(t *testing.T) {
	s := Step{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}

func TestNew(t *testing.T) {
	tpl := template.New("test")
	s := New(tpl)

	if s.template != tpl {
		t.Errorf("Wrong template expected %v actual %v", tpl, s.template)
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
