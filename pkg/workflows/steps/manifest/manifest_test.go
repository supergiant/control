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

func TestWriteManifest(t *testing.T) {
	var (
		kubernetesVersion   = "1.8.7"
		kubernetesConfigDir = "/kubernetes/conf/dir"
		RBACEnabled         = true
		etcdHost            = "127.0.0.1"
		etcdPort            = "2379"
		privateIpv4         = "12.34.56.78"
		providerString      = "aws"
		masterHost          = "127.0.0.1"
		masterPort          = "8080"

		r                   runner.Runner = &fakeRunner{}
		writeManifestScript               = `KUBERNETES_MANIFESTS_DIR={{ .KubernetesConfigDir }}/manifests

mkdir -p ${KUBERNETES_MANIFESTS_DIR}
    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-apiserver.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-apiserver
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-apiserver
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - apiserver
    - --bind-address=0.0.0.0
    - --etcd-servers=http://{{ .EtcdHost }}:{{ .EtcdPort }}
	- --allow-privileged=true
    {{if .RBACEnabled }}- --authorization-mode=Node,RBAC{{end}}
    - --service-cluster-ip-range=10.3.0.0/24
    - --secure-port=443
    - --v=2
    - --advertise-address={{ .PrivateIpv4 }}
    - --admission-control=NamespaceLifecycle,NamespaceExists,LimitRanger,ServiceAccount,ResourceQuota,DefaultStorageClass{{if .RBACEnabled }},NodeRestriction{{end}}
    - --tls-cert-file=/etc/kubernetes/ssl/apiserver.pem
    - --tls-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --client-ca-file=/etc/kubernetes/ssl/ca.pem
    - --service-account-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --basic-auth-file=/etc/kubernetes/ssl/basic_auth.csv
    - --token-auth-file=/etc/kubernetes/ssl/known_tokens.csv
    - --kubelet-preferred-address-types=InternalIP,Hostname,ExternalIP
    - --storage-backend=etcd2
    {{- .ProviderString }}
    ports:
    - containerPort: 443
      hostPort: 443
      name: https
    - containerPort: 8080
      hostPort: 8080
      name: local
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/kubernetes/addons
      name: api-addons-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /etc/kubernetes/addons
    name: api-addons-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-controller-manager.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-controller-manager
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-controller-manager
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - controller-manager
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    - --service-account-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --root-ca-file=/etc/kubernetes/ssl/ca.pem
    - --v=2
    - --cluster-cidr=10.244.0.0/14
    - --allocate-node-cidrs=true
    {{- .ProviderString }}
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10252
      initialDelaySeconds: 15
      timeoutSeconds: 1
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-scheduler.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-scheduler
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-scheduler
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - scheduler
    - --v=2
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10251
      initialDelaySeconds: 15
      timeoutSeconds: 1
EOF

    cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-proxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-proxy
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }}
    command:
    - /hyperkube
    - proxy
    - --v=2
    - --master=http://{{ .MasterHost }}:{{ .MasterPort }}
    - --proxy-mode=iptables
    securityContext:
      privileged: true
    volumeMounts:
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
`
	)

	proxyTemplate, err := template.New("manifest").Parse(writeManifestScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)
	cfg := steps.Config{
		ManifestConfig: steps.ManifestConfig{
			KubernetesVersion:   kubernetesVersion,
			KubernetesConfigDir: kubernetesConfigDir,
			RBACEnabled:         RBACEnabled,
			EtcdHost:            etcdHost,
			EtcdPort:            etcdPort,
			PrivateIpv4:         privateIpv4,
			MasterHost:          masterHost,
			MasterPort:          masterPort,
			ProviderString:      providerString,
		},
	}

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	err = j.Run(context.Background(), cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), kubernetesVersion) {
		t.Errorf("kubernetes version dir %s not found in %s", kubernetesVersion, output.String())
	}

	if RBACEnabled && !strings.Contains(output.String(), "RBAC") {
		t.Errorf("RBAC not found in %s", output.String())
	}

	if !strings.Contains(output.String(), masterHost) {
		t.Errorf("master host %s not found in %s", masterHost, output.String())
	}

	if !strings.Contains(output.String(), masterPort) {
		t.Errorf("master port %s not found in %s", masterPort, output.String())
	}

	if !strings.Contains(output.String(), etcdHost) {
		t.Errorf("etcd host %s not found in %s", etcdHost, output.String())
	}

	if !strings.Contains(output.String(), etcdPort) {
		t.Errorf("etcd port %s not found in %s", etcdPort, output.String())
	}

	if !strings.Contains(output.String(), privateIpv4) {
		t.Errorf("private ipv4 %s not found in %s", privateIpv4, output.String())
	}

	if !strings.Contains(output.String(), providerString) {
		t.Errorf("provider string %s not found in %s", providerString, output.String())
	}
}

func TestWriteManifestError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New("manifest").Parse("")
	output := new(bytes.Buffer)
	cfg := steps.Config{
		ManifestConfig: steps.ManifestConfig{},
	}

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	err = j.Run(context.Background(), cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
