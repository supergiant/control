package tiller

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

func TestPostStart(t *testing.T) {
	host := "127.0.0.1"
	port := "8080"
	username := "john"
	rbacEnabled := true

	var (
		r               runner.Runner = &fakeRunner{}
		postStartScript               = `#!/bin/bash
until $(curl --output /dev/null --silent --head --fail http://{{ .Host }}:{{ .Port }}); do printf '.'; sleep 5; done
curl -XPOST -H 'Content-type: application/json' -d'{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"kube-system"}}' http://{{ .Host }}:{{ .Port }}/api/v1/namespaces
/opt/bin/kubectl config set-cluster default-cluster --server="{{ .Host }}:{{ .Port }}"
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
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)
	cfg := steps.Config{
		PostStartConfig: steps.PostStartConfig{
			host,
			port,
			username,
			rbacEnabled,
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

	if !strings.Contains(output.String(), host) {
		t.Errorf("host %s not found in %s", host, output.String())
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

	cfg := steps.Config{
		PostStartConfig: steps.PostStartConfig{},
		Runner:          r,
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
