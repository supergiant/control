package drain

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
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwdnJH7+a1ckkqu6fc4hLfsGlNvGhRqPVTRq/YKLbQu3XsNjf
oIC4+w81FABCNvNPajwcTRTfQv6wtmj/m60vI3LyJ23DMThyv0M8czkg40EDn6ub
UOOdR7IZAqkRbc+XCEYK5/8i2j2fuHa/oLKbILvKr0xmSG/Mwso1T2bH4YO9GfEJ
ftZJ/tSk5Yzpb0Oo5cWb2V96OvJDG/MEb5BXPKbOMU8ngnx2gqs+jh1p1fQwrb0G
JjJgA83sySjr8AHkelH9QNS8HvgfHlL8bX2nmFGeaVrRxRXzspJndN15NKUWAMLU
gvWVno3ulqKMkzH+ehK6SFCuKoF5icxtqBx4wwIDAQABAoIBAAwbSulVsRja+BRI
1OKFR5nCBEx7KMRdpQusuPkTEriKXCcqVEUU5PihCYKXRYtjBLmwyV+zBwKLH4Q0
6InTdhczrZXyz/b5/IifbV4Q2lH3FH/bWtbhcEgzAkbdQj5mcZtNrI6yq32PzbLa
j7s8jF2t/MmX7udlPBeKQ2wTEjauRKy5yAS3nshHtsZ1yTa5wFdtrdv3Ulav8w+O
zcoQFWP5CM63MKi4vQ8wN/Nl+wYFlFztVJMbqfUhcWHkkQwyZRha3PfjyrbShwDS
rq7v9UoJ8Zuun0FYQkkwavWBvrG+euhFrDJSV/77m2rkvWJyymfzlNDo4rNESjqs
U+RSppECgYEA8daWHxlou0F2iSEAU3m2femj3a+MrDSgYN5Gdg4MH12+X5WC6hvs
JAvqEi7NkBqozei+FhaKuB4mkONBBwjx7mIwT/J/QFy9Ejr4Lg3BJeyachrB5WxG
6PcZn8IzPXXpyxBzZUMeo17Gif5/dQwTI8CHyFz7jGT0BKnJfpPOf0kCgYEAzTPQ
tJNb83HYAbziKzlacGZFUDrjHDTV/wr1NnMpj551jJVrFTBvz99wWv4uvMXTtiyH
QZftY+wyLbfgFGPAQnp0x7c8/YC656bEehYEa4ecrPy14seLk/4KAI6IMiUbspe1
iuUdSlsS3FaUSE1BM35QCpgaTPs3rNsZGjH626sCgYEA2BLtYG34aE6uFQl6XBsE
VW26LmkaHAaNQN94PyR/6kp8vLQ+GuPF0dMfWQ2eNuHK7ubDZ8LOQIEX3h5dzGZO
mrn6BoRY8+2oNLCha6x4ZWUH/Wkw0sYyeRXGPDpsQ76lm/xfzhrxNfCJHWRZBwA5
3Zi4+OkzC5Zre/sjf8eaGZkCgYBgqo0h091YNIQWZX2B8TW6h2MVpXgBfJ5m1Cmp
6dxlTLeBb44PYE779P0/0EgCI4tVYWqiKsjo7obA5MMJt+gFKRzETHzNywvBPt2F
ycNxSGQ1VaL1Xx1QrTbXBk4AmVyP6EncUYxXz8l1xM97s/EIKfPY2chiBWI36srL
fUn4mwKBgQCpELvH91vYlzECe+uM5SPHEKWXZf27+I+ABOttVlxMYomjSfWnuKRp
aJTJFasfzhumpKl9+T3HZ3WNx2YzzggkDs3l/CuBeJcuHmWBZ/99VSAcQ4v7bJez
ftubvUfbIYGZDXxwqpBifsVm22U3lTggt1MSzSFzUaS6RkvOw+wW2g==
-----END RSA PRIVATE KEY-----`
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

func TestDrain(t *testing.T) {
	var (
		expected               = "10.20.30.40"
		r        runner.Runner = &fakeRunner{}
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

	cfg, err := steps.NewConfig("",
		"", profile.Profile{})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	cfg.Masters = steps.NewMap(map[string]*model.Machine{
		"master-0": {
			Name:     "master-0",
			State:    model.MachineStateActive,
			PublicIp: "10.20.30.40",
		},
	})
	cfg.Kube.SSHConfig = model.SSHConfig{
		BootstrapPrivateKey: privateKey,
		User:                "root",
		Port:                ssh.DefaultPort,
	}

	task := &Step{
		script: tpl,
		getRunner: func(masterIP string, config *steps.Config) (runner.Runner, error) {
			return r, nil
		},
	}

	cfg.DrainConfig.PrivateIP = expected
	err = task.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), expected) {
		t.Errorf("nodename not found %s in %s",
			expected, output.String())
	}
}

func TestErrors(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	task := &Step{
		script: proxyTemplate,
		getRunner: func(masterIP string, config *steps.Config) (runner.Runner, error) {
			return r, nil
		},
	}

	cfg, err := steps.NewConfig("",
		"", profile.Profile{})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

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
		t.Errorf("Error message expected to contain %s actual %s",
			errMsg, err.Error())
	}
}

func TestStepName(t *testing.T) {
	s := Step{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s",
			StepName, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := Step{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v",
			s.Depends(), []string{})
	}
}

func TestStep_Rollback(t *testing.T) {
	s := Step{}
	err := s.Rollback(context.Background(), ioutil.Discard,
		&steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}

func TestNew(t *testing.T) {
	tpl := template.New("test")
	s := New(tpl)

	if s.script != tpl {
		t.Errorf("Wrong template expected %v actual %v",
			tpl, s.script)
	}

	if s.getRunner == nil {
		t.Errorf("getRunner function must not be nil")
	}

	cfg := &steps.Config{
		Kube: model.Kube{
			SSHConfig: model.SSHConfig{
				BootstrapPrivateKey: privateKey,
				User:                "root",
				Port:                ssh.DefaultPort,
			},
		},
	}

	if _, err := s.getRunner("10.20.30.40", cfg); err != nil {
		t.Errorf("Unexpected error when get runner %v", err)
	}
}

func TestNewErr(t *testing.T) {
	tpl := template.New("test")
	s := New(tpl)

	if s.script != tpl {
		t.Errorf("Wrong template expected %v actual %v",
			tpl, s.script)
	}

	if s.getRunner == nil {
		t.Errorf("getRunner function must not be nil")
	}

	cfg := &steps.Config{}

	if _, err := s.getRunner("10.20.30.40", cfg); err == nil {
		t.Errorf("Error must not be nil")
	}
}

func TestInit(t *testing.T) {
	templatemanager.SetTemplate(StepName, &template.Template{})
	Init()

	s := steps.GetStep(StepName)

	if s == nil {
		t.Error("Step not found")
	}

	templatemanager.DeleteTemplate(StepName)
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

func TestStep_Description(t *testing.T) {
	s := &Step{}

	if desc := s.Description(); desc != "drain resources from a node" {
		t.Errorf("Wrong desription expected %s actual %s",
			"drain resources from a node", desc)
	}
}
