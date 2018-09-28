package poststart

import (
	"bytes"
	"context"
	"io"
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
