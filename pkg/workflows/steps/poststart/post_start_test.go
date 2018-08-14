package poststart

import (
	"bytes"

	"context"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"time"

	"github.com/supergiant/supergiant/pkg/runner"
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

func TestPostStart(t *testing.T) {
	port := "8080"
	username := "john"
	rbacEnabled := true

	var (
		r               runner.Runner = &fakeRunner{}
		postStartScript               = `#!/bin/bash
until $(curl --output /dev/null --silent --head --fail http://127.0.0.1:{{ .Port }}); do printf '.'; sleep 5; done
curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://127.0.0.1:{{ .Port }}/api/v1/namespaces
/opt/bin/kubectl config set-cluster default-cluster --server="127.0.0.1:{{ .Port }}"
/opt/bin/kubectl config set-context default-system --cluster=default-cluster --user=default-admin
/opt/bin/kubectl config use-context default-system

{{if .RBACEnabled }}
/opt/bin/kubectl create clusterrolebinding kubelet-binding --clusterrole=system:node --user=kubelet
/opt/bin/kubectl create clusterrolebinding system:dns-admin-binding --clusterrole=cluster-admin --user=system:dns
/opt/bin/kubectl create clusterrolebinding add-ons-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
/opt/bin/kubectl create clusterrolebinding default-user-cluster-admin --clusterrole=cluster-admin --user={{ .Username }}
/opt/bin/kubectl create clusterrolebinding default-kube-system-admin --clusterrole=cluster-admin --serviceaccount=default:default --namespace=kube-system
{{end}}

/opt/bin/kubectl create -f /etc/kubernetes/addons/kube-dns.yaml
/opt/bin/kubectl create -f /etc/kubernetes/addons/cluster-monitoring
/opt/bin/kubectl create -f /etc/kubernetes/addons/default-storage-class.yaml
/opt/bin/helm init`
	)

	proxyTemplate, err := template.New(StepName).Parse(postStartScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy templatemanager %v", err)
	}

	output := new(bytes.Buffer)
	cfg := &steps.Config{
		PostStartConfig: steps.PostStartConfig{
			"127.0.0.1",
			port,
			username,
			rbacEnabled,
			120,
		},
		Runner: r,
	}

	j := &Step{
		proxyTemplate,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), port) {
		t.Errorf("port %s not found in %s", port, output.String())
	}

	if !strings.Contains(output.String(), username) && rbacEnabled {
		t.Errorf("rbac section not found in %s", output.String())
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

	cfg := &steps.Config{
		PostStartConfig: steps.PostStartConfig{
			"127.0.0.1",
			port,
			username,
			rbacEnabled,
			1,
		},
		Runner: r,
	}
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
