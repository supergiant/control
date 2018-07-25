package kubeletconf

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"context"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
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

func TestWriteKubeletConf(t *testing.T) {
	host := "127.0.0.1"
	port := "8080"

	var (
		r            runner.Runner = &fakeRunner{}
		tillerScript               = `server: http://{{ .Host }}:{{ .Port }}`
	)

	proxyTemplate, err := template.New(StepName).Parse(tillerScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)
	cfg := steps.Config{
		KubeletConfConfig: steps.KubeletConfConfig{
			Host: host,
			Port: port,
		},
	}

	j := &Step{
		r,
		proxyTemplate,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), host) {
		t.Errorf("host %s not found in %s", host, output.String())
	}

	if !strings.Contains(output.String(), port) {
		t.Errorf("port %s not found in %s", port, output.String())
	}
}

func TestWriteKubeletConfErr(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)
	cfg := steps.Config{
		KubeletConfConfig: steps.KubeletConfConfig{},
	}

	j := &Step{
		r,
		proxyTemplate,
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
