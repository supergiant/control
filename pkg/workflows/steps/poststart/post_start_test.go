package poststart

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubelet"
)

type fakeRunner struct {
	errMsg  string
	timeout time.Duration
}

func (f *fakeRunner) Run(command *runner.Command) error {
	if len(f.errMsg) > 0 {
		return errors.New(f.errMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	// Simulate command latency
	time.Sleep(f.timeout)

	return err
}

func TestPostStartMaster(t *testing.T) {
	port := "8080"
	username := "john"
	rbacEnabled := true
	r := &fakeRunner{}

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	cfg := steps.NewConfig("test", "", "test", profile.Profile{
		MasterProfiles: []profile.NodeProfile{{}},
	})
	cfg.IsMaster = true
	cfg.PostStartConfig = steps.PostStartConfig{
		Host:        "127.0.0.1",
		Port:        port,
		Username:    username,
		RBACEnabled: rbacEnabled,
		Timeout:     120,
	}
	cfg.Node = node.Node{
		PrivateIp: "10.20.30.40",
	}
	cfg.Runner = r

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while master node %v", err)
	}

	if !strings.Contains(output.String(), port) {
		t.Errorf("port %s not found in %s", port, output.String())
	}

	if !strings.Contains(output.String(), username) && rbacEnabled {
		t.Errorf("rbac section not found in %s", output.String())
	}
}

func TestPostStartNode(t *testing.T) {
	r := &fakeRunner{}

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	cfg := steps.NewConfig("test", "", "test", profile.Profile{
		MasterProfiles: []profile.NodeProfile{{}},
	})
	cfg.PostStartConfig = steps.PostStartConfig{
		Host:    "127.0.0.1",
		Timeout: 120,
	}
	cfg.Node = node.Node{
		PrivateIp: "10.20.30.40",
	}
	cfg.Runner = r

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), "docker ps") {
		t.Errorf("docker ps command  not found in %s", output.String())
	}
}

func TestPostStartError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	j := &Step{
		proxyTemplate,
	}

	cfg := &steps.Config{
		PostStartConfig: steps.PostStartConfig{
			Timeout: 1,
		},
		Runner: r,
	}
	err = j.Run(context.Background(), output, cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}

func TestPostStartTimeout(t *testing.T) {
	port := "8080"
	username := "john"
	rbacEnabled := true

	r := &fakeRunner{
		errMsg:  "",
		timeout: time.Second * 2,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	j := &Step{
		proxyTemplate,
	}

	cfg := steps.NewConfig("", "", "", profile.Profile{})
	cfg.PostStartConfig = steps.PostStartConfig{
		true,
		"127.0.0.1",
		port,
		username,
		rbacEnabled,
		1,
	}
	cfg.Runner = r

	err = j.Run(context.Background(), output, cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if errors.Cause(err) != context.DeadlineExceeded {
		t.Errorf("Wrong error cause expected %v actual %v", context.DeadlineExceeded,
			err)
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

	if len(s.Depends()) != 1 && s.Depends()[0] != kubelet.StepName {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{kubelet.StepName})
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

	if s.script != tpl {
		t.Errorf("Wrong template expected %v actual %v", tpl, s.script)
	}
}

func TestInit(t *testing.T) {
	Init()

	s := steps.GetStep(StepName)

	if s == nil {
		t.Error("Step not found")
	}
}
